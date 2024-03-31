package helper

import "encoding/json"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}

func (r Response) Error() string {
	return r.Message
}

func Success(msg string, status string, data interface{}) Response {
	return Response{
		Code:    200,
		Message: msg,
		Status:  status,
		Data: data,
	}
}

func BadRequest(msg string, status string) Response {
	return Response{
		Code:    400,
		Message: msg,
		Status:  status,
	}
}

func Marshal(data any) []byte {
	marshal, _ := json.Marshal(data)

	return marshal
}
