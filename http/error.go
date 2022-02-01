package http

import (
	"errors"
	"net/http"
)

var ErrUnknown = &ErrorResponse{
	ErrorData: ErrorData{
		Message: "Something went wrong",
		Code:    "KTBS_ERROR_319999",
	},
	ErrorType:  errors.New("unknown_error"),
	HttpStatus: http.StatusInternalServerError,
}

var ErrUnauthorized = &ErrorResponse{
	ErrorData: ErrorData{
		Message: "You are not authorized",
		Code:    "KTBS_ERROR_310001",
	},
	ErrorType:  errors.New("authorization_issue"),
	HttpStatus: http.StatusInternalServerError,
}

var ErrInvalidHeader = &ErrorResponse{
	ErrorData: ErrorData{
		Message: "Invalid or incomplete header",
		Code:    "KTBS_ERROR_310002",
	},
	ErrorType:  errors.New("request_parameter_issue"),
	HttpStatus: http.StatusBadRequest,
}

var ErrInvalidHeaderSignature = &ErrorResponse{
	ErrorData: ErrorData{
		Message: "Signature header is invalid",
		Code:    "KTBS_ERROR_319003",
	},
	ErrorType:  errors.New("signature_issue"),
	HttpStatus: http.StatusBadRequest,
}

var ErrInvalidHeaderTime = &ErrorResponse{
	ErrorData: ErrorData{
		Message: "Header time already expired",
		Code:    "KTBS_ERROR_319004",
	},
	ErrorType:  errors.New("request_issue"),
	HttpStatus: http.StatusBadRequest,
}
