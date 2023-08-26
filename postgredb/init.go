package postgredb

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "3846936720"
	dbname   = "dogcare"
)

func InitDB() *sql.DB {
	psqlConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlConnection)
	checkError(err)
	err = db.Ping()
	checkError(err)

	createTables(db)
	//createTriggers(db)
	createIndexes(db)

	return db
}

func createTables(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS user_in_segment (" +
		"user_id serial NOT NULL, " +
		"segment_id integer NOT NULL REFERENCES segment (segment_id) ON DELETE RESTRICT ON UPDATE RESTRICT, " +
		"in_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, " +
		"out_date TIMESTAMP DEFAULT NULL, " +
		"PRIMARY KEY(user_id,segment_id)" +
		")")
	checkError(err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS segment (" +
		"segment_id serial NOT NULL, " +
		"segment_name varchar(50) NOT NULL UNIQUE, " +
		"active boolean DEFAULT FALSE, " +
		"PRIMARY KEY(segment_id)" +
		")")
	checkError(err)
}

func createIndexes(db *sql.DB) {
	_, err := db.Exec("CREATE INDEX out_date_idx ON user_in_segment (out_date)")
	checkError(err)
}

//func createTriggers(db *sql.DB) {
//	_, err := db.Exec("CREATE TABLE IF NOT EXISTS segment (" +
//		"segment_id integer NOT NULL, " +
//		"segment_name varchar(50) NOT NULL UNIQUE, " +
//		"PRIMARY KEY(segment_id)" +
//		")")
//	checkError(err)
//}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
