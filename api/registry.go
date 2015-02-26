package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
	"github.com/CenturyLinkLabs/imagelayers/server"
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
}

type Repo struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

type RegistryConnection interface {
	Status() (Status, error)
	Connect(string) error
	GetImageID(string, string) (string, error)
	GetAncestry(string) ([]string, error)
	GetMetadata(string) (*registry.ImageMetadata, error)
}

type registryApi struct {
	connection RegistryConnection
}

func newRegistryApi(conn RegistryConnection) *registryApi {
	return &registryApi{connection: conn}
}

func (reg *registryApi) Routes(context string, router *server.Router) {
	routes := map[string]map[string]http.HandlerFunc{
		"GET": {
			"/status": reg.handleStatus,
		},
		"POST": {
			"/analyze": reg.handleAnalysis,
		},
	}

	router.AddCorsRoutes(context, routes)
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
	list := make([]*Response, len(images))

	// Goroutine for metadata
	for idx, image := range images {
		reg.connection.Connect(image.Name)
		id, _ := reg.connection.GetImageID(image.Name, image.Tag)
		layers, _ := reg.connection.GetAncestry(id)
		metadata := reg.loadMetaData(image, layers)
		list[idx] = metadata
	}

	return list, nil
}

func (reg *registryApi) loadMetaData(repo Repo, layers []string) *Response {
	var wg sync.WaitGroup
	res := new(Response)
	res.Repo = &repo


	list := make([]*registry.ImageMetadata, len(layers))

	for i, layerID := range layers {
		wg.Add(1)
		go func(idx int, layer string) {
			defer wg.Done()
			m, _ := reg.connection.GetMetadata(layer)
			list[idx] = m
		}(i, layerID)
	}
	wg.Wait()
	res.Layers = list
	return res
}
