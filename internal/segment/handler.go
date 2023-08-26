package segment

import (
	"AvitoTechTask/internal/handlers"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	segmentURL = "/segment/:name"
)

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *httprouter.Router) {
	router.DELETE(segmentURL, h.DeleteSegment)
	router.POST(segmentURL, h.CreateSegment)
}

func (h *handler) DeleteSegment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("segment deleted"))
}

func (h *handler) CreateSegment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(201)
	w.Write([]byte("segment created"))
}
