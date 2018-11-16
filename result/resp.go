package result

import (
	"fmt"
)

func Success(data interface{}) ApiResponse {
	return ApiResponse{100200, Resp[100200], data}
}

func ErrCode(code int) ApiResponse {
	response := ApiResponse{Code: code, Data: ""}
	if message, ok := Resp[code]; ok {
		response.Msg = message
	} else {
		response.Msg = "unknown error"
	}
	return response
}

func ErrMsg(msg string) ApiResponse {
	return ApiResponse{Code: 100102, Msg: msg, Data: ""}
}

type ApiResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (e ApiResponse) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Msg)
}
