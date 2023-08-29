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
	repo   Repository
}

func NewHandler(l *logging.Logger, r Repository) handlers.Handler {
	l.Info("register segment handler")
	return &handler{
		logger: l,
		repo:   r,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.DELETE(segmentURL, h.DeleteSegment)
	router.POST(segmentURL, h.CreateSegment)
}

func (h *handler) DeleteSegment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	if name == "" {
		h.logger.Error("there is no segment_name in context")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("there is no segment_name in context"))
	}
	h.logger.Infof("segment segment_name: %s is deleting", name)
	seg := Segment{SegmentName: name}
	err := h.repo.Delete(r.Context(), &seg)
	if err != nil {
		h.logger.Errorf("%s : something wrong with deleting user from segment", err.Error())
		http.Error(w, "Something wrong with deleting segment", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("segment deleted from users"))
}

func (h *handler) CreateSegment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	if name == "" {
		h.logger.Error("there is no segment_name in context")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("there is no segment_name in context"))
	}
	h.logger.Infof("there is segment_name: %s", name)
	seg := Segment{SegmentName: name}
	err := h.repo.Create(r.Context(), &seg)
	if err != nil {
		h.logger.Errorf("error while creating segment: %s", err.Error())
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("error with creating entity example"))
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("segment created"))
}
