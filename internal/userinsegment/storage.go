package userinsegment

import (
	"AvitoTechTask/internal/segment"
	"context"
)

type Repository interface {
	AddSegments(ctx context.Context, user *UserDTO, segment *segment.SegmentDTO) error
	AddSegmentsWithPeriod(ctx context.Context, user *UserDTO, segment *segment.SegmentDTO) error
	DeleteFromSegment(ctx context.Context, user *UserInSegmentDTO) error
	GetSegments(ctx context.Context, user *UserInSegment) (UserSegmentsList, error)
	GetUserHistory(ctx context.Context, user *UserDTO) ([]UserInSegmentsHistory, error)
}

//GetSegments::
//SELECT segment_name FROM (SELECT segment_id AS segment_id FROM user_in_segment
//WHERE user_id = $1 && out_date IS NULL OR out_date < current_date) AS segments
//INNER JOIN segment ON segments.segment_id = segment.segment_id
//
//AddSegments::
//INSERT INTO user_in_segment(user_id, segment_id, in_date, out_date)
//VALUES ($1, (SElECT segment_id FROM segment WHERE segment_name = $2), current_timestamp, DEFAULT)
//
//INSERT INTO user_in_segment(user_id, segment_id, in_date, out_date)
//VALUES ($1, (SElECT segment_id FROM segment WHERE segment_name = $2), current_timestamp, $3)
//
//DeleteSegments::
//UPDATE user_in_segment SET out_date = current_timestamp WHERE segment_id =
//(SELECT segment_id FROM segment WHERE segment_name = $1) AND user_id = $2
//
//
//Get history for user:
//SELECT * from user_in_segment WHERE user_id = $1
