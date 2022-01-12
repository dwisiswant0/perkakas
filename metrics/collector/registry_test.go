package collector

import (
	"testing"
	"time"
)

func TestRegistry(t *testing.T) {
	cases := []struct {
		name       string
		collectors []ICollector
		wantRes    int
	}{
		{
			name:       "when test registry with 1 collector, should be return 1",
			collectors: []ICollector{NewGoCollector(nil)},
			wantRes:    1,
		},
		{
			name:       "when test registry with 2 collector, should be return 2",
			collectors: []ICollector{NewGoCollector(nil), NewProcessCollector(nil)},
			wantRes:    2,
		},
		{
			name:       "when test registry with 2 duplicate collector, should be return 1",
			collectors: []ICollector{NewGoCollector(nil), NewGoCollector(nil)},
			wantRes:    1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reg := NewRegistry(0)
			for _, v := range c.collectors {
				reg.Register(v)
			}

			if reg.NumCollectors() != c.wantRes {
				t.Errorf("num collectors should be %d", c.wantRes)
			}

			go reg.Collect()

			// give time to react
			time.Sleep(2 * time.Second)

			reg.Stop()
		})
	}

}
