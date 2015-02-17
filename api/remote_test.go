package api

import (
//	"encoding/json"
//	"net/http"
//	"net/http/httptest"
	"testing"
//	"net/url"
//
//	"github.com/CenturyLinkLabs/docker-reg-client/registry"
	"github.com/stretchr/testify/mock"
//	"github.com/stretchr/testify/assert"
)

type mockHub struct {
	*mock.Mock
}

type mockRepository struct {
	*mock.Mock
}

type mockImage struct {
	*mock.Mock
}

type mockClient struct {
	Hub *mockHub
	Image *mockImage
	Repository *mockRepository
}

var (
	client *mockClient
)

func clientSetup() {
	client = new(mockClient)

	client.Hub = new(mockHub)
	client.Image = new(mockImage)
	client.Repository = new(mockRepository)
}


func TestAnalyze(t *testing.T) {
	clientSetup()

}

func TestStatus(t *testing.T) {

}
