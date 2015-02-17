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

func TestMarshalStatus (t *testing.T) {
	status := Status{Message: "active", Service: "foo"}
	jsonTest, _ := json.Marshal(status)

	assert.Equal(t, `{"message":"active","service":"foo"}`, string(jsonTest))
}

func TestMarshalRequest (t *testing.T) {
	imageList := []string{"foo", "bar", "baz"}
	repos := Request{Repos: imageList}
	jsonTest, _ := json.Marshal(repos)

	assert.Equal(t, `{"repos":["foo","bar","baz"]}`, string(jsonTest))
}

func TestRoutes (t *testing.T) {
	layerManager := new(mockLayerManager)
	api := newRegistryApi(layerManager)
	routes := api.Routes()

	assert.NotNil(t, routes["GET"])
	assert.NotNil(t, routes["GET"]["/status"])
	assert.NotNil(t, routes["POST"])
	assert.NotNil(t, routes["POST"]["/analyze"])
}

func TestAnalyzeRequest(t *testing.T) {
	// setup
	layerManager := new(mockLayerManager)
	api := newRegistryApi(layerManager)
	inBody, images := "{\"repos\":[\"foo\"]}", []string{"foo"}

	// build request
	resp := make([]*registry.ImageMetadata,1)
	req, _ := http.NewRequest("POST", "http://localhost/analyze", strings.NewReader(inBody))
	w := httptest.NewRecorder()

	// test
	layerManager.On("Analyze", images).Return(resp, nil)
	api.analyze(w, req)

	// asserts
	layerManager.AssertExpectations(t)
}

func TestStatusRequest(t *testing.T) {
	// setup
	layerManager := new(mockLayerManager)
	api := newRegistryApi(layerManager)

	// build request
	resp := Status{Message:"foo", Service:"bar"}
	req, _ := http.NewRequest("GET", "http://localhost/analyze", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	// test
	layerManager.On("Status").Return(resp, nil)
	api.status(w, req)

	// asserts
	layerManager.AssertExpectations(t)
}
