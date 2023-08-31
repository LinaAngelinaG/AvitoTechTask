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

// @Summary DeleteSegment
// @Tags segment
// @Description Delete segment with its name that is sent as parameter of http-request. Don't phisically delete segment, just set its parameter 'active' from value 'true' to 'false'. Also delete users from that segment: change null-values of 'out_date' in table user_in_segment to current_date.
// @ID delete-segment
// @Accept  json
// @Produce  json
// @Param name path string true "SEGMENT NAME"
// @Success 200 {string}  string "segment deleted from users"
// @Failure 400 {string}  string "there is no segment_name in context"
// @Failure 400 {string}  string "Something wrong with deleting segment"
// @Router /segment/:name [delete]
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

// @Summary CreateSegment
// @Tags segment
// @Description Create segment with its name that is sent as parameter of http-request. 'segment_id' is autoinremented in DB while inserting, 'active' set to default value 'true'.
// @ID delete-segment
// @Accept  json
// @Produce  json
// @Param name path string true "SEGMENT NAME"
// @Success 201 {string}  string "segment created"
// @Failure 400 {string}  string "there is no segment_name in context"
// @Failure 418 {string}  string "error with creating entity example"
// @Router /segment/:name [post]
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
