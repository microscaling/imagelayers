package registry

import (
	"fmt"
	"net/http"
)

type Authenticator interface {
	ApplyAuthentication(*http.Request)
}

type NilAuth struct{}

func (NilAuth) ApplyAuthentication(r *http.Request) {}

type BasicAuth struct {
	Username string
	Password string
}

func (a BasicAuth) ApplyAuthentication(r *http.Request) {
	r.SetBasicAuth(a.Username, a.Password)
}

type Access int

const (
	None Access = iota
	Read
	Write
	Delete
)

type TokenAuth struct {
	Host   string
	Token  string
	Access Access
}

func (a TokenAuth) ApplyAuthentication(r *http.Request) {
	token := fmt.Sprintf("Token %s", a.Token)
	r.Header.Add("Authorization", token)

	r.URL.Host = a.Host
	r.Host = a.Host
}
