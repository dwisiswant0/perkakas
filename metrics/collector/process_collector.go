package collector

import (
	"os"
	"sync"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/prometheus/procfs"
	"github.com/rs/zerolog/log"
)

const (
	openFileDescriptor           = "process.open_file_descriptor"
	maxOpenFileDescriptor        = "process.max_file_descriptor"
	processVirtualMemory         = "process.virtual_memory_bytes"
	processVirtualMemoryMaxBytes = "process.virtual_memory_max_bytes"
	residentMemory               = "process.resident_memory"
)

type ProcessCollector struct {
	pid int
	st  *statsd.Client
}

func NewProcessCollector(st *statsd.Client) *ProcessCollector {
	return &ProcessCollector{
		pid: os.Getpid(),
		st:  st,
	}
}

func (p *ProcessCollector) Name() string {
	return "process_collector"
}

func (p *ProcessCollector) Collect() {
	if p.pid == 0 || !canCollectProcess() || p.st == nil {
		return
	}

	if err := p.composer(
		procfs.NewProc,
		[]func(procfs.Proc, func(name string, value float64, tags []string, rate float64) error) func(){
			p.collectVirtualMemory,
			p.collectProcessDescriptor,
			p.collectLimit,
		},
		p.st.Gauge,
	)(); err != nil {
		log.Error().Err(err).Msg("process collector composer")
	}
}

func (p *ProcessCollector) composer(
	procFn func(pid int) (procfs.Proc, error),
	collectorFuncs []func(procfs.Proc, func(name string, value float64, tags []string, rate float64) error) func(),
	gaugeFn func(name string, value float64, tags []string, rate float64) error,
) func() error {
	return func() error {
		proc, err := procFn(p.pid)
		if err != nil {
			log.Error().Err(err).Msg("new proc")
			return err
		}

		var wg sync.WaitGroup
		wg.Add(len(collectorFuncs))
		for _, fn := range collectorFuncs {
			go func(f func(
				proc procfs.Proc,
				gaugeFn func(name string, value float64, tags []string, rate float64) error,
			) func()) {
				defer wg.Done()
				f(proc, gaugeFn)
			}(fn)
		}

		wg.Wait()

		return nil
	}
}

func (p *ProcessCollector) collectVirtualMemory(
	proc procfs.Proc,
	gauge func(name string, value float64, tags []string, rate float64) error,
) func() {
	return func() {
		stat, err := proc.Stat()
		if err != nil {
			log.Error().Err(err).Msg("proc stat")
			return
		}

		if err = gauge(
			processVirtualMemory,
			float64(stat.VirtualMemory()),
			[]string{},
			1,
		); err != nil {
			log.Error().Err(err).Msg("collect virtial memory")
		}

		if err = gauge(
			residentMemory,
			float64(stat.ResidentMemory()),
			[]string{},
			1,
		); err != nil {
			log.Error().Err(err).Msg("collect resident memory")
		}
	}
}

func (p *ProcessCollector) collectProcessDescriptor(
	proc procfs.Proc,
	gauge func(name string, value float64, tags []string, rate float64) error,
) func() {
	return func() {
		fds, err := proc.FileDescriptorsLen()
		if err != nil {
			log.Error().Err(err).Msg("file descriptor len")
			return
		}

		if err = gauge(
			openFileDescriptor,
			float64(fds),
			[]string{},
			1,
		); err != nil {
			log.Error().Err(err).Msg("collect process descriptor")
		}
	}
}

func (p *ProcessCollector) collectLimit(
	proc procfs.Proc,
	gauge func(name string, value float64, tags []string, rate float64) error,
) func() {
	return func() {
		limits, err := proc.Limits()
		if err != nil {
			log.Error().Err(err).Msg("limits proc")
			return
		}

		if err = gauge(
			maxOpenFileDescriptor,
			float64(limits.OpenFiles),
			[]string{},
			1,
		); err != nil {
			log.Error().Err(err).Msg("limit open file")
		}

		if err = gauge(
			processVirtualMemoryMaxBytes,
			float64(limits.AddressSpace),
			[]string{},
			1,
		); err != nil {
			log.Error().Err(err).Msg("limit max virtual memory")
		}

	}
}

func canCollectProcess() bool {
	info, err := os.Stat("/proc")
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return true
}
