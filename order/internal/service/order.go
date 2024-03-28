package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rulanugrh/order/internal/entity"
	"github.com/rulanugrh/order/internal/grpc/order"
	"github.com/rulanugrh/order/internal/middleware"
	"github.com/rulanugrh/order/internal/repository"
	"github.com/rulanugrh/order/internal/util"
	"github.com/rulanugrh/order/pkg"
)

type OrderServiceServer struct {
	order.UnimplementedOrderServiceServer
	repository repository.OrderInterface
	product    repository.ProductInterface
	xendit     pkg.XenditInterface
	rabbitmq pkg.RabbitMQInterface
}

func OrderService(repository repository.OrderInterface, product repository.ProductInterface, xendit pkg.XenditInterface, rabbitmq pkg.RabbitMQInterface) *OrderServiceServer {
	return &OrderServiceServer{repository: repository, product: product, xendit: xendit, rabbitmq: rabbitmq}
}

func (o *OrderServiceServer) Receipt(ctx context.Context, req *order.Request) (*order.ResponseProccess, error) {
	token, err := middleware.ReadToken()
	if err != nil {
		return util.UnauthorizedCreateOrder(err.Error()), err
	}

	input := entity.Order{
		UserID:         token.ID,
		ProductID:      uint(req.Req.GetProductId()),
		Quantity:       uint(req.Req.GetQuantity()),
		MethodPayment:  req.Req.GetMethodPayment(),
		Address:        req.Req.GetAddress(),
		RequestCurreny: req.Req.RequstCurrency,
		MobilePhone:    req.Req.MobilePhone,
		ChannelCode:    req.Req.ChannelCode,
	}

	data, find := o.product.FindID(uint(req.Req.ProductId))
	if find != nil {
		return util.NotFoundOrderCreate(find.Error()), find
	}

	result, err := o.repository.Create(input)
	if err != nil {
		return util.BadRequestOrderCreate(err.Error()), err
	}

	response := order.Data{
		Uuid:          result.UUID,
		ProductName:   data.Name,
		Price:         int64(data.Price),
		MethodPayment: result.MethodPayment,
		Address:       result.Address,
		Fname: token.Username,
	}

	marshalling, _ := json.Marshal(&response)

	err_publisher := o.rabbitmq.Publisher("order-create", marshalling, "order", "topic", token.Username)
	if err_publisher != nil {
		return util.InternalServerErrorOrderCreate(err_publisher.Error()), err_publisher
	}

	return util.SuccessOrderCreate("success create order", &response), nil
}

func (o *OrderServiceServer) Checkout(ctx context.Context, req *order.UUID) (*order.ResponseCheckout, error) {
	token, err := middleware.ReadToken()
	if err != nil {
		return util.UnauthorizedCheckout(err.Error()), err
	}

	data, err := o.repository.Checkout(req.Uuid)
	if err != nil {
		return util.BadRequestOrderCheckout(err.Error()), err
	}

	product, err_product := o.product.FindID(data.ProductID)
	if err_product != nil {
		return util.BadRequestOrderCheckout(err_product.Error()), err_product
	}

	payment, err_payment := o.xendit.PaymentRequest(*data, token.Username, product.Name, product.Description, float64(product.Price))
	if err_payment != nil {
		return util.BadRequestOrderCheckout(err_payment.Error()), err_payment
	}

	response := order.CheckOut{
		ProductName: product.Name,
		Price:       int64(product.Price),
		Quantity:    int64(data.Quantity),
		Total:       (int64(data.Quantity) * int64(product.Price)),
		LinkPayment: payment.GetCreated(),
	}

	transaction := entity.Transaction{
		OrderID:                data.ID,
		OrderUUID:              data.UUID,
		MethodPayment:          data.MethodPayment,
		Status:                 string(payment.Status),
		PaymentCreated:         payment.Created,
		PaymentUpdated:         payment.Updated,
		PaymentRequestCurrency: payment.Currency.String(),
		PaymentID:              payment.Id,
		Amount:                 *payment.Amount,
		ReferenceID:            payment.ReferenceId,
	}

	err_save := o.repository.SaveTransaction(transaction)
	if err_save != nil {
		log.Println("[*] Error saving record transaction into DB: ", err_save)
	}

	marshalling, _ := json.Marshal(&response)

	err_publisher := o.rabbitmq.Publisher("order-checkout", marshalling, "order", "topic", token.Username)
	if err_publisher != nil {
		return util.BadRequestOrderCheckout(err_publisher.Error()), err_publisher
	}

	return util.SuccessOrderCheckout("success checkout order", &response), nil

}
