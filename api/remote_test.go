package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteConnectionInterface(t *testing.T) {
	assert.Implements(t, (*RegistryConnection)(nil), new(remoteConnection))
}

func TestStatus(t *testing.T) {
	rc := new(remoteConnection)
	s, _ := rc.Status()

	assert.Equal(t, s.Message, "Connected to Registry")
	assert.Equal(t, s.Service, "Registry Image Manager")
}
