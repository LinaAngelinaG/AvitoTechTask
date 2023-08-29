package main

import (
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

func main() {
	//dbase := db.InitDB()
	//defer dbase.Close()
	logger = logging.GetLogger()
	logger.Info("create router")
	router = httprouter.New()
	config = configuration.GetConfig(logger)

	startRepositories()

	registerHandlers()
	start()
}

func startRepositories() {
	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, config)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	logger.Info("postgres client successfully connected")
	segRepo = segmentdb.NewRepository(logger, postgreSQLClient)
	userRepo = userinsegmentdb.NewRepository(logger, postgreSQLClient)

	//////////-----------------------------------------------------------------------------------------------
	//newSeg := segment.Segment{SegmentName: "AVITO_SEG_TEST_1"}
	//err = segRepo.Create(context.TODO(), &newSeg)
	//if err != nil {
	//	logger.Fatalf("%v", err)
	//}
	//logger.Info(newSeg)
	//err = segRepo.Create(context.TODO(), &segment.Segment{SegmentName: "AVITO_SEG_TEST_2"})
	//if err != nil {
	//	logger.Fatalf("%v", err)
	//}
	//
	//err = userRepo.AddSegments(context.TODO(), &userinsegment.UserInSegment{UserId: 1000}, &segment.SegmentDTO{"AVITO_SEG_TEST_2"})
	//if err != nil {
	//	logger.Fatalf("%v", err)
	//}
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
