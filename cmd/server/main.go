package main

import (
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"net/http"
	"time"
)

func main() {
	var err error
	//config init
	err = config.ParseFlags()
	if err != nil {
		logger.Sugar.Errorf("Error parsing flags: %v", err)
		panic(err)
	}

	//logger init
	err = logger.Initialize(config.LogLevel)
	if err != nil {
		logger.Sugar.Errorf("Error init logger: %v", err)
		panic(err)
	}

	//init memory repository
	repo := storage.NewMemRepo()

	//init file storage
	fileStorage := storage.NewFileStorage(repo)

	//restore from file
	if config.Restore {
		err = fileStorage.Restore()
		if err != nil {
			logger.Sugar.Errorf("Error restoring data: %v", err)
			panic(err)
		}
	}

	if config.StoreInterval != 0 {
		//run store to file go-routine
		go func() {
			storeTicker := time.NewTicker(time.Duration(config.StoreInterval) * time.Second)
			for range storeTicker.C {
				err := fileStorage.Save()
				if err != nil {
					logger.Sugar.Errorf("Error saving file: %v", err)
					return
				}
			}
		}()
	}

	r := app.NewRouter(repo, fileStorage)
	logger.Sugar.Infof("Starting server. Listening on %s", config.RunAddr)
	logger.Sugar.Infof("Store file: %s, interval: %d, restore: %t",
		config.FileStoragePath,
		config.StoreInterval,
		config.Restore,
	)
	err = http.ListenAndServe(config.RunAddr, r)

	if err != nil {
		logger.Sugar.Errorf("Error starting server: %v", err)
		panic(err)
	}
}
