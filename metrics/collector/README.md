# Perkakas Collector
This library helps you to send go collector and process collector backed by statsd

# Component 
* go_collector: go collector including goroutine, gc, etc
* process_collector: process collector including virtual mem, resident mem, and file descriptor. only available in linux since it reads from `/proc`

# How To use
```
st, err := statsd.New(telegrafURL)
if err != nil {
	return err
}

reg := collector.NewRegistry(time.Second)

// register go collector
reg.Register(collector.NewGoCollector(st))

// register process collector
reg.Register(collector.NewProcessCollector(st))

// run it
go reg.Collect()

```