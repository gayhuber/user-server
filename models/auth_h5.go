package models

import (
	"user-server/lib"
)

type h5Auth struct {
	Src string
	Ext map[string]interface{}
}

func (auth *h5Auth) getName() string {
	return "h5Auth"
}

func (auth *h5Auth) register() (code int, obj interface{}) {
	return 200, lib.H{
		"result": "this is h5 handler",
		"src":    auth.Src,
		"other":  auth.Ext,
	}
}

func (auth *h5Auth) login() (code int, obj interface{}) {
	return
}

func (auth *h5Auth) info() (code int, obj interface{}) {
	return
}

func (auth *h5Auth) home() (code int, obj interface{}) {
	return
}

func (auth *h5Auth) setParams(params map[string]interface{}) {
	auth.Src = params["src"].(string)
	if auth.Src != "fanli" {
		auth.Ext = params
	}

	return
}
