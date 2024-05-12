package app

import (
	"context"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"time"
)

func BackupMetrics(ctx context.Context, fileStorage *storage.FileStorage, log *logger.Logger) {
	if config.StoreInterval != 0 {
		//run store to file go-routine
		go func(ctx context.Context) {
			storeTicker := time.NewTicker(time.Duration(config.StoreInterval) * time.Second)
			for {
				select {
				case <-ctx.Done():
					return
				case <-storeTicker.C:
					err := fileStorage.Save(ctx)
					if err != nil {
						log.Sugar.Errorf("Error saving file: %v", err)
						return
					}
				}

			}
		}(ctx)
	}
}
