package app

import (
	"errors"
	"github.com/mylastgame/yp-metrics-service/internal/domain/counter"
	"github.com/mylastgame/yp-metrics-service/internal/domain/gauge"
	counterRepo "github.com/mylastgame/yp-metrics-service/internal/storage/counter"
	gaugeRepo "github.com/mylastgame/yp-metrics-service/internal/storage/gauge"
)

type App struct {
	gaugeRepo   gauge.Repo
	counterRepo counter.Repo
}

func New() *App {
	return &App{
		gaugeRepo:   gaugeRepo.NewMemRepo(),
		counterRepo: counterRepo.NewMemRepo(),
	}
}

func (app *App) Save(mtype, title, val string) error {
	if mtype == gauge.Key {
		return app.SaveGauge(title, val)
	}

	if mtype == counter.Key {
		return app.SaveCounter(title, val)
	}

	return errors.New("bad metric type")
}

func (app *App) SaveGauge(title, val string) error {
	g, err := gauge.FromString(title, val)
	if err != nil {
		return err
	}
	return app.gaugeRepo.Save(g)
}

func (app *App) SaveCounter(title, val string) error {
	c, err := counter.FromString(title, val)
	if err != nil {
		return err
	}
	return app.counterRepo.Save(c)
}
