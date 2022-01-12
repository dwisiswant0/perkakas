package collector

type ICollector interface {
	Collect()
	Name() string
}
