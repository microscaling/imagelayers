package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus (t *testing.T) {
	reg := NewRemoteRegistry()
	s, _ := reg.Status()

	assert.Equal(t, s.Message, "Connected to Registry")
	assert.Equal(t, s.Service, "Connected to Registry")

}
