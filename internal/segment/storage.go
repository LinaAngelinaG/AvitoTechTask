package segment

import "context"

type Repository interface {
	Create(ctx context.Context, segment *Segment) error
	Delete(ctx context.Context, segment *Segment) error
}

//Create::
//INSERT INTO segment(segment_id, segment_name, active) VALUES (DEFAULT, $1, DEFAULT)

//Delete::
//UPDATE segment SET (active) = false WHERE segment_name = $1
//UPDATE userinsegment SET (out_date) = current_timestamp WHERE segment_id =
//(SELECT segment_id FROM segment WHERE segment_name = $1)
