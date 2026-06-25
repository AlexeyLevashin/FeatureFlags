package handlers

import "net/http"

type FlagHandler struct {
	// service service.FlagService
}

func (h *FlagHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
