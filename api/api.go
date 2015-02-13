package api

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
)

type Version struct {
	Version string `json:"version"`
	Name	string `json:"name"`
}

type LayerManager interface {
	Version() (Version, error)
	Analyze([]string) ([]registry.ImageMetadata, error)
}

type api struct {
	manager LayerManager
}

func NewApi(mgr LayerManager) *api {
	newApi := new(api)
	newApi.manager = mgr

	return newApi
}

func (a *api) Routes() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc {
		"GET": {
			"/version":a.version,
		},
		"POST": {
			"/analyze": a.analyze,
		},
	}
}


func (a *api) analyze(w http.ResponseWriter, r *http.Request) {
	test := new(registry.ImageMetadata)
	test.ID = "Me"

	json.NewEncoder(w).Encode(test)
}

func (a *api) version(w http.ResponseWriter, r *http.Request) {
	version, _ := a.manager.Version()
	log.Printf("Version: %s", version.Name)

	json.NewEncoder(w).Encode(version)
}
