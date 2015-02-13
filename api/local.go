package api

import (
	"github.com/CenturyLinkLabs/docker-reg-client/registry"
)

type localMgr struct {
}

func NewLocalMgr() *localMgr {
	return new(localMgr)
}

func (r *localMgr) Analyze(images []string) ([]registry.ImageMetadata, error) {
	var image registry.ImageMetadata
	list := make([]registry.ImageMetadata, 1)
	list[0] = image
	return list, nil
}

func (r *localMgr) Version() (Version, error) {
	var v Version
	v.Version = "0.1.0"
	v.Name = "Local Image Manager"
	return v, nil
}
