package main

import (
	"AvitoTechTask/internal/configuration"
	"AvitoTechTask/internal/segment"
	"AvitoTechTask/internal/userinsegment"
	"AvitoTechTask/pkg/logging"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	//dbase := db.InitDB()
	//defer dbase.Close()
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()
	config := configuration.GetConfig(logger)

	registerHandlers(router, logger)
	start(router, logger, config)
}

func registerHandlers(router *httprouter.Router, logger *logging.Logger) {
	userHandler := userinsegment.NewHandler(logger)
	userHandler.Register(router)

	segmentHandler := segment.NewHandler(logger)
	segmentHandler.Register(router)
}

func start(router *httprouter.Router, logger *logging.Logger, config *configuration.Config) {
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
