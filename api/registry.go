package api // import "github.com/CenturyLinkLabs/imagelayers/api"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
	"github.com/CenturyLinkLabs/imagelayers/server"
	"github.com/gorilla/mux"
	"github.com/pmylund/go-cache"
)

const (
	cacheDuration = 15 * time.Minute
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
		imageCache: cache.New(cacheDuration, 0),
	}
}

func (reg *registryApi) Routes(context string, router *server.Router) {
	routes := map[string]map[string]http.HandlerFunc{
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

	res, _ := reg.connection.GetTags(image)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (reg *registryApi) handleSearch(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("name")

	res, _ := reg.connection.Search(value)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (reg *registryApi) handleStatus(w http.ResponseWriter, r *http.Request) {
	status, _ := reg.connection.Status()
	log.Printf("Status: %s", status.Service)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (reg *registryApi) handleAnalysis(w http.ResponseWriter, r *http.Request) {
	var request Request

	body, err := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(body, &request)
	if err != nil {
		panic(err)
	}
	res, _ := reg.inspectImages(request.Repos)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (reg *registryApi) inspectImages(images []Repo) ([]*Response, error) {
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
	return list, nil
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
