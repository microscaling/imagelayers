package registry

import (
	"fmt"
	"net/http"
)

// Authenticator is the interface used to abstract various HTTP authentication
// methods.
//
// ApplyAuthentication accepts an http.Request and will apply the appropriate
// authentication headers to the request.
type Authenticator interface {
	ApplyAuthentication(*http.Request)
}

// NilAuth is an Authenticator implementation that can be used when no
// authentication is required.
type NilAuth struct{}

// No-op implementation. Makes NO changes to the specified request.
func (NilAuth) ApplyAuthentication(r *http.Request) {}

// BasicAuth is an Authenticator implementation which can be used for HTTP
// basic authentication.
type BasicAuth struct {
	Username string
	Password string
}

// Sets the Authorization header on the request to use HTTP Basic
// Authentication with the username and password provided.
func (a BasicAuth) ApplyAuthentication(r *http.Request) {
	r.SetBasicAuth(a.Username, a.Password)
}

// Access is an enum type representing the access type provided by a TokenAuth
// authenticator.
type Access int

// Enum values for the Access type: None, Read, Write, Delete
const (
	None Access = iota
	Read
	Write
	Delete
)

// TokenAuth is an Authenticator implementation which implements the token
// authentication scheme used by the Docker Hub.
type TokenAuth struct {
	Host   string
	Token  string
	Access Access
}

// Sets the Authorization header on the request to use the provided token.
// Will also re-write the host value of the request URL to use the value
// provided in the TokenAuth.
func (a TokenAuth) ApplyAuthentication(r *http.Request) {
	token := fmt.Sprintf("Token %s", a.Token)
	r.Header.Add("Authorization", token)

	r.URL.Host = a.Host
	r.Host = a.Host
}
