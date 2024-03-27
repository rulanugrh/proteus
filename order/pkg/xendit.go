package pkg

import (
	"context"
	"strconv"

	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/entity"
	"github.com/rulanugrh/order/internal/util/constant"
	xdt "github.com/xendit/xendit-go/v4"
	"github.com/xendit/xendit-go/v4/payment_request"
)

type XenditInterface interface {
	PaymentRequest(req entity.Order) (*payment_request.PaymentRequest, error)
}

type xendit struct {
	client *xdt.APIClient
	conf   *config.App
}

func XenditPluggin(client *xdt.APIClient, conf *config.App) XenditInterface {
	return &xendit{client: client, conf: conf}
}

func (x *xendit) PaymentRequest(req entity.Order) (*payment_request.PaymentRequest, error) {
	productID := strconv.Itoa(int(req.ProductID))
	pay := payment_request.PaymentRequestParameters{
		ReferenceId: &req.UUID,
		Currency:    payment_request.PaymentRequestCurrency(req.RequestCurreny),
		Items: append([]payment_request.PaymentRequestBasketItem{},
			payment_request.PaymentRequestBasketItem{
				Name:        req.Product.Name,
				Description: &req.Product.Description,
				Quantity:    float64(req.Quantity),
				Price:       float64(req.Product.Price),
				Currency:    req.RequestCurreny,
				ReferenceId: &productID,
			},
		),
		Customer: map[string]interface{}{
			"id": req.UserID,
			"address": req.Address,
		},
		Metadata: map[string]interface{}{
			"type":     "pay",
			"product":  req.Product.Name,
			"quantity": req.Quantity,
		},
		PaymentMethod: &payment_request.PaymentMethodParameters{
			Type:        payment_request.PaymentMethodType(req.MethodPayment),
			ReferenceId: &productID,
			Reusability: payment_request.PaymentMethodReusability("ONE_TIME_USE"),
		},
		ChannelProperties: &payment_request.PaymentRequestParametersChannelProperties{
			FailureReturnUrl: &x.conf.Xendit.FailureURL,
			SuccessReturnUrl: &x.conf.Xendit.SuccessURL,
			CancelReturnUrl:  &x.conf.Xendit.CancelURL,
		},
	}

	response, r, err := x.client.PaymentRequestApi.CreatePaymentRequest(context.Background()).
		IdempotencyKey(req.UUID).
		ForUserId(strconv.Itoa(int(req.UserID))).
		PaymentRequestParameters(pay).
		Execute()

	if r.StatusCode == 500 {
		return nil, constant.InternalServerError("sorry internal server error", err)
	} else if r.StatusCode == 400 {
		return nil, constant.BadRequest("bad request to xendit", err)
	} else if r.StatusCode == 429 {
		return nil, constant.OverloadRequest("too many request to xendit")
	} else {
		return response, nil
	}

}
