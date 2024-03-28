package pkg

import (
	"context"
	"strconv"
	"time"

	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/entity"
	"github.com/rulanugrh/order/internal/util/constant"
	xdt "github.com/xendit/xendit-go/v4"
	"github.com/xendit/xendit-go/v4/payment_request"
)

type XenditInterface interface {
	PaymentRequest(req entity.Order, username string, product_name string, product_desc string, product_price float64) (*payment_request.PaymentRequest, error)
}

type xendit struct {
	client *xdt.APIClient
	conf   *config.App
}

func XenditPluggin(client *xdt.APIClient, conf *config.App) XenditInterface {
	return &xendit{client: client, conf: conf}
}

func (x *xendit) PaymentRequest(req entity.Order, username string, product_name string, product_desc string, product_price float64) (*payment_request.PaymentRequest, error) {
	productID := strconv.Itoa(int(req.ProductID))
	pay := payment_request.PaymentRequestParameters{
		ReferenceId: &req.UUID,
		Currency:    payment_request.PaymentRequestCurrency(req.RequestCurreny),
		Items: []payment_request.PaymentRequestBasketItem{
			{
				Name:        product_name,
				Description: &product_desc,
				Quantity:    float64(req.Quantity),
				Price:       float64(product_price),
				Currency:    req.RequestCurreny,
				ReferenceId: &productID,
			},
		},
		Customer: map[string]interface{}{
			"id":      req.UserID,
			"address": req.Address,
		},
		Metadata: map[string]interface{}{
			"type":     "pay",
			"product":  product_name,
			"quantity": req.Quantity,
		},
		PaymentMethod: x.paymentMethod(req.UUID, req.MethodPayment, req.ChannelCode, req.MobilePhone, username),
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

func (x *xendit) paymentMethod(uuid string, method string, channel_code string, mobile_phone string, customer string) (request *payment_request.PaymentMethodParameters) {
	if method == "EWALLET" {
		request = &payment_request.PaymentMethodParameters{
			Type:        payment_request.PaymentMethodType(method),
			Reusability: payment_request.PAYMENTMETHODREUSABILITY_ONE_TIME_USE,
			ReferenceId: &uuid,
			Ewallet:     *payment_request.NewNullableEWalletParameters(x.ewallet(mobile_phone, channel_code)),
		}
	} else if method == "QRIS" {
		request = &payment_request.PaymentMethodParameters{
			Type:        payment_request.PaymentMethodType(method),
			Reusability: payment_request.PAYMENTMETHODREUSABILITY_ONE_TIME_USE,
			ReferenceId: &uuid,
			QrCode:      *payment_request.NewNullableQRCodeParameters(x.qrcode(channel_code)),
		}
	} else if method == "VA" {
		request = &payment_request.PaymentMethodParameters{
			Type:           payment_request.PaymentMethodType(method),
			Reusability:    payment_request.PAYMENTMETHODREUSABILITY_ONE_TIME_USE,
			ReferenceId:    &uuid,
			VirtualAccount: *payment_request.NewNullableVirtualAccountParameters(x.virtual_account(channel_code, customer)),
		}
	}

	return request
}

func (x *xendit) ewallet(mobile_phone string, channel_code string) (ewallet *payment_request.EWalletParameters) {

	if channel_code == "DANA" {
		ewallet = &payment_request.EWalletParameters{
			ChannelCode: payment_request.EWALLETCHANNELCODE_DANA.Ptr(),
			ChannelProperties: &payment_request.EWalletChannelProperties{
				SuccessReturnUrl: &x.conf.Xendit.SuccessURL,
				FailureReturnUrl: &x.conf.Xendit.FailureURL,
				CancelReturnUrl:  &x.conf.Xendit.CancelURL,
				MobileNumber:     &mobile_phone,
			},
		}
	} else if channel_code == "GOPAY" {
		ewallet = &payment_request.EWalletParameters{
			ChannelCode: payment_request.EWALLETCHANNELCODE_GCASH.Ptr(),
			ChannelProperties: &payment_request.EWalletChannelProperties{
				SuccessReturnUrl: &x.conf.Xendit.SuccessURL,
				FailureReturnUrl: &x.conf.Xendit.FailureURL,
				CancelReturnUrl:  &x.conf.Xendit.CancelURL,
				MobileNumber:     &mobile_phone,
			},
		}
	} else if channel_code == "OVO" {
		ewallet = &payment_request.EWalletParameters{
			ChannelCode: payment_request.EWALLETCHANNELCODE_OVO.Ptr(),
			ChannelProperties: &payment_request.EWalletChannelProperties{
				SuccessReturnUrl: &x.conf.Xendit.SuccessURL,
				FailureReturnUrl: &x.conf.Xendit.FailureURL,
				CancelReturnUrl:  &x.conf.Xendit.CancelURL,
				MobileNumber:     &mobile_phone,
			},
		}
	} else if channel_code == "SHOPEE" {
		ewallet = &payment_request.EWalletParameters{
			ChannelCode: payment_request.EWALLETCHANNELCODE_SHOPEEPAY.Ptr(),
			ChannelProperties: &payment_request.EWalletChannelProperties{
				SuccessReturnUrl: &x.conf.Xendit.SuccessURL,
				FailureReturnUrl: &x.conf.Xendit.FailureURL,
				CancelReturnUrl:  &x.conf.Xendit.CancelURL,
				MobileNumber:     &mobile_phone,
			},
		}
	}

	return ewallet
}

func (x *xendit) qrcode(channel_code string) (response *payment_request.QRCodeParameters) {
	var expire time.Time = time.Now().Add(10 * time.Minute)
	if channel_code == "QRIS" {
		response = &payment_request.QRCodeParameters{
			ChannelCode: *payment_request.NewNullableQRCodeChannelCode(payment_request.QRCODECHANNELCODE_QRIS.Ptr()),
			ChannelProperties: &payment_request.QRCodeChannelProperties{
				QrString:  &payment_request.NewCaptureWithDefaults().PaymentRequestId,
				ExpiresAt: &expire,
			},
		}
	} else if channel_code == "DANA" {
		response = &payment_request.QRCodeParameters{
			ChannelCode: *payment_request.NewNullableQRCodeChannelCode(payment_request.QRCODECHANNELCODE_DANA.Ptr()),
			ChannelProperties: &payment_request.QRCodeChannelProperties{
				QrString:  &payment_request.NewCaptureWithDefaults().PaymentRequestId,
				ExpiresAt: &expire,
			},
		}
	}

	return response
}

func (x *xendit) virtual_account(channel_code string, customer string) (response *payment_request.VirtualAccountParameters) {
	var expire time.Time = time.Now().Add(25 * time.Minute)
	if channel_code == "BRI" {
		response = &payment_request.VirtualAccountParameters{
			ChannelCode: payment_request.VIRTUALACCOUNTCHANNELCODE_BRI,
			ChannelProperties: payment_request.VirtualAccountChannelProperties{
				CustomerName: customer,
				ExpiresAt:    &expire,
			},
		}
	} else if channel_code == "BCA" {
		response = &payment_request.VirtualAccountParameters{
			ChannelCode: payment_request.VIRTUALACCOUNTCHANNELCODE_BCA,
			ChannelProperties: payment_request.VirtualAccountChannelProperties{
				CustomerName: customer,
				ExpiresAt:    &expire,
			},
		}
	} else if channel_code == "BNI" {
		response = &payment_request.VirtualAccountParameters{
			ChannelCode: payment_request.VIRTUALACCOUNTCHANNELCODE_BNI,
			ChannelProperties: payment_request.VirtualAccountChannelProperties{
				CustomerName: customer,
				ExpiresAt:    &expire,
			},
		}
	} else if channel_code == "MANDIRI" {
		response = &payment_request.VirtualAccountParameters{
			ChannelCode: payment_request.VIRTUALACCOUNTCHANNELCODE_MANDIRI,
			ChannelProperties: payment_request.VirtualAccountChannelProperties{
				CustomerName: customer,
				ExpiresAt:    &expire,
			},
		}
	} else if channel_code == "PERMATA" {
		response = &payment_request.VirtualAccountParameters{
			ChannelCode: payment_request.VIRTUALACCOUNTCHANNELCODE_PERMATA,
			ChannelProperties: payment_request.VirtualAccountChannelProperties{
				CustomerName: customer,
				ExpiresAt:    &expire,
			},
		}
	} else if channel_code == "BJB" {
		response = &payment_request.VirtualAccountParameters{
			ChannelCode: payment_request.VIRTUALACCOUNTCHANNELCODE_BJB,
			ChannelProperties: payment_request.VirtualAccountChannelProperties{
				CustomerName: customer,
				ExpiresAt:    &expire,
			},
		}
	} else if channel_code == "BSI" {
		response = &payment_request.VirtualAccountParameters{
			ChannelCode: payment_request.VIRTUALACCOUNTCHANNELCODE_BSI,
			ChannelProperties: payment_request.VirtualAccountChannelProperties{
				CustomerName: customer,
				ExpiresAt:    &expire,
			},
		}
	}

	return response
}
