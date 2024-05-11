package main

import (
	"context"
	"fmt"
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
		fmt.Printf("Error parsing flags: %v/n", err)
		panic(err)
	}

	//logger init
	log, err := logger.NewLogger(config.LogLevel)
	if err != nil {
		fmt.Printf("Error init logger: %v/n", err)
		panic(err)
	}

	//init memory repository
	repo := storage.NewMemRepo()

	//init file storage
	fileStorage := storage.NewFileStorage(repo, log)

	//restore from file
	if config.Restore {
		err = fileStorage.Restore()
		if err != nil {
			log.Sugar.Errorf("Error restoring data: %v", err)
			panic(err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if config.StoreInterval != 0 {
		//run store to file go-routine
		go func(ctx context.Context) {
			storeTicker := time.NewTicker(time.Duration(config.StoreInterval) * time.Second)
			for {
				select {
				case <-ctx.Done():
					return
				case <-storeTicker.C:
					err := fileStorage.Save()
					if err != nil {
						log.Sugar.Errorf("Error saving file: %v", err)
						return
					}
				}

			}
		}(ctx)
	}

	r := app.NewRouter(repo, fileStorage, log)
	log.Sugar.Infof("Starting server. Listening on %s", config.RunAddr)
	log.Sugar.Infof("Store file: %s, interval: %d, restore: %t",
		config.FileStoragePath,
		config.StoreInterval,
		config.Restore,
	)
	err = http.ListenAndServe(config.RunAddr, r)

	if err != nil {
		log.Sugar.Errorf("Error starting server: %v", err)
		panic(err)
	}
}
