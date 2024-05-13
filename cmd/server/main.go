package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/core"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"net/http"
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

	//Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//init memory repository
	var repo storage.Repo
	if config.DBConnect == "" {
		repo = storage.NewMemRepo()
	} else {
		var db *sql.DB
		db, err = sql.Open("pgx", config.DBConnect)
		if err != nil {
			log.Sugar.Errorf("Error connecting to database: %v", err)
			panic(err)
		}
		defer db.Close()

		//create DB repo
		err = core.Retry("init DB repo", 3, func() error {
			repo, err = storage.NewDBRepo(ctx, db)
			return err
		}, log)

		if err != nil {
			log.Sugar.Error(err)
			panic(err)
		}
	}

	//init file storage
	fileStorage := storage.NewFileStorage(repo, log)

	//restore from file
	if config.Restore {
		err = fileStorage.Restore(ctx)
		if err != nil {
			log.Sugar.Errorf("Error restoring data: %v", err)
			panic(err)
		}
	}

	app.BackupMetrics(ctx, fileStorage, log)

	r := app.NewRouter(repo, fileStorage, log)
	log.Sugar.Infof("Starting server. Listening on %s", config.RunAddr)
	log.Sugar.Infof("Store file: %s, interval: %d, restore: %t",
		config.FileStoragePath,
		config.StoreInterval,
		config.Restore,
	)
	log.Sugar.Infof("DB connecting data: %s", config.DBConnect)

	err = http.ListenAndServe(config.RunAddr, r)

	if err != nil {
		log.Sugar.Errorf("Error starting server: %v", err)
		panic(err)
	}
}
