package payload

import (
	"helloworld/shared/model/apperror"
)

type Response struct {
	Success      bool   `json:"success"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Data         any    `json:"data"`
	TraceID      string `json:"traceId"`
}

func NewSuccessResponse(data any, traceID string) any {
	var res Response
	res.Success = true
	res.Data = data
	res.TraceID = traceID
	return res
}

func NewErrorResponse(err error, traceID string) any {
	var res Response
	res.Success = false
	res.TraceID = traceID

	et, ok := err.(apperror.ErrorType)
	if !ok {
		res.ErrorCode = "UNDEFINED"
		res.ErrorMessage = err.Error()
		return res
	}

	res.ErrorCode = et.Code()
	res.ErrorMessage = et.Error()
	return res
}
