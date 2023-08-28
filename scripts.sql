SET timezone TO 'Europe/Moscow';

CREATE TABLE IF NOT EXISTS segment (
    segment_id serial NOT NULL,
    segment_name varchar(50) NOT NULL UNIQUE,
    active boolean DEFAULT TRUE,
    PRIMARY KEY(segment_id)
);

CREATE TABLE IF NOT EXISTS user_in_segment (
    user_id serial NOT NULL,
    segment_id integer NOT NULL REFERENCES segment (segment_id)
    ON DELETE RESTRICT ON UPDATE RESTRICT,
    in_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

CREATE INDEX IF NOT EXISTS out_date_idx ON user_in_segment (out_date);


INSERT INTO segment VALUES (DEFAULT, 'AVITO_1'),(DEFAULT, 'AVITO_2'),(DEFAULT, 'AVITO_3');

INSERT INTO user_in_segment VALUES (1000, 1);

SELECT segment_name FROM (SELECT segment_id AS s_id
                          FROM user_in_segment
                          WHERE user_id = 1000
                            AND (out_date IS NULL
                                     OR out_date < current_date)) AS active_u_segments
    INNER JOIN segment ON s_id = segment.segment_id;


INSERT INTO user_in_segment(user_id, segment_id, in_date, out_date)
VALUES ($1, (SElECT segment_id FROM segment WHERE segment_name = $2),
        current_timestamp, DEFAULT);

INSERT INTO user_in_segment(user_id, segment_id, in_date, out_date)
VALUES ($1, (SElECT segment_id FROM segment WHERE segment_name = $2),
        current_timestamp, $3);

UPDATE user_in_segment SET out_date = current_timestamp
                       WHERE segment_id =
                             (SELECT segment_id
                              FROM segment
                              WHERE segment_name = $1)
                         AND user_id = $2;



SELECT in_date,
       (SELECT segment_name
        from segment
        where user_in_segment.segment_id = segment.segment_id)
from user_in_segment
WHERE user_id = $1 AND in_date >= $2 AND in_date < $3;


SELECT out_date,
       (SELECT segment_name
        from segment
        where user_in_segment.segment_id = segment.segment_id)
from user_in_segment
WHERE user_id = $1 AND out_date >= $2 AND out_date < $3;

UPDATE segment SET active = false WHERE segment_name = $1;

UPDATE user_in_segment
SET out_date = current_timestamp
WHERE segment_id = (SELECT segment_id FROM segment WHERE segment_name = $1);
