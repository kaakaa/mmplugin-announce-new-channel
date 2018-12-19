package main

import (
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/plugin"
)

const (
	iconFileName = "default_icon.png"
)

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	router.HandleFunc("/"+iconFileName, p.handleIcon).Methods("GET")

	router.ServeHTTP(w, r)
}

func (p *Plugin) handleIcon(w http.ResponseWriter, r *http.Request) {
	box := packr.NewBox("./assets")

	b, err := box.Find(iconFileName)
	if err != nil {
		p.API.LogError("Failed to load icon file", "details", err)
	}
	if _, err := w.Write(b); err != nil {
		p.API.LogError("Failed to write response", "details", err)
	}
}
