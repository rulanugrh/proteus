package service

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rulanugrh/order/internal/entity"
	"github.com/rulanugrh/order/internal/grpc/cart"
	"github.com/rulanugrh/order/internal/middleware"
	"github.com/rulanugrh/order/internal/repository"
	"github.com/rulanugrh/order/internal/util"
	"github.com/rulanugrh/order/internal/util/constant"
	"github.com/rulanugrh/order/pkg"
	"github.com/rulanugrh/order/pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CartServiceServer struct {
	cart.UnimplementedCartServiceServer
	repository repository.CartInterface
	product    repository.ProductInterface
	metric     *pkg.Metrict
	log        logger.ILogrus
}

func CartService(repository repository.CartInterface, product repository.ProductInterface, metric *pkg.Metrict, log logger.ILogrus) *CartServiceServer {
	return &CartServiceServer{repository: repository, product: product, metric: metric, log: log}
}

func (c *CartServiceServer) AddToCart(ctx context.Context, req *cart.Request) (*cart.Response, error) {
	_, err := c.product.FindID(uint(req.Req.ProductId))
	if err != nil {
		c.log.RecordGRPC("/cart.CartService/AddToCart", "POST", 404).Error(err.Error())
		c.metric.Histogram.With(prometheus.Labels{"code": "404", "method": "POST", "type": "add", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())
		return util.NotFoundCart(err.Error()), err
	}

	token, err_token := middleware.ReadToken()
	if err_token != nil {
		c.log.RecordGRPC("/cart.CartService/AddToCart", "POST", 401).Error(err_token.Error())
		c.metric.Histogram.With(prometheus.Labels{"code": "401", "method": "POST", "type": "add", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())
		return util.UnauthorizedCart(err_token.Error()), err
	}

	input := entity.Cart{
		Quantity:  uint(req.Req.GetQuantity()),
		ProductID: uint(req.Req.GetProductId()),
		UserID:    token.ID,
	}

	err_create := c.repository.AddToCart(input)
	if err_create != nil {
		c.log.RecordGRPC("/cart.CartService/AddToCart", "POST", 400).Error(err_create.Error())
		c.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "POST", "type": "add", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())
		return util.BadRequestCart(err_create.Error()), err_create
	}

	c.log.RecordGRPC("/cart.CartService/AddToCart", "POST", 400).Info("success add to cart")
	c.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "POST", "type": "add", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())
	c.metric.Counter.With(prometheus.Labels{"type": "add", "service": "cart"}).Inc()
	return util.SuccessCart("success add to cart"), nil
}

func (c *CartServiceServer) DeleteCart(ctx context.Context, req *cart.ID) (*cart.Response, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		c.log.RecordGRPC("/cart.CartService/DeleteCart", "DELETE", 500).Error("error converting request ID")
		c.metric.Histogram.With(prometheus.Labels{"code": "500", "method": "DELETE", "type": "delete", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())

		return util.InternalServerErrorCart(err.Error()), err
	}

	err_delete := c.repository.Delete(uint(id))
	if err_delete != nil {
		c.log.RecordGRPC("/cart.CartService/DeleteCart", "DELETE", 400).Error(err_delete.Error())
		c.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "DELETE", "type": "delete", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())

		return util.BadRequestCart(err_delete.Error()), err_delete
	}

	c.log.RecordGRPC("/cart.CartService/DeleteCart", "DELETE", 200).Info("success deleted cart")
	c.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "DELETE", "type": "delete", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())

	return util.DeletedCart("success delete cart by this id"), nil
}

func (c *CartServiceServer) ListCart(empty *emptypb.Empty, stream cart.CartService_ListCartServer) error {
	token, err := middleware.ReadToken()
	if err != nil {
		c.log.RecordGRPC("/cart.CartService/ListCart", "GET", 401).Error(err.Error())
		c.metric.Histogram.With(prometheus.Labels{"code": "401", "method": "GET", "type": "getAll", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())

		return constant.Unauthorized(err.Error())
	}

	data, err_list := c.repository.ListCart(token.ID)
	if err_list != nil {
		c.log.RecordGRPC("/cart.CartService/ListCart", "GET", 400).Error(err_list.Error())
		c.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "GET", "type": "getAll", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())

		return constant.BadRequest(err_list.Error(), err_list)
	}

	for _, result := range *data {
		product, err := c.product.FindID(result.ID)
		if err != nil {
			constant.BadRequest(err.Error(), err)
		}

		stream.Send(&cart.CartList{
			ProductName:  product.Name,
			ProductDesc:  product.Description,
			ProductPrice: uint64(product.Price),
			Quantity:     int32(result.Quantity),
		})
	}

	c.log.RecordGRPC("/cart.CartService/ListCart", "GET", 200).Info("success get all cart")
	c.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "GET", "type": "getAll", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())

	return constant.Success("success get all product by this user id")
}

func (c *CartServiceServer) Proccesses(ctx context.Context, req *cart.RequestProcess) (*cart.Created, error) {
	input := entity.Updates{
		MethodType:     req.Req.GetMethodPayment(),
		Address:        req.Req.GetAddress(),
		ChannelCode:    req.Req.GetChannelCode(),
		RequestCurreny: req.Req.GetRequestCurrency(),
		MobilePhone:    req.Req.GetMobilePhone(),
	}

	data, err := c.repository.ProcessCart(uint(req.Id), input)
	if err != nil {
		c.log.RecordGRPC("/cart.CartService/Proccesses", "POST", 400).Error(err.Error())
		c.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "POST", "type": "proccess", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())

		return util.FailureCreatedCart(err.Error()), err
	}

	product, err_find := c.product.FindID(data.ProductID)
	if err_find != nil {
		c.log.RecordGRPC("/cart.CartService/Proccesses", "POST", 404).Error(err_find.Error())

		c.metric.Histogram.With(prometheus.Labels{"code": "404", "method": "POST", "type": "proccess", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())

		return util.NotFoundCartCreated(err_find.Error()), err_find
	}

	response := cart.Data{
		Uuid:        data.UUID,
		Quantity:    int64(data.Quantity),
		ProductName: product.Name,
		Price:       int64(product.Price),
		Total:       int64(product.Price * uint32(data.Quantity)),
	}

	c.log.RecordGRPC("/cart.CartService/Proccesses", "POST", 200).Info("success proccess cart to order")
	c.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "POST", "type": "proccess", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())
	c.metric.Counter.With(prometheus.Labels{"type": "process", "service": "cart"}).Inc()
	return util.SuccessCreatedCart("success created for order", &response), nil
}

func (c *CartServiceServer) Update(ctx context.Context, req *cart.RequestUpdate) (*cart.Response, error) {
	request := entity.Cart{
		Quantity: uint(req.Req.GetQuantity()),
	}

	err := c.repository.Update(uint(req.Id), request)
	if err != nil {
		c.log.RecordGRPC("/cart.CartService/Update", "PUT", 400).Error(err.Error())
		c.metric.Histogram.With(prometheus.Labels{"code": "400", "method": "PUT", "type": "update", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())

		return util.BadRequestCart(err.Error()), err
	}

	c.log.RecordGRPC("/cart.CartService/Update", "PUT", 200).Info("success update cart with this id")
	c.metric.Histogram.With(prometheus.Labels{"code": "200", "method": "PUT", "type": "update", "service": "cart"}).Observe(time.Since(time.Now()).Seconds())
	c.metric.Counter.With(prometheus.Labels{"type": "update", "service": "cart"}).Inc()
	return util.SuccessCart("success update cart"), nil
}
