package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

type HttpHandlerContext struct {
	M Meta
	E map[error]*ErrorResponse
}

func NewContextHandler(meta Meta) HttpHandlerContext {
	var errMap map[error]*ErrorResponse = map[error]*ErrorResponse{
		// register general error here, so if there are new general error you must add it here
		ErrInvalidHeader:          ErrInvalidHeader,
		ErrUnauthorized:           ErrUnauthorized,
		ErrInvalidHeaderSignature: ErrInvalidHeaderSignature,
		ErrInvalidHeaderTime:      ErrInvalidHeaderTime,
	}

	return HttpHandlerContext{
		M: meta,
		E: errMap,
	}
}

func (hctx HttpHandlerContext) AddError(key error, value *ErrorResponse) {
	hctx.E[key] = value
}

func (hctx HttpHandlerContext) AddErrorMap(errMap map[error]*ErrorResponse) {
	for k, v := range errMap {
		hctx.E[k] = v
	}
}

type CustomWriter struct {
	C HttpHandlerContext
}

func (c *CustomWriter) Write(w http.ResponseWriter, data interface{}, nextPage *string) {
	var successResp SuccessResponse
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

	successResp.ResponseCode = "000000"
	successResp.Next = nextPage
	successResp.Meta = c.C.M

	writeSuccessResponse(w, successResp)
}

// WriteError sending error response based on err type
func (c *CustomWriter) WriteError(w http.ResponseWriter, err error) {
	if len(c.C.E) > 0 {
		errorResponse := LookupError(c.C.E, err)
		if errorResponse == nil {
			errorResponse = ErrUnknown
		}

		errorResponse.Meta = c.C.M
		writeErrorResponse(w, errorResponse)
	} else {
		var errorResponse *ErrorResponse = &ErrorResponse{}
		if errors.As(err, &errorResponse) {
			errorResponse.Meta = c.C.M
			writeErrorResponse(w, errorResponse)
		} else {
			errorResponse = ErrUnknown
			errorResponse.Meta = c.C.M
			writeErrorResponse(w, errorResponse)
		}
	}
}

func writeResponse(w http.ResponseWriter, response interface{}, contentType string, httpStatus int) {
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to unmarshal"))
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(httpStatus)
	w.Write(res)
}

func writeSuccessResponse(w http.ResponseWriter, response SuccessResponse) {
	writeResponse(w, response, "application/json", http.StatusOK)
}

func writeErrorResponse(w http.ResponseWriter, errorResponse *ErrorResponse) {
	writeResponse(w, errorResponse, "application/json", errorResponse.HttpStatus)
}

// LookupError will get error message based on error type, with variables if you want give dynamic message error
func LookupError(lookup map[error]*ErrorResponse, err error) (res *ErrorResponse) {
	if msg, ok := lookup[err]; ok {
		res = msg
	}

	return
}
