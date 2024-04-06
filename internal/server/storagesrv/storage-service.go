package storagesrv

import (
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/storage"
	"regexp"
)

type StorageServiceI interface {
	Save(string, string, string) error
}

type StorageService struct {
	storage storage.StorageI
}

func New(s storage.StorageI) *StorageService {
	return &StorageService{storage: s}
}

// Convert and save metric in storage
func (srv *StorageService) Save(mtype, mtitle, mval string) error {
	fmt.Println("mtype: ", mtype, ", mtitle: ", mtitle, ", mval: ", mval)
	//check metrics type
	if !metrics.TypeExists(mtype) {
		return fmt.Errorf("type %s not exists", mtype)
	}

	//check metrics title
	checkTitle := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	if !checkTitle(mtitle) {
		return fmt.Errorf("Title %s is not valid!", mtitle)
	}

	if mtype == metrics.GaugeKey {
		value, err := metrics.ConvertToGauge(mval)
		if err != nil {
			return err
		}
		srv.storage.AddGauge(mtitle, value)
	}

	if mtype == metrics.CounterKey {
		value, err := metrics.ConvertToCounter(mval)
		if err != nil {
			return err
		}
		srv.storage.AddCounter(mtitle, value)
	}

	fmt.Println(srv.storage)

	return nil
}
