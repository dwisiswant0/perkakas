package collector

import (
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/rs/zerolog/log"
)

const (
	goroutineStats      = "goroutine"
	gcStats             = "gc"
	threadCreationStats = "thread"

	memstatsAlloc         = "memstats.alloc_bytes"
	memstatsTotalAlloc    = "memstats.total_alloc_bytes"
	memStatsSysBytes      = "memstats.sys_bytes"
	memstatsHeapAllocate  = "memstats.heap_alloc"
	memstatsHeapIdle      = "memstats.heap_idle"
	memstatsHeapInUse     = "memstats.heap_inuse"
	memstatsHeapSysBytes  = "memstats.heap_sys_bytes"
	memstatsStackInUse    = "memstats.stack_inuse"
	memstatsStackSysBytes = "memstats.stack_sys_bytes"
	memstatsGcSysBytes    = "memstats.gc_sys_bytes"
	memstatsCpuFraction   = "memstats.cpu_fraction"
)

type metricAggregate string

const (
	gaugeAggregate metricAggregate = "gauge"
	countAggregate metricAggregate = "count"
)

type goCollector struct {
	sync.Mutex
	st                 *statsd.Client
	readMem            func(*runtime.MemStats)
	memLast            *runtime.MemStats
	memLastTime        time.Time
	memMaxWait         time.Duration
	memMaxAge          time.Duration
	memStatsCollection map[string]memStatsCollection
}

func NewGoCollector(st *statsd.Client) *goCollector {
	return &goCollector{
		st:         st,
		readMem:    runtime.ReadMemStats,
		memMaxWait: time.Second,
		memMaxAge:  5 * time.Minute,
		memStatsCollection: map[string]memStatsCollection{
			memstatsAlloc: {
				aggregateType: gaugeAggregate,
				get: func(ms *runtime.MemStats) float64 {
					return float64(ms.Alloc)
				},
			},
			memstatsTotalAlloc: {
				aggregateType: countAggregate,
				get: func(ms *runtime.MemStats) float64 {
					return float64(ms.TotalAlloc)
				},
			},
			memStatsSysBytes: {
				aggregateType: gaugeAggregate,
				get: func(ms *runtime.MemStats) float64 {
					return float64(ms.Sys)
				},
			},
			memstatsHeapAllocate: {
				aggregateType: gaugeAggregate,
				get: func(ms *runtime.MemStats) float64 {
					return float64(ms.HeapAlloc)
				},
			},
			memstatsHeapIdle: {
				aggregateType: gaugeAggregate,
				get: func(ms *runtime.MemStats) float64 {
					return float64(ms.HeapIdle)
				},
			},
			memstatsHeapInUse: {
				aggregateType: gaugeAggregate,
				get: func(ms *runtime.MemStats) float64 {
					return float64(ms.HeapInuse)
				},
			},
			memstatsHeapSysBytes: {
				aggregateType: gaugeAggregate,
				get: func(ms *runtime.MemStats) float64 {
					return float64(ms.HeapSys)
				},
			},
			memstatsStackInUse: {
				aggregateType: gaugeAggregate,
				get: func(ms *runtime.MemStats) float64 {
					return float64(ms.StackInuse)
				},
			},
			memstatsStackSysBytes: {
				aggregateType: gaugeAggregate,
				get: func(ms *runtime.MemStats) float64 {
					return float64(ms.StackSys)
				},
			},
			memstatsGcSysBytes: {
				aggregateType: gaugeAggregate,
				get: func(ms *runtime.MemStats) float64 {
					return float64(ms.GCSys)
				},
			},
			memstatsCpuFraction: {
				aggregateType: gaugeAggregate,
				get: func(ms *runtime.MemStats) float64 {
					return float64(ms.GCCPUFraction)
				},
			},
		},
	}
}

func (g *goCollector) Collect() {
	if g.st == nil {
		return
	}

	g.composer([]func(){
		g.collectNumOfGoroutine(g.st.Gauge, g.st.Count),
		g.collectGC(g.st.Gauge, g.st.Count),
		g.collectThread(g.st.Gauge, g.st.Count),
		g.collectMemory(g.st.Gauge, g.st.Count),
	})()
}

func (g *goCollector) Name() string {
	return "go_collector"
}

func (g *goCollector) composer(collectorFuncs []func()) func() {
	return func() {
		var wg sync.WaitGroup
		wg.Add(len(collectorFuncs))

		for _, cFn := range collectorFuncs {
			go func(fn func()) {
				defer wg.Done()
				fn()
			}(cFn)
		}

		wg.Wait()
	}
}

func (g *goCollector) collectNumOfGoroutine(
	gaugeFn func(name string, value float64, tags []string, rate float64) error,
	_ func(name string, value int64, tags []string, rate float64) error,
) func() {
	return func() {
		if err := gaugeFn(goroutineStats, float64(runtime.NumGoroutine()), []string{}, 1); err != nil {
			log.Error().Err(err).Msg("collect number of goroutine")
		}
	}
}

func (g *goCollector) collectGC(
	gaugeFn func(name string, value float64, tags []string, rate float64) error,
	_ func(name string, value int64, tags []string, rate float64) error,
) func() {
	return func() {
		var stats debug.GCStats
		debug.ReadGCStats(&stats)
		if err := gaugeFn(gcStats, float64(stats.NumGC), []string{}, 1); err != nil {
			log.Error().Err(err).Msg("collect gc")
		}
	}
}

func (g *goCollector) collectThread(
	gaugeFn func(name string, value float64, tags []string, rate float64) error,
	_ func(name string, value int64, tags []string, rate float64) error,
) func() {
	return func() {
		n, _ := runtime.ThreadCreateProfile(nil)
		if err := gaugeFn(threadCreationStats, float64(n), []string{}, 1); err != nil {
			log.Error().Err(err).Msg("collect thread")
		}
	}
}

func (g *goCollector) collectMemory(
	gaugeFn func(name string, value float64, tags []string, rate float64) error,
	countFn func(name string, value int64, tags []string, rate float64) error,
) func() {
	return func() {
		var (
			ms   = runtime.MemStats{}
			done = make(chan struct{})
		)

		go func() {
			g.readMem(&ms)
			g.Lock()
			g.memLast = &ms
			g.memLastTime = time.Now()
			g.Unlock()
			close(done)
		}()

		timer := time.NewTimer(g.memMaxWait)
		select {
		case <-done:
			timer.Stop()
			g.deletegateMemCollector(&ms, gaugeFn, countFn)
			return
		case <-timer.C:
		}

		g.Lock()
		if time.Since(g.memLastTime) < g.memMaxAge {
			g.deletegateMemCollector(g.memLast, gaugeFn, countFn)
			g.Unlock()
			return
		}
		g.Unlock()
		<-done

		g.deletegateMemCollector(&ms, gaugeFn, countFn)
	}
}

func (g *goCollector) deletegateMemCollector(
	memstats *runtime.MemStats,
	gaugeFn func(name string, value float64, tags []string, rate float64) error,
	countFn func(name string, value int64, tags []string, rate float64) error,
) {
	for k, v := range g.memStatsCollection {
		if v.aggregateType == gaugeAggregate {
			if err := gaugeFn(k, v.get(memstats), []string{}, 1); err != nil {
				log.Error().Err(err).Msgf("collect mem: %s", k)
			}
		} else if v.aggregateType == countAggregate {
			if err := countFn(k, int64(v.get(memstats)), []string{}, 1); err != nil {
				log.Error().Err(err).Msgf("collect mem: %s", k)
			}
		}
	}
}

type memStatsCollection struct {
	aggregateType metricAggregate
	get           func(ms *runtime.MemStats) float64
}
