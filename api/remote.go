package api

import (
	"fmt"
	"sync"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
	"github.com/pmylund/go-cache"
)

type remoteConnection struct {
	client     *registry.Client
	layerCache *cache.Cache
}

func NewRemoteRegistry() *registryApi {
	conn := &remoteConnection{
		client:     registry.NewClient(),
		layerCache: cache.New(cacheDuration, cacheCleanupInterval),
	}
	reg := newRegistryApi(conn)

	return reg
}

func (rc *remoteConnection) GetTags(name string) (registry.TagMap, error) {
	auth, err := rc.client.Hub.GetReadToken(name)
	if err != nil {
		return nil, err
	}

	return rc.client.Repository.ListTags(name, auth)
}

func (rc *remoteConnection) Search(name string) (*registry.SearchResults, error) {
	return rc.client.Search.Query(name, 1, 10)
}

func (rc *remoteConnection) Status() (Status, error) {
	var s Status

	s.Message = "Connected to Registry"
	s.Service = "Registry Image Manager"
	return s, nil
}

func (rc *remoteConnection) GetImageLayers(name, tag string) ([]*registry.ImageMetadata, error) {
	var wg sync.WaitGroup

	auth, err := rc.client.Hub.GetReadToken(name)
	if err != nil {
		return nil, err
	}

	id, err := rc.client.Repository.GetImageID(name, tag, auth)
	if err != nil {
		return nil, err
	}

	layerIDs, err := rc.client.Image.GetAncestry(id, auth)
	if err != nil {
		return nil, err
	}

	list := make([]*registry.ImageMetadata, len(layerIDs))

	for i, layerID := range layerIDs {
		wg.Add(1)
		go func(idx int, lid string) {
			defer wg.Done()
			var err error
			var m *registry.ImageMetadata

			val, found := rc.layerCache.Get(lid)
			if found {
				m = val.(*registry.ImageMetadata)
			} else {
				m, err = rc.client.Image.GetMetadata(lid, auth)
				if err == nil && m != nil {
					rc.layerCache.Set(lid, m, cache.DefaultExpiration)
				}
			}

			list[idx] = m
		}(i, layerID)
	}
	wg.Wait()

	// Return an error if any of the layers are missing
	for i, layer := range list {
		if layer == nil {
			return nil, fmt.Errorf("Error retrieving layer: %s", layerIDs[i])
		}
	}

	return list, nil
}
