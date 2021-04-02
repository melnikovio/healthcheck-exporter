package exporter

import (
	"fmt"
	"github.com/healthcheck-exporter/cmd/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
)

type Exporter struct {
	config   *model.Config
	counters []Counter
}

type Counter struct {
	id      string
	counter prometheus.Counter
}

func NewExporter(config *model.Config) *Exporter {
	ex := Exporter{
		config: config,
	}

	if config != nil {
		counters := make([]Counter, len(config.Jobs))
		for i := 0; i < len(config.Jobs); i++ {
			counter := promauto.NewCounter(prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_downtime", config.Jobs[i].Id),
				Help: config.Jobs[i].Description,
			})
			counters[i] = Counter{
				id:      config.Jobs[i].Id,
				counter: counter,
			}

			log.Info(fmt.Sprintf("Registered counter %s", config.Jobs[i].Id))
		}

		ex.counters = counters
	}

	return &ex
}

func (ex *Exporter) SetCounter(id string, value int64) {
	for i := 0; i < len(ex.counters); i++ {
		if ex.counters[i].id == id {
			val := float64(value)
			m := dto.Metric{
				Counter: &dto.Counter{
					Value: &val,
				},
			}

			err := ex.counters[i].counter.Write(&m)
			if err != nil {
				log.Error(fmt.Sprintf("Error writing metrics: %s", err.Error()))
			}
		}
	}
}

func (ex *Exporter) AddCounter(id string, value int64) {
	for i := 0; i < len(ex.counters); i++ {
		if ex.counters[i].id == id {
			val := float64(value)
			ex.counters[i].counter.Add(val)
		}
	}
}
