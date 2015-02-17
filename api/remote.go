package api

import (
	"sync"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
)

type registryMgr struct {
	client *registry.Client
}

func NewRemoteRegistry() *registryApi {
	mgr := &registryMgr{client: registry.NewClient()}
	reg := newRegistryApi(mgr)

	return reg
}

func (mgr *registryMgr) Status() (Status, error) {
	var s Status

	s.Message = "Connected to Registry"
	s.Service = "Registry Image Manager"
	return s, nil
}

func (mgr *registryMgr) Analyze(images []string) ([]*registry.ImageMetadata, error) {
	list := make([]*registry.ImageMetadata,0)
	client := mgr.client
	// Goroutine for metadata
	for _, image := range images {
		auth, _ := client.Hub.GetReadToken(image)
		id, _ := client.Repository.GetImageID(image, "latest", auth)
		layers, _ := client.Image.GetAncestry(id, auth)
		metadata := mgr.loadImageData(client, auth, layers)
		list = append(list, metadata...)
	}

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
