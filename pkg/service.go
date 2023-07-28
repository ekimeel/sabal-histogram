package tq

import (
	"context"
	"github.com/ekimeel/sabal-pb/pb"
	"github.com/ekimeel/sabal-plugin/pkg/metric_utils"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	PluginName    = "histogram"
	PluginVersion = "v1.0"
)

var (
	singletonService *Service
	onceService      sync.Once
)

type Service struct {
	dao *dao
}

func GetService() *Service {

	onceService.Do(func() {
		singletonService = &Service{}
		singletonService.dao = getDao()
	})

	return singletonService
}

func (s *Service) Run(ctx context.Context, metrics []*pb.Metric) {

	unitOfWork := metric_utils.GroupMetricsByPointId(metrics)
	var wg sync.WaitGroup

	for pointId, items := range unitOfWork {
		wg.Add(1)
		go func(pointId uint32, items []*pb.Metric) {
			defer wg.Done()
			s.compute(pointId, items)
		}(pointId, items)
	}

	wg.Wait()
}

func (s *Service) compute(pointId uint32, metrics []*pb.Metric) {
	hist, err := s.dao.selectByPointId(pointId)
	if err != nil {
		log.WithField("plugin", PluginName).Errorf("dao error: %s", err)
	}

	if hist == nil {
		hist = &Histogram{
			PointId:     pointId,
			LastUpdated: time.Now(),
			Histogram:   make(map[string]int, 0),
		}

		hist.update(metrics)
		_, err := s.dao.insert(hist)
		if err != nil {
			log.WithField("plugin", PluginName).Errorf("failed to insert: %s", err)
		}
	} else {
		hist.update(metrics)
		_, err := s.dao.update(hist)
		if err != nil {
			log.WithField("plugin", PluginName).Errorf("failed ot update: %s", err)
		}
	}

}
