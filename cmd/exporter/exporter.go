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

	counters := make([]Counter, len(config.Functions))
	for i := 0; i < len(config.Functions); i++ {
		counter := promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_downtime", config.Functions[i].Id),
			Help: config.Functions[i].Description,
		})
		counters[i] = Counter{
			id:      config.Functions[i].Id,
			counter: counter,
		}
	}

	ex.counters = counters

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
