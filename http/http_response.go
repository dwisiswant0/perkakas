package http

// Meta defines meta format format for api format
type Meta struct {
	Version string `json:"version" mapstructure:"version"`
	Status  string `json:"api_status" mapstructure:"api_status"`
	APIEnv  string `json:"api_env" mapstructure:"api_env"`
}

type SuccessResponse struct {
	SuccessData
	Next *string `json:"next,omitempty" mapstructure:"next,omitempty"`
	Meta Meta    `json:"meta" mapstructure:"meta"`
}

type SuccessData struct {
	Status  string      `json:"status" mapstructure:"status"`
	Message string      `json:"message" mapstructure:"message"`
	Data    interface{} `json:"data,omitempty" mapstructure:"data,omitempty"`
}

type ErrorResponse struct {
	ErrorData  `json:"error" mapstructure:"error"`
	ErrorType  error `json:"-"`
	HttpStatus int   `json:"-"`
}

type ErrorData struct {
	Message string `json:"message"`
	Code    string `json:"error_code"`
}

func (e *ErrorResponse) Error() string {
	return e.ErrorType.Error()
}
