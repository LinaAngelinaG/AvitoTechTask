package userinsegment

import (
	"AvitoTechTask/internal/segment"
	"context"
)

type Repository interface {
	AddSegments(ctx context.Context, user *UserInSegment) error
	DeleteSegments(ctx context.Context, user *UserInSegment) error
	GetSegments(ctx context.Context, user *UserInSegment) ([]segment.Segment, error)
}

//GetSegments::
//SELECT segment_name FROM (SELECT segment_id AS segment_id FROM userinsegment
//WHERE user_id = $1 && out_date IS NULL OR out_date < current_date) AS segments
//INNER JOIN segment ON segments.segment_id = segment.segment_id
//
//AddSegments::
//INSERT INTO segment(segment_id, segment_name, active) VALUES (DEFAULT, $1, DEFAULT)
//
//DeleteSegments::
// UPDATE userinsegment SET (out_date) = current_timestamp WHERE segment_id =
//(SELECT segment_id FROM segment WHERE segment_name = $1) AND user_id = $2
//
//UPDATE segment SET (active) =
// (SELECT active FROM segment WHERE segment_name = $1)
//
