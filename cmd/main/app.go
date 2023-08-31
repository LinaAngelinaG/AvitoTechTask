package main

import (
	_ "AvitoTechTask/docs"
	"AvitoTechTask/internal/configuration"
	"AvitoTechTask/internal/handlers"
	"AvitoTechTask/internal/segment"
	segmentdb "AvitoTechTask/internal/segment/db"
	"AvitoTechTask/internal/userinsegment"
	userinsegmentdb "AvitoTechTask/internal/userinsegment/db"
	postgresql "AvitoTechTask/pkg/client/postgres"
	"AvitoTechTask/pkg/logging"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
)

var logger *logging.Logger
var segRepo segment.Repository
var userRepo userinsegment.Repository
var router *httprouter.Router
var config *configuration.Config
var userHandler handlers.Handler
var segmentHandler handlers.Handler

// @title AvitoTechTask App API
// @version 1.0
// @description API Server for AvitoTechTask Application

// @host localhost:1234
// @BasePath /

func main() {
	logger = logging.GetLogger()
	logger.Info("create router")
	router = httprouter.New()
	config = configuration.GetConfig(logger)

	err := initDB()
	if err != nil {
		logger.Error(err.Error())
	}

	startRepositories()

	registerHandlers()
	start()
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func initDB() error {
	q1 := `
	CREATE TABLE IF NOT EXISTS segment (
	   segment_id serial NOT NULL,
	   segment_name varchar(50) NOT NULL UNIQUE,
	   active boolean DEFAULT TRUE,
	   PRIMARY KEY(segment_id)
	);
`
	q2 := `
	CREATE TABLE IF NOT EXISTS user_in_segment (
	   user_id serial NOT NULL,
	   segment_id serial NOT NULL REFERENCES segment (segment_id)
		   ON DELETE RESTRICT ON UPDATE RESTRICT,
	   in_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	   out_date TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	   primary key (user_id, segment_id)
	);
	`
	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, config)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	if _, err := postgreSQLClient.Query(context.TODO(), q1); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}
	logger.Info("Table ")
	if _, err := postgreSQLClient.Query(context.TODO(), q2); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			logger.Errorf("SQL error: %s, details: %s, where: %s, code: %s, SQL-state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}
	return nil
}

func startRepositories() {
	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, config)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	logger.Info("postgres client successfully connected")
	segRepo = segmentdb.NewRepository(logger, postgreSQLClient)
	userRepo = userinsegmentdb.NewRepository(logger, postgreSQLClient)
}

func registerHandlers() {
	userHandler = userinsegment.NewHandler(logger, userRepo)
	userHandler.Register(router)

	segmentHandler = segment.NewHandler(logger, segRepo)
	segmentHandler.Register(router)

}

func start() {

	logger.Info("start application")

	logger.Info("listen tcp")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Listen.BindIp, config.Listen.Port))
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Info("server is listening port %s:%s", config.Listen.BindIp, config.Listen.Port)
	log.Fatal(server.Serve(listener))
}
