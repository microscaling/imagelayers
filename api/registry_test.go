package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockConnection struct {
	mock.Mock
}

func (m *mockConnection) GetTags(name string) (registry.TagMap, error) {
	args := m.Mock.Called(name)
	return args.Get(0).(registry.TagMap), nil
}

func (m *mockConnection) Search(name string) (*registry.SearchResults, error) {
	args := m.Mock.Called(name)
	return args.Get(0).(*registry.SearchResults), nil
}

func (m *mockConnection) Status() (Status, error) {
	args := m.Mock.Called()
	return args.Get(0).(Status), nil
}

func (m *mockConnection) GetImageLayers(name, tag string) ([]*registry.ImageMetadata, error) {
	args := m.Mock.Called(name, tag)
	return args.Get(0).([]*registry.ImageMetadata), nil
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

	assert.Equal(t, `{"repos":[{"name":"foo","tag":"latest","size":0,"count":0},{"name":"bar","tag":"latest","size":0,"count":0}]}`, string(jsonTest))
}

func TestAnalyzeRequest(t *testing.T) {
	// setup
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)
	inBody := "{\"repos\":[{\"name\":\"foo\",\"tag\":\"latest\"}]}"

	// build request
	req, _ := http.NewRequest("POST", "http://localhost/analyze", strings.NewReader(inBody))
	w := httptest.NewRecorder()

	// build response
	metadata := new(registry.ImageMetadata)
	resp := make([]*registry.ImageMetadata, 1)
	resp[0] = metadata

	// test
	fakeConn.On("GetImageLayers", "foo", "latest").Return([]*registry.ImageMetadata{metadata})
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

func TestSearchRequest(t *testing.T) {
	// setup
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)

	// build request
	req, _ := http.NewRequest("GET", "http://localhost/search?name=foo", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	// test
	results := new(registry.SearchResults)
	fakeConn.On("Search", "foo").Return(results, nil)
	api.handleSearch(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
}

func TestGetTagsRequestWithSlash(t *testing.T) {
	// setup
	image := "centurylink/dray"
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)
	var res registry.TagMap
	fakeConn.On("GetTags", image).Return(res, nil)

	// build request
	req, _ := http.NewRequest("GET", "http://localhost/images/centurylink%2Fdray/tags", strings.NewReader("{}"))
	w := httptest.NewRecorder()
	m := mux.NewRouter()
	m.HandleFunc("/images/{front}/{tail}/tags", api.handleTags).Methods("GET")

	// test
	m.ServeHTTP(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
}

func TestGetTagsRequest(t *testing.T) {
	// setup
	image := "redis"
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)
	var res registry.TagMap
	fakeConn.On("GetTags", image).Return(res, nil)

	// build request
	req, _ := http.NewRequest("GET", "http://localhost/images/redis/tags", strings.NewReader("{}"))
	w := httptest.NewRecorder()
	m := mux.NewRouter()
	m.HandleFunc("/images/{front}/tags", api.handleTags).Methods("GET")

	// test
	m.ServeHTTP(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
}
