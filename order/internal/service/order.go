package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rulanugrh/order/internal/entity"
	"github.com/rulanugrh/order/internal/grpc/order"
	"github.com/rulanugrh/order/internal/middleware"
	"github.com/rulanugrh/order/internal/repository"
	"github.com/rulanugrh/order/internal/util"
	"github.com/rulanugrh/order/pkg"
	"github.com/rulanugrh/order/pkg/logger"
)

type OrderServiceServer struct {
	order.UnimplementedOrderServiceServer
	repository repository.OrderInterface
	product    repository.ProductInterface
	xendit     pkg.XenditInterface
	rabbitmq   pkg.RabbitMQInterface
	metric     *pkg.Metrict
	log logger.ILogrus
}

func OrderService(repository repository.OrderInterface, product repository.ProductInterface, xendit pkg.XenditInterface, rabbitmq pkg.RabbitMQInterface, metric *pkg.Metrict, log logger.ILogrus) *OrderServiceServer {
	return &OrderServiceServer{repository: repository, product: product, xendit: xendit, rabbitmq: rabbitmq, metric: metric, log: log}
}

func (o *OrderServiceServer) Receipt(ctx context.Context, req *order.Request) (*order.ResponseProccess, error) {
	token, err := middleware.ReadToken()
	if err != nil {
		o.log.RecordGRPC("/order.OrderService/Receipt", "POST", 401).Error(err.Error())
		o.metric.Histogram.With(prometheus.Labels{"code": "401", "method": "POST", "type": "create", "service": "order"}).Observe(time.Since(time.Now()).Seconds())
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
		o.log.RecordGRPC("/order.OrderService/Receipt", "POST", 404).Error(find.Error())
		o.metric.Histogram.With(prometheus.Labels{"code": "404", "method": "POST", "type": "create", "service": "order"}).Observe(time.Since(time.Now()).Seconds())

		return util.NotFoundOrderCreate(find.Error()), find
	}

	result, err_create := o.repository.Create(input)
	if err_create != nil {
		o.log.RecordGRPC("/order.OrderService/Receipt", "POST", 400).Error(err_create.Error())
		o.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "POST", "type": "create", "service": "order"}).Observe(time.Since(time.Now()).Seconds())

		return util.BadRequestOrderCreate(err_create.Error()), err
	}

	response := order.Data{
		Uuid:          result.UUID,
		ProductName:   data.Name,
		Price:         int64(data.Price),
		MethodPayment: result.MethodPayment,
		Address:       result.Address,
		Fname:         token.Username,
	}

	marshalling, _ := json.Marshal(&response)

	err_publisher := o.rabbitmq.Publisher("order-create", marshalling, "order", "topic", token.Username)
	if err_publisher != nil {
		o.log.RecordGRPC("/order.OrderService/Receipt", "POST", 500).Error(err_publisher.Error())

		o.metric.Histogram.With(prometheus.Labels{"code": "500", "method": "POST", "type": "create", "service": "order"}).Observe(time.Since(time.Now()).Seconds())

		return util.InternalServerErrorOrderCreate(err_publisher.Error()), err_publisher
	}

	o.log.RecordGRPC("/order.OrderService/Receipt", "POST", 200).Info("success create receipt")
	o.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "POST", "type": "create", "service": "order"}).Observe(time.Since(time.Now()).Seconds())
	o.metric.Counter.With(prometheus.Labels{"type": "create", "service": "order"}).Inc()
	return util.SuccessOrderCreate("success create order", &response), nil
}

func (o *OrderServiceServer) Checkout(ctx context.Context, req *order.UUID) (*order.ResponseCheckout, error) {
	token, err := middleware.ReadToken()
	if err != nil {
		o.log.RecordGRPC("/order.OrderService/Checkout", "POST", 401).Error(err.Error())
		o.metric.Histogram.With(prometheus.Labels{"code": "401", "method": "POST", "type": "checkout", "service": "order"}).Observe(time.Since(time.Now()).Seconds())
		return util.UnauthorizedCheckout(err.Error()), err
	}

	data, err := o.repository.Checkout(req.Uuid)
	if err != nil {
		o.log.RecordGRPC("/order.OrderService/Checkout", "POST", 400).Error(err.Error())
		o.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "POST", "type": "checkout", "service": "order"}).Observe(time.Since(time.Now()).Seconds())

		return util.BadRequestOrderCheckout(err.Error()), err
	}

	product, err_product := o.product.FindID(data.ProductID)
	if err_product != nil {
		o.log.RecordGRPC("/order.OrderService/Checkout", "POST", 404).Error(err_product.Error())
		o.metric.Histogram.With(prometheus.Labels{"code": "404", "method": "POST", "type": "checkout", "service": "order"}).Observe(time.Since(time.Now()).Seconds())

		return util.BadRequestOrderCheckout(err_product.Error()), err_product
	}

	payment, err_payment := o.xendit.PaymentRequest(*data, token.Username, product.Name, product.Description, float64(product.Price))
	if err_payment != nil {
		o.log.RecordGRPC("/order.OrderService/Checkout", "POST", 400).Error(err_payment.Error())

		o.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "POST", "type": "payment", "service": "order"}).Observe(time.Since(time.Now()).Seconds())

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
		PaymentCreated:         payment.Created,
		PaymentUpdated:         payment.Updated,
		PaymentRequestCurrency: payment.Currency.String(),
		PaymentID:              payment.Id,
		Amount:                 *payment.Amount,
		ReferenceID:            payment.ReferenceId,
	}

	err_save := o.repository.SaveTransaction(transaction)
	if err_save != nil {
		o.log.RecordGRPC("/order.OrderService/Checkout", "POST", 500).Error(err_save.Error())
		o.metric.Histogram.With(prometheus.Labels{"code": "500", "method": "POST", "type": "payment", "service": "order"}).Observe(time.Since(time.Now()).Seconds())

		log.Println("[*] Error saving RecordGRPC transaction into DB: ", err_save)
	}

	marshalling, _ := json.Marshal(&response)

	err_publisher := o.rabbitmq.Publisher("order-checkout", marshalling, "order", "topic", token.Username)
	if err_publisher != nil {
		o.log.RecordGRPC("/order.OrderService/Checkout", "POST", 400).Error(err_publisher.Error())

		o.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "POST", "type": "checkout", "service": "order"}).Observe(time.Since(time.Now()).Seconds())

		return util.BadRequestOrderCheckout(err_publisher.Error()), err_publisher
	}

	o.log.RecordGRPC("/order.OrderService/Checkout", "POST", 200).Info("success checkout payment")
	o.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "POST", "type": "checkout", "service": "order"}).Observe(time.Since(time.Now()).Seconds())
	o.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "POST", "type": "payment", "service": "order"}).Observe(time.Since(time.Now()).Seconds())
	o.metric.Counter.With(prometheus.Labels{"type": "checkout", "service": "order"}).Inc()
	return util.SuccessOrderCheckout("success checkout order", &response), nil

}
