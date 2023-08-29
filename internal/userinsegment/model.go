package userinsegment

import "time"

type UserInSegment struct {
	UserId    int       `json:"user_id"`
	SegmentId int       `json:"segment_id"`
	InDate    time.Time `json:"in_date"`
	OutDate   time.Time `json:"out_date"`
}

type UserInSegmentDTO struct {
	UserId      int    `json:"user_id"`
	SegmentName string `json:"segment_name"`
}

type UserDTO struct {
	UserId int `json:"user_id"`
	Period int `json:"period"`
}

type UserSegmentsList struct {
	UserId       int      `json:"user_id"`
	SegmentNames []string `json:"segment_list"`
	Period       int      `json:"period"`
}

type UserInSegmentsHistory struct {
	UserId      int       `json:"user_id"`
	SegmentName string    `json:"segment_name"`
	Event       string    `json:"event"`
	EventDate   time.Time `json:"event_date"`
}
