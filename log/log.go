package log

type Level uint32

const (
	FieldLogID           = "log_id"
	FieldHTTPStatus      = "http_status"
	FieldEndpoint        = "endpoint"
	FieldMethod          = "method"
	FieldServiceName     = "service"
	FieldUserID          = "user_id"
	FieldRequestBody     = "request_body"
	FieldRequestHeaders  = "request_headers"
	FieldResponseBody    = "response_body"
	FieldResponseHeaders = "response_headers"
)

type message struct {
	Message  interface{} `json:"message"`
	Level    Level       `json:"level"`
	File     string      `json:"file"`
	FuncName string      `json:"func"`
	Line     int         `json:"line"`
}
