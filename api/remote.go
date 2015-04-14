package api

import (
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
		layerCache: cache.New(cacheDuration, 0),
	}
	reg := newRegistryApi(conn)

	return reg
}

func (rc *remoteConnection) GetTags(name string) (registry.TagMap, error) {
	auth, err := rc.getAuth(name)
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

	auth, err := rc.getAuth(name)
	if err != nil {
		return nil, err
	}

	id, err := rc.getImageID(auth, name, tag)
	if err != nil {
		return nil, err
	}

	layerIDs, err := rc.getAncestry(auth, id)
	if err != nil {
		return nil, err
	}

	list := make([]*registry.ImageMetadata, len(layerIDs))

	for i, layerID := range layerIDs {
		wg.Add(1)
		go func(idx int, lid string) {
			defer wg.Done()
			var m *registry.ImageMetadata

			val, found := rc.layerCache.Get(lid)
			if found {
				m = val.(*registry.ImageMetadata)
			} else {
				m, _ = rc.getMetadata(auth, lid)
				rc.layerCache.Set(lid, m, cache.DefaultExpiration)
			}

			list[idx] = m
		}(i, layerID)
	}
	wg.Wait()

	return list, nil
}

func (rc *remoteConnection) getAuth(repo string) (registry.Authenticator, error) {
	return rc.client.Hub.GetReadToken(repo)
}

func (rc *remoteConnection) getImageID(auth registry.Authenticator, name string, tag string) (string, error) {
	id, err := rc.client.Repository.GetImageID(name, tag, auth)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (rc *remoteConnection) getAncestry(auth registry.Authenticator, id string) ([]string, error) {
	layers, err := rc.client.Image.GetAncestry(id, auth)
	if err != nil {
		return nil, err
	}

	return layers, nil

}

func (rc *remoteConnection) getMetadata(auth registry.Authenticator, layer string) (*registry.ImageMetadata, error) {
	metadata, err := rc.client.Image.GetMetadata(layer, auth)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}
