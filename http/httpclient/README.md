# HttpClient
Http client with resillience factor built in. Utilizing [Heimdall](https://github.com/gojek/heimdall) as it library. Default retrier is `Constant Retrier`. Please see [Heimdall usage](https://github.com/gojek/heimdall#usage) for more details.

`NewHttpClient(*configuration)` for creating new http client. If `configuration` is nil, default value will be used.

Example:
```go
conf := new(HttpClientConf)
conf.BackoffInterval = 2 * time.Millisecond       // 2ms
conf.MaximumJitterInterval = 5 * time.Millisecond // 5ms
conf.Timeout = 15000 * time.Millisecond           // 15s
conf.RetryCount = 3  // 3 times

h := NewHttpClient(conf)
resp, err := h.Client.Get("http://some-url", headers)
// Do something with response
```