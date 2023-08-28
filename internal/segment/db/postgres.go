package db

import (
	"AvitoTechTask/internal/segment"
	postgresql "AvitoTechTask/pkg/client/postgres"
	"AvitoTechTask/pkg/logging"
	"context"
	"github.com/jackc/pgconn"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, segment *segment.Segment) error {
	q := `
		INSERT INTO segment 
		    (segment_name)
		VALUES ($1) 
		RETURNING segment_id, active
		`
	r.logger.Tracef("SQL query: %s", QueryToString(q))
	if err := r.client.QueryRow(ctx, q, segment.SegmentName).Scan(segment.SegmentId, segment.Active); err != nil {
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
	if err := r.client.QueryRow(ctx, q, segment.SegmentName).Scan(segment.SegmentId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}
	return nil
}

func NewRepository(l *logging.Logger, c postgresql.Client) segment.Repository {
	return &repository{
		client: c,
		logger: l,
	}
}
