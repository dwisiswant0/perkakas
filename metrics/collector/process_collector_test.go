package collector

import (
	"os"
	"runtime"
	"testing"

	"github.com/prometheus/procfs"
)

func TestProcessCollectorComposer(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip()
	}

	p := &ProcessCollector{
		pid: os.Getpid(),
	}

	p.composer(
		procfs.NewProc,
		[]func(procfs.Proc, gaugeFunc){
			p.collectVirtualMemory,
			p.collectProcessDescriptor,
			p.collectLimit,
		},
		func(name string, value float64, tags []string, rate float64) error {
			return nil
		},
	)()
}
