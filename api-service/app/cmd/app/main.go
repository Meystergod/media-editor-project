package main

import (
	"context"
	"github.com/Meystergod/media-editor-project/api-service/internal/app"
	"github.com/Meystergod/media-editor-project/api-service/internal/config"
	"github.com/Meystergod/media-editor-project/api-service/pkg/logging"
)

func main() {
	cfg := config.GetConfig()

	logging.Init(cfg.AppConfig.LogLevel)
	logger := logging.GetLogger()
	logger.Info("logger and config initialized")

	application, err := app.NewApp(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("running api server")
	if application.Run(context.Background()) != nil {
		logger.Fatal(err)
	}
}
