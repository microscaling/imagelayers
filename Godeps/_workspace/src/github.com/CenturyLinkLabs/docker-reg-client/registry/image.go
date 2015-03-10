package registry

import (
	"fmt"
	"io"
	"time"
)

// ImageService gives access to the /images portion of the Registry API.
type ImageService struct {
	client *Client
}

type ImageMetadata struct {
	ID              string           `json:"id"`
	Parent          string           `json:"parent"`
	Comment         string           `json:"Comment"`
	Created         time.Time        `json:"created"`
	Container       string           `json:"container"`
	ContainerConfig ContainerConfig  `json:"container_config"`
	DockerVersion   string           `json:"docker_version"`
	Author          string           `json:"author"`
	Config          *ContainerConfig `json:"config"`
	Architecture    string           `json:"architecture"`
	OS              string           `json:"os"`
	Size            int64            `json:"Size"`
}

type ContainerConfig struct {
	HostName        string              `json:"Hostname"`
	DomainName      string              `json:"Domainname"`
	User            string              `json:"User"`
	Memory          int64               `json:"Memory"`
	MemorySwap      int64               `json:"MemorySwap"`
	CPUShares       int64               `json:"CpuShares"`
	CPUSet          string              `json:"Cpuset"`
	AttachStdin     bool                `json:"AttachStdin"`
	AttachStdout    bool                `json:"AttachStdout"`
	AttachStderr    bool                `json:"AttachStderr"`
	PortSpecs       []string            `json:"PortSpecs"`
	ExposedPorts    map[string]struct{} `json:"ExposedPorts"`
	TTY             bool                `json:"Tty"`
	OpenStdin       bool                `json:"OpenStdin"`
	StdinOnce       bool                `json:"StdinOnce"`
	Env             []string            `json:"Env"`
	Cmd             []string            `json:"Cmd"`
	DNS             []string            `json:"Dns"`
	Image           string              `json:"Image"`
	Volumes         map[string]struct{} `json:"Volumes"`
	VolumesFrom     string              `json:"VolumesFrom"`
	WorkingDir      string              `json:"WorkingDir"`
	Entrypoint      []string            `json:"Entrypoint"`
	NetworkDisabled bool                `json:"NetworkDisabled"`
	SecurityOpts    []string            `json:"SecurityOpts"`
	OnBuild         []string            `json:"OnBuild"`
}

func (i *ImageService) GetLayer(imageID string, layer io.Writer, auth Authenticator) error {
	path := fmt.Sprintf("images/%s/layer", imageID)
	req, err := i.client.newRequest("GET", path, auth, layer)
	if err != nil {
		return err
	}

	_, err = i.client.do(req, layer)
	return err
}

func (i *ImageService) AddLayer(imageID string, layer io.Reader, auth Authenticator) error {
	path := fmt.Sprintf("images/%s/layer", imageID)
	req, err := i.client.newRequest("PUT", path, auth, layer)
	if err != nil {
		return err
	}

	_, err = i.client.do(req, nil)
	return err
}

func (i *ImageService) GetMetadata(imageID string, auth Authenticator) (*ImageMetadata, error) {
	path := fmt.Sprintf("images/%s/json", imageID)
	req, err := i.client.newRequest("GET", path, auth, nil)
	if err != nil {
		return nil, err
	}

	meta := &ImageMetadata{}
	_, err = i.client.do(req, meta)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (i *ImageService) AddMetadata(imageID string, metadata *ImageMetadata, auth Authenticator) error {
	path := fmt.Sprintf("images/%s/json", imageID)
	req, err := i.client.newRequest("PUT", path, auth, metadata)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	_, err = i.client.do(req, nil)
	return err
}

func (i *ImageService) GetAncestry(imageID string, auth Authenticator) ([]string, error) {
	path := fmt.Sprintf("images/%s/ancestry", imageID)
	req, err := i.client.newRequest("GET", path, auth, nil)
	if err != nil {
		return nil, err
	}

	var ancestry []string
	_, err = i.client.do(req, &ancestry)
	if err != nil {
		return nil, err
	}

	return ancestry, nil
}
