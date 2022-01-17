package collector

type ICollector interface {
	Collect()
	Name() string
}

type (
	gaugeFunc func(name string, value float64, tags []string, rate float64) error
	countFunc func(name string, value int64, tags []string, rate float64) error
)
