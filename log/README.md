# Log
This logger will help you logging on a request - response cycle and show the stack trace of the error when you add log message.
Logging level not supported. Level in this logger only for showing the severity level when add log messages.

`NewLogger("service_name")` will create a logger for you.

`SetRequest(req interface{})` is function expecting *http.Request, then will extract required information from the request struct to fill the log.
Later will support graphql and grpc also.
Input can also a primitive type, and should be string. It will treat the input as the request body in string.

`SetResponse(resp interface{}, body []byte)` will extract the information for the log from resp, expecting *http.Response or http.ResponseWriter.
Body will fill the response body in the log message.

`AddMessage(lv Level, msg interface{})` add message to the logger along with severity level.

`Print()` print the log message, then flush. When no message added, it will not print anything.

Example:
```go
logger := NewLogger("service_name")
	
logger.SetRequest(req) // req expected as *http.Request
logger.SetResponse(resp, string(b)) // resp expectes as *http.Response
logger.AddMessage(TraceLevel, "This is trace")
logger.AddMessage(DebugLevel, "This is debug")
logger.AddMessage(InfoLevel, "This is info")
logger.AddMessage(WarnLevel, "This is warning")
logger.AddMessage(ErrorLevel, errors.New("This is error"))
logger.AddMessage(FatalLevel, "This is fatal")
logger.AddMessage(PanicLevel, "This is panic")
logger.Print()
```

Output Example:
```json
{"level":"trace","log_id":"7715f1d9-e983-11e9-9b3e-d0c5d396697d","service":"test_service","log_message":"This is trace","request_body":{"Mock":{},"Error":null,"UseNetwork":false,"StatusCode":200,"Header":{"Content-Type":["application/json"],"Time":["2019-09-30T10:25:14+07:00"],"X-Test-Response":["Hello"]},"Cookies":null,"BodyBuffer":"eyJpZCI6ImExMjM0LWFiY2QiLCJxdW90ZXMiOiJFdmVyeXRoaW5nIHRoZSBsaWdodCB0b3VjaGVzLCBpcyBvdXIga2luZ2RvbSJ9Cg==","ResponseDelay":0,"Mappers":null,"Filters":null},"response_body":"{\"id\":\"a1234-abcd\",\"quotes\":\"Everything the light touches, is our kingdom\"}\n","response_headers":{"Content-Type":["application/json"],"Time":["2019-09-30T10:25:14+07:00"],"X-Test-Response":["Hello"]},"stack":[{"message":"This is trace","level":"trace","file":"/data/works/perkakas/log/log_test.go","func":"github.com/kitabisa/perkakas-log-poc/log.(*LogTestSuite).TestLog","line":63},{"message":"This is debug","level":"debug","file":"/data/works/perkakas/log/log_test.go","func":"github.com/kitabisa/perkakas-log-poc/log.(*LogTestSuite).TestLog","line":64},{"message":"This is info","level":"info","file":"/data/works/perkakas/log/log_test.go","func":"github.com/kitabisa/perkakas-log-poc/log.(*LogTestSuite).TestLog","line":65},{"message":"This is warning","level":"warning","file":"/data/works/perkakas/log/log_test.go","func":"github.com/kitabisa/perkakas-log-poc/log.(*LogTestSuite).TestLog","line":66},{"message":"This is error","level":"error","file":"/data/works/perkakas/log/log_test.go","func":"github.com/kitabisa/perkakas-log-poc/log.(*LogTestSuite).TestLog","line":67},{"message":"This is fatal","level":"fatal","file":"/data/works/perkakas/log/log_test.go","func":"github.com/kitabisa/perkakas-log-poc/log.(*LogTestSuite).TestLog","line":68},{"message":"This is panic","level":"panic","file":"/data/works/perkakas/log/log_test.go","func":"github.com/kitabisa/perkakas-log-poc/log.(*LogTestSuite).TestLog","line":69}]}
```