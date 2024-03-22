package web

type Response struct {
  Code int16 `json:"code"`
  Message string `json:"message"`
  Data interface{} `json:"data"`
}

func (r Response) Error() string {
  return r.Message
}

func Created(data interface{}, msg string) Response {
  return Response{
    Code: 201,
    Message: msg,
    Data: data,
  }
}

func Success(data interface{}, msg string) Response {
  return Response{
    Code: 200,
    Message: msg,
    Data: data,
  }
}

func NotFound(msg string) Response {
  return Response{
    Code: 404,
    Message: msg,
  }
}

func BadRequest(msg string) Response {
  return Response{
    Code: 400,
    Message: msg,
  }
}

func Forbidden(msg string) Response {
  return Response{
    Code: 403,
    Message: msg,
  }
}

func Unauthorized(msg string) Response {
  return Response{
    Code: 401,
    Message: msg,
  }
}

