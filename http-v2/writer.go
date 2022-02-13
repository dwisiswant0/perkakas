package http

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type HttpHandlerContext struct {
	M Meta
	E map[string]*ErrorResponse
}

//NewContextHandler add base error response
func NewContextHandler(meta Meta) HttpHandlerContext {
	var errMap = map[string]*ErrorResponse{
		// register general error here, so if there are new general error you must add it here
		ErrInvalidHeader.ErrorType.Error():          ErrInvalidHeader,
		ErrUnauthorized.ErrorType.Error():           ErrUnauthorized,
		ErrInvalidHeaderSignature.ErrorType.Error(): ErrInvalidHeaderSignature,
		ErrInvalidHeaderTime.ErrorType.Error():      ErrInvalidHeaderTime,
	}

	return HttpHandlerContext{
		M: meta,
		E: errMap,
	}
}

// AddError map error string with the error response
func (hctx HttpHandlerContext) AddError(key error, value *ErrorResponse) {
	hctx.E[key.Error()] = value
}

// AddErrorMap populate error response from error map
func (hctx HttpHandlerContext) AddErrorMap(errMap map[error]*ErrorResponse) {
	for k, v := range errMap {
		hctx.E[k.Error()] = v
	}
}

type CustomWriter struct {
	C HttpHandlerContext
}

func (c *CustomWriter) Write(w http.ResponseWriter, data interface{}, nextPage *string) {
	var successResp SuccessResponse
	c.setSuccessRespData(data, &successResp)

	successResp.SuccessData.Status = "200"
	successResp.SuccessData.Message = "OK"
	successResp.Next = nextPage
	successResp.Meta = c.C.M

	writeSuccessResponse(w, successResp)
}

//TODO add explanation detail of what this block code is doing
func (c *CustomWriter) setSuccessRespData(data interface{}, successResp *SuccessResponse) {
	voData := reflect.ValueOf(data)
	arrayData := []interface{}{}

	if voData.Kind() != reflect.Slice {
		if voData.IsValid() {
			arrayData = []interface{}{data}
		}
		successResp.Data = arrayData
	} else {
		if voData.Len() != 0 {
			successResp.Data = data
		} else {
			successResp.Data = arrayData
		}
	}

}

// WriteError sending error response based on err type
func (c *CustomWriter) WriteError(w http.ResponseWriter, err error) {
	if len(c.C.E) > 0 {
		errorResponse := LookupError(c.C.E, err)
		if errorResponse == nil {
			errorResponse = ErrUnknown
		}

		writeErrorResponse(w, errorResponse)
	} else {
		var errorResponse = &ErrorResponse{
			ErrorType: err,
		}
		writeErrorResponse(w, errorResponse)
	}
}

func writeJSONResponse(w http.ResponseWriter, response interface{}, httpStatus int) {
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to unmarshal"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(res)
}

func writeSuccessResponse(w http.ResponseWriter, response SuccessResponse) {
	writeJSONResponse(w, response, http.StatusOK)
}

func writeErrorResponse(w http.ResponseWriter, errorResponse *ErrorResponse) {
	writeJSONResponse(w, errorResponse, errorResponse.HttpStatus)
}

// LookupError will get error message based on error type, with variables if you want give dynamic message error
func LookupError(lookup map[string]*ErrorResponse, err error) (res *ErrorResponse) {
	if msg, ok := lookup[err.Error()]; ok {
		res = msg
	}

	return
}
