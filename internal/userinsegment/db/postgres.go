package userinsegmentdb

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

func (r repository) GetUserHistory(ctx context.Context, user *userinsegment.UserInSegment) (history []userinsegment.UserInSegmentsHistory, err error) {
	q := `
		SELECT TO_CHAR(in_date, 'YYYY/MM/DD HH24:MI:SS'),
		       (SELECT segment_name
		        FROM segment
		        WHERE user_in_segment.segment_id = segment.segment_id)
		FROM user_in_segment
		WHERE user_id = $1 
		  AND in_date >= $2
		  AND in_date < TO_TIMESTAMP(cast($2 as TEXT),'YYYY/MM/DD HH24:MI:SS') + interval '1 month'
		`
	history = make([]userinsegment.UserInSegmentsHistory, 0)

	if err = r.addToHistory(q, "inserted", user, ctx, &history); err != nil {
		return history, err
	}
	q = `
		SELECT TO_CHAR(out_date, 'YYYY/MM/DD HH24:MI:SS'),
		       (SELECT segment_name
		        FROM segment
		        WHERE user_in_segment.segment_id = segment.segment_id)
		FROM user_in_segment
		WHERE user_id = $1 
		  AND out_date >= $2 
		  AND out_date < TO_TIMESTAMP(cast($2 as TEXT),'YYYY/MM/DD HH24:MI:SS') + interval '1 month'
		`
	if err = r.addToHistory(q, "deleted", user, ctx, &history); err != nil {
		return history, err
	}

	return history, nil
}

func (r repository) addToHistory(q string, event string, user *userinsegment.UserInSegment, ctx context.Context, history *[]userinsegment.UserInSegmentsHistory) error {
	r.logger.Tracef("SQL query: %s", segmentdb.QueryToString(q))
	rows, err := r.client.Query(ctx, q, user.UserId, user.InDate)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}

	for rows.Next() {
		data := userinsegment.UserInSegmentsHistory{UserId: user.UserId, Event: event}
		if err = rows.Scan(&data.EventDate, &data.SegmentName); err != nil {
			r.logger.Errorf("Something wrong with scanning segmentName from row of SQL-request")
			return err
		}
		*history = append(*history, data)
	}
	if err = rows.Err(); err != nil {
		r.logger.Errorf("Something wrong with getting rows of SQL-request")
		return err
	}
	return nil
}

func (r repository) DeleteFromSegment(ctx context.Context, user *userinsegment.UserInSegmentDTO) error {
	q := `
		UPDATE user_in_segment 
		SET out_date = current_timestamp
        WHERE user_id = $1 
          AND segment_id = (SELECT segment_id
                              FROM segment
                              WHERE segment_name = $2) 
          AND out_date IS NULL
		`

	r.logger.Tracef("SQL query: %s", segmentdb.QueryToString(q))
	if _, err := r.client.Query(ctx, q, user.UserId, user.SegmentName); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}
	return nil
}

func (r repository) AddSegments(ctx context.Context, user *userinsegment.UserDTO, segment *segment.SegmentDTO) error {
	q := `
		INSERT INTO user_in_segment
		    (user_id, segment_id, in_date, out_date)
		VALUES ($1, (SElECT segment_id FROM segment WHERE segment_name = $2),
        current_timestamp, DEFAULT)
		`
	r.logger.Tracef("SQL query: %s", segmentdb.QueryToString(q))

	if _, err := r.client.Query(ctx, q, user.UserId, segment.SegmentName); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}

	return nil
}

func (r repository) AddSegmentsWithPeriod(ctx context.Context, user *userinsegment.UserDTO, segment *segment.SegmentDTO) error {
	q := `
		INSERT INTO user_in_segment
		    (user_id, segment_id, in_date, out_date)
		VALUES ($1, (SElECT segment_id FROM segment WHERE segment_name = $2),
        current_timestamp, current_timestamp + $3 * (interval '1 day'))
		`
	r.logger.Tracef("SQL query: %s", segmentdb.QueryToString(q))
	if _, err := r.client.Query(ctx, q, user.UserId, segment.SegmentName, user.Period); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}
	return nil
}

func (r repository) GetSegments(ctx context.Context, user *userinsegment.UserInSegment) (userinsegment.UserSegmentsList, error) {
	q := `
		SELECT segment_name 
		FROM (SELECT segment_id AS s_id
		      FROM user_in_segment
			  WHERE user_id = $1
			    AND (out_date IS NULL
			             OR out_date > current_date)
			  ) AS active_u_segments
		    INNER JOIN segment ON s_id = segment.segment_id
		`
	r.logger.Tracef("SQL query: %s", segmentdb.QueryToString(q))
	rows, err := r.client.Query(ctx, q, user.UserId)
	segments := userinsegment.UserSegmentsList{UserId: user.UserId, SegmentNames: []string{}}
	if err != nil {
		return segments, err
	}
	for rows.Next() {
		var segName string
		if err = rows.Scan(&segName); err != nil {
			r.logger.Errorf("Something wrong with scanning segmentName from row of SQL-request")
			return segments, err
		}
		segments.SegmentNames = append(segments.SegmentNames, segName)
	}
	if err = rows.Err(); err != nil {
		r.logger.Errorf("Something wrong with getting rows of SQL-request")
		return segments, err
	}
	return segments, nil
}

func NewRepository(l *logging.Logger, c postgresql.Client) userinsegment.Repository {
	return &repository{
		logger: l,
		client: c,
	}
}
