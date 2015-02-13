package api

import (
	"sync"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
)

type registryMgr struct {
}

func NewRegistry() *api {
	mgr := new(registryMgr)
	reg := newApi(mgr)

	return reg
}

func (mgr *registryMgr) Status() (Status, error) {
	var v Status

	v.Version = "0.1.0"
	v.Name = "Registry Image Manager"
	return v, nil
}

func (mgr *registryMgr) Analyze(images []string) ([]*registry.ImageMetadata, error) {
	repo := "centurylink/image-graph"
	client := registry.NewClient()
	auth,_ := client.Hub.GetReadToken(repo)

	id, _ := client.Repository.GetImageID(repo, "latest", auth)

	layers, _ := client.Image.GetAncestry(id, auth)
	// Goroutine for metadata
	list := mgr.loadImageData(client, auth, layers)

	return list, nil
}

func (mgr *registryMgr) loadImageData(client *registry.Client, auth *registry.TokenAuth, layers []string) []*registry.ImageMetadata {
	var wg sync.WaitGroup
	list := make([]*registry.ImageMetadata, len(layers))

	for i, layerID := range layers {
		wg.Add(1)
		go func(idx int, layer string) {
			defer wg.Done()
			m, _ := client.Image.GetMetadata(layerID, auth)
			list[idx] = m
		}(i, layerID)
	}
	wg.Wait()
	return list
}
