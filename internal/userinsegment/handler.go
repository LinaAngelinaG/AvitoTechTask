package userinsegment

import (
	"AvitoTechTask/internal/handlers"
	"AvitoTechTask/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var _ handlers.Handler = &handler{}

const (
	userURL         = "/user/:uid"
	userSegmentsURL = "/user/segments/:uid"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(l *logging.Logger) handlers.Handler {
	l.Info("register user handler")
	return &handler{
		logger: l,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(userURL, h.GetListOfSegments)
	router.DELETE(userSegmentsURL, h.DeleteListOfSegments)
	router.POST(userSegmentsURL, h.AddUserSegments)
}

func (h *handler) GetListOfSegments(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("this is the list of user segments"))
}

func (h *handler) DeleteListOfSegments(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("user deleted from segments"))
}

func (h *handler) AddUserSegments(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(201)
	w.Write([]byte("segments added to user"))
}
