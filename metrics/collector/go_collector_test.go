package collector

import "testing"

func TestGoCollectorComposer(t *testing.T) {
	g := NewGoCollector(nil)
	gaugeFn := func(name string, value float64, tags []string, rate float64) error {
		return nil
	}

	countFn := func(name string, value int64, tags []string, rate float64) error {
		return nil
	}

	g.composer([]func(){
		g.collectNumOfGoroutine(gaugeFn),
		g.collectGC(gaugeFn),
		g.collectThread(gaugeFn),
		g.collectMemory(gaugeFn, countFn),
	})()
}
