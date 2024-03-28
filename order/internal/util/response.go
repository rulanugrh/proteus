package util

import (
	"github.com/rulanugrh/order/internal/grpc/cart"
	"github.com/rulanugrh/order/internal/grpc/order"
)

func BadRequestOrderCreate(msg string) *order.ResponseProccess {
	return &order.ResponseProccess{
		Code: 400,
		Message: msg,
	}
}

func InternalServerErrorOrderCreate(msg string) *order.ResponseProccess {
	return &order.ResponseProccess{
		Code: 500,
		Message: msg,
	}
}

func NotFoundOrderCreate(msg string) *order.ResponseProccess {
	return &order.ResponseProccess{
		Code: 404,
		Message: msg,
	}
}

func UnauthorizedCreateOrder(msg string) *order.ResponseProccess {
	return &order.ResponseProccess{
		Code: 401,
		Message: msg,
	}
}

func SuccessOrderCreate(msg string, data *order.Data) *order.ResponseProccess {
	return &order.ResponseProccess{
		Code: 201,
		Message: msg,
		Data: data,
	}
}

func BadRequestOrderCheckout(msg string) *order.ResponseCheckout {
	return &order.ResponseCheckout{
		Code: 400,
		Message: msg,
	}
}

func UnauthorizedCheckout(msg string) *order.ResponseCheckout {
	return &order.ResponseCheckout{
		Code: 401,
		Message: msg,
	}
}

func SuccessOrderCheckout(msg string, data *order.CheckOut) *order.ResponseCheckout {
	return &order.ResponseCheckout{
		Code: 200,
		Message: msg,
		Data: data,
	}
}

func SuccessCart(msg string) *cart.Response {
	return &cart.Response{
		Code: 200,
		Message: msg,
	}
}

func NotFoundCart(msg string) *cart.Response {
	return &cart.Response{
		Code: 400,
		Message: msg,
	}
}

func InternalServerErrorCart(msg string) *cart.Response {
	return &cart.Response{
		Code: 500,
		Message: msg,
	}
}

func UnauthorizedCart(msg string) *cart.Response {
	return &cart.Response{
		Code: 401,
		Message: msg,
	}
}

func DeletedCart(msg string) *cart.Response {
	return &cart.Response{
		Code: 204,
		Message: msg,
	}
}

func BadRequestCart(msg string) *cart.Response {
	return &cart.Response{
		Code: 400,
		Message: msg,
	}
}

func SuccessCreatedCart(msg string, data *cart.Data) *cart.Created {
	return &cart.Created{
		Code: 201,
		Message: msg,
		Data: data,
	}
}

func FailureCreatedCart(msg string) *cart.Created {
	return &cart.Created{
		Code: 400,
		Message: msg,
	}
}