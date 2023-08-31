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