package userinsegment

import (
	"AvitoTechTask/internal/handlers"
	"AvitoTechTask/internal/segment"
	"AvitoTechTask/pkg/logging"
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

var _ handlers.Handler = &handler{}

const (
	userURL         = "/user/:uid"
	userSegmentsURL = "/user/segments"
)

type handler struct {
	logger *logging.Logger
	repo   Repository
}

func NewHandler(l *logging.Logger, r Repository) handlers.Handler {
	l.Info("register user handler")
	return &handler{
		logger: l,
		repo:   r,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	h.logger.Info("register user handler methods")
	router.GET(userURL, h.GetListOfSegments)
	router.DELETE(userSegmentsURL, h.DeleteListOfSegments)
	router.POST(userSegmentsURL, h.AddUserSegments)
}

func (h *handler) GetListOfSegments(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	userId, err := strconv.Atoi(params.ByName("uid"))
	if err != nil || userId <= 0 {
		h.logger.Error("there is no user_id in context or its wrong value")
		w.WriteHeader(418)
		w.Write([]byte("there is no user_id in context or its wrong value"))
	}
	h.logger.Infof("geting list of segment for user with id %d started", userId)
	u := UserInSegment{UserId: userId}
	segs, err := h.repo.GetSegments(r.Context(), &u)
	if err != nil {
		w.WriteHeader(418)
		h.logger.Error("error with getting list of segments")
		w.Write([]byte("error with getting list of segments"))
	}
	segsBytes, err := json.Marshal(segs)
	if err != nil {
		w.WriteHeader(418)
		h.logger.Error("error with marshalling list of segments")
		w.Write([]byte("error with marshalling list of segments"))
	}
	h.logger.Error("got list of segments")
	w.WriteHeader(200)
	w.Write(segsBytes)
}

func (h *handler) DeleteListOfSegments(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var segments UserSegmentsList
	err := decoder.Decode(&segments)

	if err != nil {
		h.logger.Error("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	} else {
		h.logger.Info("start delete-user-segments")
		names := segments.SegmentNames
		for _, segName := range names {
			err = h.repo.DeleteFromSegment(r.Context(), &UserInSegmentDTO{
				UserId:      segments.UserId,
				SegmentName: segName,
			})
			if err != nil {
				h.logger.Errorf("%s : something wrong with deleting user from segment", err.Error())
				http.Error(w, "Something wrong with deleting segment", http.StatusBadRequest)
				return
			}
		}

	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user deleted from segments"))
}

func (h *handler) AddUserSegments(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var segments UserSegmentsList
	err := decoder.Decode(&segments)

	if err != nil {
		h.logger.Error("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	} else if segments.Period != 0 {
		h.logger.Info("start add-user-segments with period")
		err = addSegments(&UserDTO{UserId: segments.UserId, Period: segments.Period},
			segments.SegmentNames,
			r.Context(),
			h.repo.AddSegmentsWithPeriod)
	} else {
		h.logger.Info("start add-user-segments without period")
		err = addSegments(&UserDTO{UserId: segments.UserId},
			segments.SegmentNames,
			r.Context(),
			h.repo.AddSegments)
	}
	if err != nil {
		h.logger.Errorf("%s : something wrong with addding segment to user ", err.Error())
		http.Error(w, "Something wrong with addding segment", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("segments added to user"))
}

func addSegments(u *UserDTO, names []string, context context.Context, repoAdd func(ctx context.Context, user *UserDTO, segment *segment.SegmentDTO) error) (err error) {
	for _, segName := range names {
		err = repoAdd(context, u, &segment.SegmentDTO{SegmentName: segName})
		if err != nil {
			return err
		}
	}
	return err
}
