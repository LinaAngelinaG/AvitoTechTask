package userinsegment

import (
	"AvitoTechTask/internal/handlers"
	"AvitoTechTask/internal/segment"
	"AvitoTechTask/pkg/logging"
	"context"
	"encoding/json"
	"github.com/gocarina/gocsv"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

var _ handlers.Handler = &handler{}

const (
	userURL         = "/user/:uid"
	historyURL      = "/history/:uid/:year/:month"
	userSegmentsURL = "/user/segments"
	dateFormat      = "2006-01-02 15:04:05"
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
	router.GET(historyURL, h.GetUserHistory)
	router.DELETE(userSegmentsURL, h.DeleteListOfSegments)
	router.POST(userSegmentsURL, h.AddUserSegments)
}

// @Summary GetListOfSegments
// @Tags user_in_segment
// @Description Provide JSON with values {"userId": someId, "segments":['seg_name_1','seg_name_2',...]]}, that contains pare of user_id and all active segments of this user.
// @ID get-list-of-user-segments
// @Accept  json
// @Produce  json
// @Param uid path string true "USER_ID"
// @Success 200 {object} UserSegmentsList
// @Failure 400 {string} string "there is no user_id in context or its wrong value"
// @Failure 418 {string} string "error with getting list of segments"
// @Failure 405 {string} string "error with marshalling list of segments"
// @Router /user/:uid [get]
func (h *handler) GetListOfSegments(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	userId, err := strconv.Atoi(params.ByName("uid"))
	if err != nil || userId <= 0 {
		h.logger.Error("there is no user_id in context or its wrong value")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("there is no user_id in context or its wrong value"))
	}
	h.logger.Infof("geting list of segment for user with id %d started", userId)
	u := UserInSegment{UserId: userId}
	segs, err := h.repo.GetSegments(r.Context(), &u)
	if err != nil {
		w.WriteHeader(http.StatusTeapot)
		h.logger.Error("error with getting list of segments")
		w.Write([]byte("error with getting list of segments"))
	}
	segsBytes, err := json.Marshal(segs)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		h.logger.Error("error with marshalling list of segments")
		w.Write([]byte("error with marshalling list of segments"))
	}
	h.logger.Info("got list of segments")
	w.WriteHeader(200)
	w.Write(segsBytes)
}

// @Summary DeleteListOfSegments
// @Tags user_in_segment
// @Description Iterates throuth the list of sent segment_names and delete each of row that have such user_id value and segment_id value. Not literally delete, but update value 'out_date' to current_date if it's null.
// @ID delete-user_in_segment
// @Accept  json
// @Produce  json
// @Param segments body UserSegmentsList true "SEGMENTS LIST FOR USER"
// @Success 200 {string}  string "user deleted from segments"
// @Failure 400 {string}  string "Invalid request body"
// @Failure 400 {string}  string "Something wrong with deleting segment"
// @Router /user/segments [delete]
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

// @Summary AddUserSegments
// @Tags user_in_segment
// @Description Create new user_in_segment sequence with provided data, only if provided segment exists. If body contains 'period' value - creating user_in_segment sequence with out_date = current_date + period.
// @ID add-segments-to-user
// @Accept  json
// @Produce  json
// @Param segments body UserSegmentsList true "SEGMENTS LIST FOR USER"
// @Success 201 {string}  string "segments added to user"
// @Failure 400 {string}  string "Invalid request body"
// @Failure 400 {string}  string "Something wrong with addding segment"
// @Router /user/segments [post]
func (h *handler) AddUserSegments(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

// @Summary GetUserHistory
// @Tags user_in_segment
// @Description Returns the link to download the CSV-file with history of all user's activities in segments: adding or deleting in period of seted month and the year.
// @ID get-history-of-user
// @Accept  json
// @Produce  json
// @Param uid path string true "USER_ID"
// @Param year path int true "YEAR"
// @Param month path int true "MONTH"
// @Success 200 {array} UserInSegmentsHistory "got list of segments"
// @Failure 400 {string} string "there is no 'parameter_name' in context or it's wrong value"
// @Failure 418 {string} string "error with getting history"
// @Router /history/:uid/:year/:month [get]
func (h *handler) GetUserHistory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userId, err := strconv.Atoi(params.ByName("uid"))
	if err != nil || userId <= 0 {
		h.logger.Error("there is no user_id in context or its wrong value")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("there is no user_id in context or its wrong value"))
		return
	}
	year, err := strconv.Atoi(params.ByName("year"))
	curY := time.Now().Year()
	if err != nil || year <= 0 || year > curY {
		h.logger.Error("there is no year in context or its wrong value")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("there is no year in context or its wrong value"))
		return
	}
	month, err := strconv.Atoi(params.ByName("month"))
	if err != nil || month <= 0 || month > 12 {
		h.logger.Error("there is no month in context or its wrong value")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("there is no month in context or its wrong value"))
		return
	}
	date := formDate(strconv.Itoa(year), month)
	if !relevantDate(date) {
		h.logger.Error("date is not relevant")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("date is not relevant"))
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	history, err := h.repo.GetUserHistory(r.Context(), &UserInSegment{
		UserId: userId,
		InDate: date,
	})
	if err != nil {
		w.WriteHeader(http.StatusTeapot)
		h.logger.Errorf("error with getting history for user %d for year:%d for month:%d; error: %s", userId, year, month, err.Error())
		w.Write([]byte("error with getting history"))
	}
	h.logger.Info("got list of segments")
	w.WriteHeader(200)
	gocsv.Marshal(history, w)
}

func formDate(y string, m int) string {
	if m <= 10 {
		return y + "-0" + strconv.Itoa(m) + "-01 00:00:00"
	}
	return y + "-" + strconv.Itoa(m) + "-01 00:00:00"

}

func relevantDate(date string) bool {
	d, err := time.Parse(dateFormat, date)
	if err != nil || d.Format(dateFormat) != date {
		return false
	}
	return true
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
