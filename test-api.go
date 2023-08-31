package AvitoTechTask

import (
	"AvitoTechTask/internal/configuration"
	"AvitoTechTask/internal/segment"
	segmentdb "AvitoTechTask/internal/segment/db"
	postgresql "AvitoTechTask/pkg/client/postgres"
	"AvitoTechTask/pkg/logging"
	"context"
	"github.com/julienschmidt/httprouter"
	"testing"
)

var logger *logging.Logger
var router *httprouter.Router
var config *configuration.Config

func TestSegmentHandler(t *testing.T) {
	logger = logging.GetLogger()
	logger.Info("create router")
	router = httprouter.New()
	config = configuration.GetConfig(logger)

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, config)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	segRepo := segmentdb.NewRepository(logger, postgreSQLClient)
	segmentHandler := segment.NewHandler(logger, segRepo)
	segmentHandler.Register(router)

}
