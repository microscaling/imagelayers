package registry

import (
	"fmt"
)

// ImageService gives access to the /repositories portion of the Registry API.
type RepositoryService struct {
	client *Client
}

// TagMap represents a collection of image tags and their corresponding image IDs.
type TagMap map[string]string

// Retrieves a list of tags for the specified repository.
//
// Docker Registry API docs: https://docs.docker.com/reference/api/registry_api/#list-repository-tags
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

// Retrieves the ID of the image identified by the given repository and tag.
//
// Docker Registry API docs: https://docs.docker.com/reference/api/registry_api/#get-image-id-for-a-particular-tag
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

// Sets a tag for the the specified repository and image ID.
//
// Docker Registry API docs: https://docs.docker.com/reference/api/registry_api/#set-a-tag-for-a-specified-image-id
func (r *RepositoryService) SetTag(repo, imageID, tag string, auth Authenticator) error {
	path := fmt.Sprintf("repositories/%s/tags/%s", repo, tag)
	req, err := r.client.newRequest("PUT", path, auth, imageID)
	if err != nil {
		return err
	}

	_, err = r.client.do(req, nil)
	return err
}

// Deletes a tag from the specified repository.
//
// Docker Registry API docs: https://docs.docker.com/reference/api/registry_api/#delete-a-repository-tag
func (r *RepositoryService) DeleteTag(repo, tag string, auth Authenticator) error {
	path := fmt.Sprintf("repositories/%s/tags/%s", repo, tag)
	req, err := r.client.newRequest("DELETE", path, auth, nil)
	if err != nil {
		return err
	}

	_, err = r.client.do(req, nil)
	return err
}

// Deletes the specified repository.
//
// Docker Registry API docs: https://docs.docker.com/reference/api/registry_api/#delete-a-repository
func (r *RepositoryService) Delete(repo string, auth Authenticator) error {
	path := fmt.Sprintf("repositories/%s/", repo)
	req, err := r.client.newRequest("DELETE", path, auth, nil)
	if err != nil {
		return err
	}

	_, err = r.client.do(req, nil)
	return err
}
