package db

import (
	"AvitoTechTask/internal/segment"
	"AvitoTechTask/internal/segment/db"
	"AvitoTechTask/internal/userinsegment"
	postgresql "AvitoTechTask/pkg/client/postgres"
	"AvitoTechTask/pkg/logging"
	"context"
	"github.com/jackc/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r repository) DeleteSegmentFromUsers(ctx context.Context, user *userinsegment.UserInSegment, segment *segment.Segment) error {
	q := `
		UPDATE user_in_segment
		SET out_date = current_timestamp
		WHERE segment_id = $1
		RETURNING out_date
		`
	r.logger.Tracef("SQL query: %s", db.QueryToString(q))
	if err := r.client.QueryRow(ctx, q, segment.SegmentId).Scan(user.OutDate); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}
	return nil
}

func (r repository) DeleteFromSegment(ctx context.Context, user *userinsegment.UserInSegment) error {
	q := `
		UPDATE user_in_segment 
		SET out_date = current_timestamp
        WHERE user_id = $1 AND 
              segment_id = (SELECT segment_id
                              FROM segment
                              WHERE segment_name = $2)
		RETURNING out_date
		`
	r.logger.Tracef("SQL query: %s", db.QueryToString(q))
	if err := r.client.QueryRow(ctx, q, user.UserId, user.SegmentId).Scan(user.OutDate); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}
	return nil
}

func (r repository) AddSegments(ctx context.Context, user *userinsegment.UserInSegment, segmentName string) error {
	q := `
		INSERT INTO user_in_segment
		    (user_id, segment_id, in_date, out_date)
		VALUES ($1, (SElECT segment_id FROM segment WHERE segment_name = $2),
        current_timestamp, DEFAULT)
		RETURNING segment_id, in_date, out_date
		`
	r.logger.Tracef("SQL query: %s", db.QueryToString(q))
	if err := r.client.QueryRow(ctx, q, user.UserId, segmentName).Scan(user.SegmentId, user.InDate, user.OutDate); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}
	return nil
}

func (r repository) AddSegmentsWithPeriod(ctx context.Context, user *userinsegment.UserInSegment, segmentName string) error {
	q := `
		INSERT INTO user_in_segment
		    (user_id, segment_id, in_date, out_date)
		VALUES ($1, (SElECT segment_id FROM segment WHERE segment_name = $2),
        current_timestamp, $3)
		RETURNING segment_id, in_date, out_date
		`
	r.logger.Tracef("SQL query: %s", db.QueryToString(q))
	if err := r.client.QueryRow(ctx, q, user.UserId, segmentName, user.OutDate).Scan(user.SegmentId, user.InDate, user.OutDate); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}
	return nil
}

func (r repository) GetSegments(ctx context.Context, user *userinsegment.UserInSegment) (segments []SegmentDTO, err error) {
	q := `
		SELECT segment_name 
		FROM (SELECT segment_id AS s_id
		      FROM user_in_segment
			  WHERE user_id = 1000
			    AND (out_date IS NULL
			             OR out_date < current_date)
			  ) AS active_u_segments
		    INNER JOIN segment ON s_id = segment.segment_id
		`
	r.logger.Tracef("SQL query: %s", db.QueryToString(q))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {

	}
}

func NewDataBase(l *logging.Logger, c postgresql.Client) userinsegment.Repository {
	return &repository{
		logger: l,
		client: c,
	}
}
