package metric

import (
	"net/http"
)

const (
	URL = "/api/heartbeat"
)

type Handler struct {
}

type HandlerFunc interface {
	HandlerFunc(method, path string, handler http.HandlerFunc)
}

func (h *Handler) Register(router HandlerFunc) {
	router.HandlerFunc(http.MethodGet, URL, h.Heartbeat)
}

func (h *Handler) Heartbeat(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Heartbeat"))
}
