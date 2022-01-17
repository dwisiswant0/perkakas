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
		g.collectNumOfGoroutine(gaugeFn, countFn),
		g.collectGC(gaugeFn, countFn),
		g.collectThread(gaugeFn, countFn),
		g.collectMemory(gaugeFn, countFn),
	})()
}
