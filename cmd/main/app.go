package main

import (
	"AvitoTechTask/internal/segment"
	"AvitoTechTask/internal/userinsegment"
	"AvitoTechTask/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
)

type User struct {
	UserId   int      `json:"user_id"`
	Segments []string `json:"segments"`
}

//const CreateSeg = "newSegment"
//const DeleteSeg = "delSegment"
//const AddUserSeg = "newUserSegment"
//const DelUserSeg = "newSegment"

func main() {
	//dbase := db.InitDB()
	//defer dbase.Close()
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	logger.Info("register user handler")
	userHandler := userinsegment.NewHandler()
	userHandler.Register(router)

	log.Println("register segment handler")
	segmentHandler := segment.NewHandler()
	segmentHandler.Register(router)

	start(router)
}

func start(router *httprouter.Router) {
	logger := logging.GetLogger()
	logger.Info("start application")

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Info("server is listening port 0.0.0.0:1234")
	log.Fatal(server.Serve(listener))
}
