package constant

import "fmt"

type Response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (r Response) Error() string {
	return r.Message
}

func Success(msg string) Response {
	return Response{
		Code:    200,
		Message: msg,
	}
}

func Created(msg string) Response {
	return Response{
		Code:    201,
		Message: msg,
	}
}

func NotFound(msg string) Response {
	return Response{
		Code:    404,
		Message: msg,
	}
}

func BadRequest(msg string, err error) Response {
	return Response{
		Code:    400,
		Message: fmt.Sprintf("%s: %s", msg, err.Error()),
	}
}

func InternalServerError(msg string, err error) Response {
	return Response{
		Code:    500,
		Message: fmt.Sprintf("%s: %s", msg, err.Error()),
	}
}

func OverloadRequest(msg string) Response {
	return Response{
		Code:    429,
		Message: msg,
	}
}