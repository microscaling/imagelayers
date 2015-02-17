package api

import (
	"github.com/CenturyLinkLabs/docker-reg-client/registry"
)

type remoteConnection struct {
	client *registry.Client
	auth   *registry.TokenAuth
}

func NewRemoteRegistry() *registryApi {
	conn := &remoteConnection{client: registry.NewClient()}
	reg := newRegistryApi(conn)

	return reg
}

func (rc *remoteConnection) Status() (Status, error) {
	var s Status

	s.Message = "Connected to Registry"
	s.Service = "Registry Image Manager"
	return s, nil
}

func (rc *remoteConnection) Connect(repo string) error {
	auth, err := rc.client.Hub.GetReadToken(repo)
	if err != nil {
		return err
	}

	rc.auth = auth
	return nil
}

func (rc *remoteConnection) GetImageID(name string, tag string) (string, error) {
	id, err := rc.client.Repository.GetImageID(name, tag, rc.auth)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (rc *remoteConnection) GetAncestry(id string) ([]string, error) {
	layers, err := rc.client.Image.GetAncestry(id, rc.auth)
	if err != nil {
		return nil, err
	}

	return layers, nil

}

func (rc *remoteConnection) GetMetadata(layer string) (*registry.ImageMetadata, error) {
	metadata, err := rc.client.Image.GetMetadata(layer, rc.auth)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}
