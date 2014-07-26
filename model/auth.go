package model

import (
	"net/http"

	"../cfg"
)

type Auth struct {
	AuthName		string
	AuthPassword	string
}

func (auth *Auth) Bind(req *http.Request) error {
	return bindFormToStruct(req, auth);
}

func (auth *Auth) IsValid() bool {
	return auth.AuthName == cfg.String("auth.name") && auth.AuthPassword == cfg.String("auth.password")
}
