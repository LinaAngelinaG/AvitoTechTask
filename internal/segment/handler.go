package segment

import (
	"AvitoTechTask/internal/handlers"
	"AvitoTechTask/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	segmentURL = "/segment/:name"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(l *logging.Logger) handlers.Handler {
	l.Info("register segment handler")
	return &handler{
		logger: l,
	}
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
