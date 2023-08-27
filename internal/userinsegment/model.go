package userinsegment

import "time"

type UserInSegment struct {
	UserId    int       `json:"user_id"`
	SegmentId int       `json:"segment_id"`
	InDate    time.Time `json:"in_date"`
	OutDate   time.Time `json:"out_date"`
}
