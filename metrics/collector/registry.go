package collector

import (
	"time"

	"github.com/rs/zerolog/log"
)

var (
	defaultDuration = 2 * time.Second
)

type Registry struct {
	collectors map[string]ICollector
	interval   time.Duration
	done       chan struct{}
}

func NewRegistry(interval time.Duration) *Registry {
	if interval == 0 {
		interval = defaultDuration
	}

	return &Registry{
		collectors: make(map[string]ICollector),
		interval:   interval,
		done:       make(chan struct{}),
	}
}

func (r *Registry) Register(c ICollector) {
	r.collectors[c.Name()] = c
}

func (r *Registry) Collect() {
	t := time.NewTicker(r.interval)

	for {
		r.runCollectors()
		select {
		case <-t.C:
		case <-r.done:
			log.Info().Msg("stopping collectors")
			t.Stop()
			return
		}
	}
}

func (r *Registry) runCollectors() {
	for _, collector := range r.collectors {
		collector.Collect()
	}
}

func (r *Registry) NumCollectors() int {
	return len(r.collectors)
}

func (r *Registry) Stop() {
	r.done <- struct{}{}
}
