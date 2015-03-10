package registry

import (
	"fmt"
)

// ImageService gives access to the /repositories portion of the Registry API.
type RepositoryService struct {
	client *Client
}

type TagMap map[string]string

func (r *RepositoryService) ListTags(repo string, auth Authenticator) (TagMap, error) {
	path := fmt.Sprintf("repositories/%s/tags", repo)
	req, err := r.client.newRequest("GET", path, auth, nil)
	if err != nil {
		return nil, err
	}

	tags := TagMap{}
	_, err = r.client.do(req, &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *RepositoryService) GetImageID(repo, tag string, auth Authenticator) (string, error) {
	path := fmt.Sprintf("repositories/%s/tags/%s", repo, tag)
	req, err := r.client.newRequest("GET", path, auth, nil)
	if err != nil {
		return "", err
	}

	var id string
	_, err = r.client.do(req, &id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *RepositoryService) SetTag(repo, imageID, tag string, auth Authenticator) error {
	path := fmt.Sprintf("repositories/%s/tags/%s", repo, tag)
	req, err := r.client.newRequest("PUT", path, auth, imageID)
	if err != nil {
		return err
	}

	_, err = r.client.do(req, nil)
	return err
}

func (r *RepositoryService) DeleteTag(repo, tag string, auth Authenticator) error {
	path := fmt.Sprintf("repositories/%s/tags/%s", repo, tag)
	req, err := r.client.newRequest("DELETE", path, auth, nil)
	if err != nil {
		return err
	}

	_, err = r.client.do(req, nil)
	return err
}

func (r *RepositoryService) Delete(repo string, auth Authenticator) error {
	path := fmt.Sprintf("repositories/%s/", repo)
	req, err := r.client.newRequest("DELETE", path, auth, nil)
	if err != nil {
		return err
	}

	_, err = r.client.do(req, nil)
	return err
}
