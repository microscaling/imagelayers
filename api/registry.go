package api

import (
	"github.com/CenturyLinkLabs/docker-reg-client/registry"
)

type registryMgr struct {
}

func NewRegistryMgr() *registryMgr {
	return new(registryMgr)
}

func (r *registryMgr) Analyze(images []string) ([]registry.ImageMetadata, error) {
	var image registry.ImageMetadata
	list := make([]registry.ImageMetadata, 1)
	list[0] = image
	return list, nil
}

func (r *registryMgr) Version() (Version, error) {
	var v Version
	v.Version = "0.1.0"
	v.Name = "Registry Image Manager"
	return v, nil
}
