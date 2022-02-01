package metrics

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/rs/zerolog/log"
)

// IMetric ...
type IMetric interface {
	IncrErr(tag []string)
	IncrSuccess(tag []string)

	IncrErrRedis()
	IncrErrDB()
	IStatsD
}

type IStatsD  interface {
	GetStatsdConfig() *StatsDConfig
	GetStatsdClient() *statsd.Client
}

// Metric ...
type Metric struct {
	statsDClient *statsd.Client
	StatsDConfig
}

type StatsDConfig struct {
	Host string
	Namespace string
	Port int
	IsEnable bool
}

// table metrics
var (
	TableErr     = "ERROR"
	TableSuccess = "SUCCESS"
)

// tags
var (
	TErrRedis   = []string{"type:redis"}
	TErrDB      = []string{"type:db"}
	TErrService = []string{"type:service"}
)

// NewMetric ...
func NewMetric(statsDConfig StatsDConfig) *Metric {
	telegrafMetrics := &Metric{
		StatsDConfig: statsDConfig,
	}

	if statsDConfig.IsEnable {
		host := fmt.Sprintf("%s:%d", statsDConfig.Host, statsDConfig.Port)

		c, err := statsd.New(host)
		if err != nil {
			panic(err)
		}

		c.Namespace = telegrafMetrics.Namespace
		telegrafMetrics.statsDClient = c
	}

	return telegrafMetrics
}

// IncrErr general error metric
func (m *Metric) IncrErr(tag []string) {
	if !m.StatsDConfig.IsEnable {
		return
	}

	err := m.statsDClient.Incr(TableErr, tag, 1)
	if err != nil {
		log.Err(err).Send()
	}
}

// IncrSuccess general success metric
func (m *Metric) IncrSuccess(tag []string) {
	if !m.StatsDConfig.IsEnable {
		return
	}

	err := m.statsDClient.Incr(TableSuccess, tag, 1)
	if err != nil {
		log.Err(err).Send()
	}
}

// IncrErrRedis for metric error redis
func (m *Metric) IncrErrRedis() {
	if !m.StatsDConfig.IsEnable {
		return
	}

	err := m.statsDClient.Incr(TableErr, TErrRedis, 1)
	if err != nil {
		log.Err(err).Send()
	}
}

// IncrErrDB for metric error database
func (m *Metric) IncrErrDB() {
	if !m.StatsDConfig.IsEnable {
		return
	}

	err := m.statsDClient.Incr(TableErr, TErrDB, 1)
	if err != nil {
		log.Err(err).Send()
	}
}

func (m *Metric) GetStatsdConfig() *StatsDConfig {
	return &m.StatsDConfig
}

func (m *Metric) GetStatsdClient() *statsd.Client {
	return m.statsDClient
}

