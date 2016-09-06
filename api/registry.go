package api // import "github.com/microscaling/imagelayers/api"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
	"github.com/gorilla/mux"
	"github.com/microscaling/imagelayers/server"
	"github.com/pmylund/go-cache"
)

const (
	cacheDuration        = 15 * time.Minute
	cacheCleanupInterval = 5 * time.Minute
)

type Status struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

type Request struct {
	Repos []Repo `json:"repos"`
}

type Response struct {
	Repo   *Repo                     `json:"repo"`
	Layers []*registry.ImageMetadata `json:"layers"`
	Status int                       `json:"status"`
}

type Repo struct {
	Name  string `json:"name"`
	Tag   string `json:"tag"`
	Size  int64  `json:"size"`
	Count int    `json:"count"`
}

type RegistryConnection interface {
	Status() (Status, error)
	GetTags(string) (registry.TagMap, error)
	Search(string) (*registry.SearchResults, error)
	GetImageLayers(name, tag string) ([]*registry.ImageMetadata, error)
}

type registryApi struct {
	connection RegistryConnection
	imageCache *cache.Cache
}

func newRegistryApi(conn RegistryConnection) *registryApi {
	return &registryApi{
		connection: conn,
		imageCache: cache.New(cacheDuration, cacheCleanupInterval),
	}
}

func (reg *registryApi) Routes(context string, router *server.Router) {
	routes := server.RouteMap{
		"GET": {
			"/status":                     reg.handleStatus,
			"/search":                     reg.handleSearch,
			"/images/{front}/tags":        reg.handleTags,
			"/images/{front}/{tail}/tags": reg.handleTags,
		},
		"POST": {
			"/analyze": reg.handleAnalysis,
		},
	}

	router.AddCorsRoutes(context, routes)
}

func (reg *registryApi) handleTags(w http.ResponseWriter, r *http.Request) {
	image := mux.Vars(r)["front"]
	tail := mux.Vars(r)["tail"]

	if tail != "" {
		image = image + "/" + tail
	}

	res, err := reg.connection.GetTags(image)
	if err != nil {
		respondToError(w, err)
		return
	}

	respondWithJSON(w, res)
}

func (reg *registryApi) handleSearch(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("name")

	res, err := reg.connection.Search(value)
	if err != nil {
		respondToError(w, err)
		return
	}

	respondWithJSON(w, res)
}

func (reg *registryApi) handleStatus(w http.ResponseWriter, r *http.Request) {
	res, err := reg.connection.Status()
	if err != nil {
		respondToError(w, err)
		return
	}
	log.Printf("Status: %s", res.Service)

	respondWithJSON(w, res)
}

func (reg *registryApi) handleAnalysis(w http.ResponseWriter, r *http.Request) {
	var request Request

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondToError(w, err)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		respondToError(w, err)
		return
	}

	res := reg.inspectImages(request.Repos)

	respondWithJSON(w, res)
}

func (reg *registryApi) inspectImages(images []Repo) []*Response {
	var wg sync.WaitGroup
	list := make([]*Response, len(images))

	for i, image := range images {
		wg.Add(1)
		go func(idx int, img Repo) {
			defer wg.Done()
			var resp *Response

			key := fmt.Sprintf("%s:%s", img.Name, img.Tag)
			val, found := reg.imageCache.Get(key)
			if found {
				resp = val.(*Response)
			} else {
				resp = reg.loadMetaData(img)
				if resp.Status == http.StatusOK {
					reg.imageCache.Set(key, resp, cache.DefaultExpiration)
				}
			}

			list[idx] = resp
		}(i, image)
	}

	wg.Wait()
	return list
}

func (reg *registryApi) loadMetaData(repo Repo) *Response {
	resp := new(Response)
	resp.Repo = &repo

	layers, err := reg.connection.GetImageLayers(repo.Name, repo.Tag)
	if err == nil {
		resp.Status = http.StatusOK
		resp.Layers = layers
		resp.Repo.Count = len(resp.Layers)

		for _, layer := range resp.Layers {
			resp.Repo.Size += layer.Size
		}
	} else {
		switch e := err.(type) {
		case registry.RegistryError:
			resp.Status = e.Code
		default:
			resp.Status = http.StatusInternalServerError
		}
		log.Printf("Error: %s", err)
	}

	return resp
}

func respondToError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func respondWithJSON(w http.ResponseWriter, o interface{}) {
	if err := json.NewEncoder(w).Encode(o); err != nil {
		respondToError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}
