package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
)

type Status struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

type Request struct {
	Repos []Repo `json:"repos"`
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

func (reg *registryApi) Routes() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"GET": {
			"/status": reg.handleStatus,
		},
		"POST": {
			"/analyze": reg.handleAnalysis,
		},
	}
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
	layers, _ := reg.inspectImages(request.Repos)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(layers)
}

func (reg *registryApi) inspectImages(images []Repo) ([]*registry.ImageMetadata, error) {
	list := make([]*registry.ImageMetadata, 0)

	// Goroutine for metadata
	for _, image := range images {
		reg.connection.Connect(image.Name)
		id, _ := reg.connection.GetImageID(image.Name, image.Tag)
		layers, _ := reg.connection.GetAncestry(id)
		metadata := reg.loadMetaData(layers)
		list = append(list, metadata...)
	}

	return list, nil
}

func (reg *registryApi) loadMetaData(layers []string) []*registry.ImageMetadata {
	var wg sync.WaitGroup
	list := make([]*registry.ImageMetadata, len(layers))

	for i, layerID := range layers {
		wg.Add(1)
		go func(idx int, layer string) {
			defer wg.Done()
			m, _ := reg.connection.GetMetadata(layerID)
			list[idx] = m
		}(i, layerID)
	}
	wg.Wait()
	return list
}
