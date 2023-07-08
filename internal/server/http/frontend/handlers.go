package frontend

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akrillis/k8storage/internal/entities"
	"github.com/akrillis/k8storage/internal/service"
	"github.com/go-playground/validator/v10"
)

const (
	errParse    = "failed to parse request"
	errNotFound = "not found"

	pClientID = "client_id"
	pName     = "name"
)

type handlers struct {
	receiver service.Receiver
	restorer service.Restorer
}

// putFile handler is responsible for putting file to storage.
func (h *handlers) putFile(w http.ResponseWriter, r *http.Request) {
	var req entities.PutFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("[putFile] %s: %s", errParse, err.Error())
		return
	}

	if err := validator.New().Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("[putFile] %s: %s", errParse, err.Error())
		return
	}

	if err := h.receiver.Put(r.Context(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("[putFile] failed to put file: %s", err.Error())
		return
	}

	writeOk(w, nil)
	log.Printf("clientID %s file %s has been put", req.ClientID, req.Name)
}

// getFile handler is responsible for getting file from storage.
func (h *handlers) getFile(w http.ResponseWriter, r *http.Request) {
	req := &entities.GetFileRequest{
		ClientID: r.URL.Query().Get(pClientID),
		Name:     r.URL.Query().Get(pName),
	}

	if err := validator.New().Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("[getFile] %s: %s", errParse, err.Error())
		return
	}

	file, err := h.restorer.Get(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("[getFile] failed to get file: %s", err.Error())
		return
	}

	if len(file) == 0 {
		http.Error(w, errNotFound, http.StatusNoContent)
		log.Printf("[getFile] file %s not found", req.Name)
		return
	}

	writeOk(
		w,
		entities.GetFileResponse{
			Name: req.Name,
			Data: file,
		},
	)
	log.Printf("clientID %s file %s has been got", req.ClientID, req.Name)
}

func writeOk(w http.ResponseWriter, body interface{}) {
	answerBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("jsonResponse body: %v error: %s", body, err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(answerBody); err != nil {
		log.Printf("write error: %s", err.Error())
	}
}
