package frontend

import (
	"net/http"

	"github.com/akrillis/k8storage/internal/service"
	"github.com/gorilla/mux"
)

const (
	urlFiles = "/files"
)

// NewRoutes returns a new router with all routes.
func NewRoutes(receiver service.Receiver, restorer service.Restorer) *mux.Router {
	m := mux.NewRouter().StrictSlash(true)
	h := &handlers{receiver: receiver, restorer: restorer}

	m.HandleFunc(urlFiles, h.putFile).Methods(http.MethodPost)
	m.HandleFunc(urlFiles, h.getFile).Methods(http.MethodGet)

	return m
}
