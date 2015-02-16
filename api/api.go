package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
)

type Status struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

type Repositories struct {
	Repos []string `json:"repos"`
}

type LayerManager interface {
	Status() (Status, error)
	Analyze([]string) ([]*registry.ImageMetadata, error)
}

type api struct {
	manager LayerManager
}

func newApi(mgr LayerManager) *api {
	return &api{manager: mgr}
}

func (a *api) Routes() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"GET": {
			"/status": a.status,
		},
		"POST": {
			"/analyze": a.analyze,
		},
	}
}

func (a *api) analyze(w http.ResponseWriter, r *http.Request) {
	var repos Repositories

	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &repos)
	if err != nil {
		panic(err)
	}

	layers, _ := a.manager.Analyze(repos.Repos)

	json.NewEncoder(w).Encode(layers)
}

func (a *api) status(w http.ResponseWriter, r *http.Request) {
	status, _ := a.manager.Status()
	log.Printf("Status: %s", status.Name)

	json.NewEncoder(w).Encode(status)
}
