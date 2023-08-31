package segmentdb

import (
	"AvitoTechTask/internal/segment"
	postgresql "AvitoTechTask/pkg/client/postgres"
	"AvitoTechTask/pkg/logging"
	"context"
	"github.com/jackc/pgconn"
	"log"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(l *logging.Logger, c postgresql.Client) segment.Repository {
	return &repository{
		client: c,
		logger: l,
	}
}

func (r *repository) Create(ctx context.Context, segment *segment.Segment) error {
	q := `
		INSERT INTO segment 
		    (segment_name)
		VALUES ($1) 
		`
	r.logger.Tracef("SQL query: %s", QueryToString(q))
	if _, err := r.client.Query(ctx, q, segment.SegmentName); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}
	return nil
}

func QueryToString(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", "")
}

func (r *repository) Delete(ctx context.Context, segment *segment.Segment) error {
	q := `
		UPDATE segment 
		SET active = false 
		WHERE segment_name = $1
		RETURNING segment_id
		`
	r.logger.Tracef("SQL query: %s", QueryToString(q))
	if err := r.client.QueryRow(ctx, q, segment.SegmentName).Scan(&segment.SegmentId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}
	return r.deleteFromUsers(ctx, segment)
}

func (r *repository) deleteFromUsers(ctx context.Context, segment *segment.Segment) error {
	q := `
		UPDATE user_in_segment
		SET out_date = current_timestamp
		WHERE segment_id = $1 AND out_date is null
		`
	r.logger.Tracef("SQL query: %s", QueryToString(q))
	if _, err := r.client.Query(ctx, q, segment.SegmentId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		log.Println(err.Error())

		return err
	}
	return nil
}
