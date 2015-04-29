package api

import (
	"encoding/json"
	"errors"
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
	return args.Get(0).(registry.TagMap), args.Error(1)
}

func (m *mockConnection) Search(name string) (*registry.SearchResults, error) {
	args := m.Mock.Called(name)
	var r *registry.SearchResults
	if a := args.Get(0); a != nil {
		r = a.(*registry.SearchResults)
	}
	return r, args.Error(1)
}

func (m *mockConnection) Status() (Status, error) {
	args := m.Mock.Called()
	return args.Get(0).(Status), args.Error(1)
}

func (m *mockConnection) GetImageLayers(name, tag string) ([]*registry.ImageMetadata, error) {
	args := m.Mock.Called(name, tag)
	return args.Get(0).([]*registry.ImageMetadata), args.Error(1)
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
	inBody := `{"repos":[{"name":"foo","tag":"latest"}]}`

	// build request
	req, _ := http.NewRequest("POST", "http://localhost/analyze", strings.NewReader(inBody))
	w := httptest.NewRecorder()

	// build response
	metadata := new(registry.ImageMetadata)
	resp := make([]*registry.ImageMetadata, 1)
	resp[0] = metadata

	// test
	fakeConn.On("GetImageLayers", "foo", "latest").Return([]*registry.ImageMetadata{metadata}, nil)
	api.handleAnalysis(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnalyzeRequestErroredBadParameters(t *testing.T) {
	// setup
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)
	inBody := "HAIL AND WELL MET!!! I AM BAD JSON!!1!!"

	// build request
	req, _ := http.NewRequest("POST", "http://localhost/analyze", strings.NewReader(inBody))
	w := httptest.NewRecorder()

	// test
	api.handleAnalysis(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")
}

func TestAnalyzeRequestWithRegistryError(t *testing.T) {
	// setup
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)
	inBody := `{"repos":[{"name":"foo","tag":"latest"}]}`

	// build request
	req, _ := http.NewRequest("POST", "http://localhost/analyze", strings.NewReader(inBody))
	w := httptest.NewRecorder()

	// test
	err := errors.New("test error")
	fakeConn.On("GetImageLayers", "foo", "latest").Return([]*registry.ImageMetadata{}, err)
	api.handleAnalysis(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
	if assert.NotNil(t, w.Body) {
		var r []*Response
		assert.NoError(t, json.Unmarshal([]byte(w.Body.String()), &r))

		if assert.Len(t, r, 1) {
			assert.Equal(t, http.StatusInternalServerError, r[0].Status)
		}
	}
}

func TestStatusRequest(t *testing.T) {
	// setup
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)

	// build request
	resp := Status{Message: "foo", Service: "bar"}
	req, _ := http.NewRequest("GET", "http://localhost/status", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	// test
	fakeConn.On("Status").Return(resp, nil)
	api.handleStatus(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"foo","service":"bar"}`, strings.TrimSpace(w.Body.String()))
}

func TestStatusRequestWithRegistryError(t *testing.T) {
	// setup
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)

	// build request
	req, _ := http.NewRequest("GET", "http://localhost/status", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	// test
	fakeConn.On("Status").Return(Status{}, errors.New("test error"))
	api.handleStatus(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "test error", w.Body.String())
}

func TestSearchRequest(t *testing.T) {
	// setup
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)

	// build request
	req, _ := http.NewRequest("GET", "http://localhost/search?name=foo", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	// test
	results := registry.SearchResults{Query: "Test Query"}
	fakeConn.On("Search", "foo").Return(&results, nil)
	api.handleSearch(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Query")
}

func TestSearchRequestWithRegistryError(t *testing.T) {
	// setup
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)

	// build request
	req, _ := http.NewRequest("GET", "http://localhost/search?name=foo", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	// test
	fakeConn.On("Search", "foo").Return(nil, errors.New("test error"))
	api.handleSearch(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "test error", w.Body.String())
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
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTagsRequest(t *testing.T) {
	// setup
	image := "redis"
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)
	res := registry.TagMap{"Key": "Value"}
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
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"Key":"Value"}`, strings.TrimSpace(w.Body.String()))
}

func TestGetTagsRequestWithRegistryError(t *testing.T) {
	// setup
	image := "centurylink/dray"
	fakeConn := new(mockConnection)
	api := newRegistryApi(fakeConn)
	fakeConn.On("GetTags", image).Return(registry.TagMap{}, errors.New("test error"))

	// build request
	req, _ := http.NewRequest("GET", "http://localhost/images/centurylink%2Fdray/tags", strings.NewReader("{}"))
	w := httptest.NewRecorder()
	m := mux.NewRouter()
	m.HandleFunc("/images/{front}/{tail}/tags", api.handleTags).Methods("GET")

	// test
	m.ServeHTTP(w, req)

	// asserts
	fakeConn.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "test error", w.Body.String())
}
