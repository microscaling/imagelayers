package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

type mockLayerManager struct {
	mock.Mock
}

func (m *mockLayerManager) Analyze(images []string) ([]*registry.ImageMetadata, error) {
	args := m.Mock.Called(images)
	return args.Get(0).([]*registry.ImageMetadata), nil
}

func (m *mockLayerManager) Status() (Status, error) {
	args := m.Mock.Called()
	return args.Get(0).(Status), nil
}

func TestStatus (t *testing.T) {
	status := Status{Status: "active", Name: "foo"}
	jsonTest, _ := json.Marshal(status)

	assert.Equal(t, `{"status":"active","name":"foo"}`, string(jsonTest))
}

func TestRespositories (t *testing.T) {
	imageList := []string{"foo", "bar", "baz"}
	repos := Repositories{Repos: imageList}
	jsonTest, _ := json.Marshal(repos)

	assert.Equal(t, `{"repos":["foo","bar","baz"]}`, string(jsonTest))
}

func TestRoutes (t *testing.T) {
	layerManager := new(mockLayerManager)
	api := newApi(layerManager)
	routes := api.Routes()

	assert.NotNil(t, routes["GET"])
	assert.NotNil(t, routes["GET"]["/status"])
	assert.NotNil(t, routes["POST"])
	assert.NotNil(t, routes["POST"]["/analyze"])
}

func TestAnalyze(t *testing.T) {
	layerManager := new(mockLayerManager)
	api := newApi(layerManager)
	inBody, images := "{\"repos\":[\"foo\"]}", []string{"foo"}

	resp := make([]*registry.ImageMetadata,1)
	req, _ := http.NewRequest("POST", "http://localhost/analyze", strings.NewReader(inBody))
	w := httptest.NewRecorder()

	layerManager.On("Analyze", images).Return(resp, nil)

	api.analyze(w, req)

	layerManager.AssertExpectations(t)
}
