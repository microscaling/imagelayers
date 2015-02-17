package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
)

type Status struct {
	Message string `json:"message"`
	Service   string `json:"service"`
}

type Request struct {
	Repos []string `json:"repos"`
}

type LayerManager interface {
	Status() (Status, error)
	Analyze([]string) ([]*registry.ImageMetadata, error)
}

type registryApi struct {
	manager LayerManager
}

func newRegistryApi(mgr LayerManager) *registryApi {
	return &registryApi{manager: mgr}
}

func (reg *registryApi) Routes() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"GET": {
			"/status": reg.status,
		},
		"POST": {
			"/analyze": reg.analyze,
		},
	}
}

func (reg *registryApi) analyze(w http.ResponseWriter, r *http.Request) {
	var repos Request

	body, err := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(body, &repos)
	if err != nil {
		panic(err)
	}
	layers, _ := reg.manager.Analyze(repos.Repos)

	json.NewEncoder(w).Encode(layers)
}

func (reg *registryApi) status(w http.ResponseWriter, r *http.Request) {
	status, _ := reg.manager.Status()
	log.Printf("Status: %s", status.Service)

	json.NewEncoder(w).Encode(status)
}
