package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockConnection struct {
	mock.Mock
}

func (m *mockConnection) Status() (Status, error) {
	args := m.Mock.Called()
	return args.Get(0).(Status), nil
}

func (m *mockConnection) Connect(repo string) error {
	args := m.Mock.Called(repo)
	return args.Error(0)
}

func (m *mockConnection) GetImageID(image string, tag string) (string, error) {
	args := m.Mock.Called(image, tag)
	return args.Get(0).(string), nil
}

func (m *mockConnection) GetAncestry(id string) ([]string, error) {
	args := m.Mock.Called(id)
	return args.Get(0).([]string), nil
}

func (m *mockConnection) GetMetadata(layer string) (*registry.ImageMetadata, error) {
	args := m.Mock.Called(layer)
	return args.Get(0).(*registry.ImageMetadata), nil
}

func TestMarshalStatus(t *testing.T) {
	status := Status{Message: "active", Service: "foo"}
	jsonTest, _ := json.Marshal(status)

	assert.Equal(t, `{"message":"active","service":"foo"}`, string(jsonTest))
}

func TestMarshalRequest(t *testing.T) {
	repo1 := Repo{Name: "foo", Tag: "latest"}
	repo2 := Repo{Name: "bar", Tag: "latest"}
	imageList := []Repo{repo1, repo2}
	repos := Request{Repos: imageList}
	jsonTest, _ := json.Marshal(repos)

	assert.Equal(t, `{"repos":[{"name":"foo","tag":"latest"},{"name":"bar","tag":"latest"}]}`, string(jsonTest))
}

func TestAnalyzeRequest(t *testing.T) {
	// setup
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)
	layers := []string{"baz"}
	inBody, image := "{\"repos\":[{\"name\":\"foo\",\"tag\":\"latest\"}]}", "foo"

	// build request
	req, _ := http.NewRequest("POST", "http://localhost/analyze", strings.NewReader(inBody))
	w := httptest.NewRecorder()

	// build response
	metadata := new(registry.ImageMetadata)
	resp := make([]*registry.ImageMetadata, 1)
	resp[0] = metadata

	// test
	fakeConn.On("Connect", image).Return(nil)
	fakeConn.On("GetImageID", image, "latest").Return("fooID")
	fakeConn.On("GetAncestry", "fooID").Return(layers)
	fakeConn.On("GetMetadata", "baz").Return(metadata)
	api.handleAnalysis(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
}

func TestStatusRequest(t *testing.T) {
	// setup
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)

	// build request
	resp := Status{Message: "foo", Service: "bar"}
	req, _ := http.NewRequest("GET", "http://localhost/analyze", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	// test
	fakeConn.On("Status").Return(resp, nil)
	api.handleStatus(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
}
