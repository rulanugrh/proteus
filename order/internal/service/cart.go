package service

import (
	"context"
	"strconv"

	"github.com/rulanugrh/order/internal/entity"
	"github.com/rulanugrh/order/internal/grpc/cart"
	"github.com/rulanugrh/order/internal/middleware"
	"github.com/rulanugrh/order/internal/repository"
	"github.com/rulanugrh/order/internal/util"
	"github.com/rulanugrh/order/internal/util/constant"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CartServiceServer struct {
	cart.UnimplementedCartServiceServer
	repository repository.CartInterface
	product    repository.ProductInterface
}

func CartService(repository repository.CartInterface, product repository.ProductInterface) *CartServiceServer {
	return &CartServiceServer{repository: repository, product: product}
}

func (c *CartServiceServer) AddToCart(ctx context.Context, req *cart.Request) (*cart.Response, error) {
	_, err := c.product.FindID(uint(req.Req.ProductId))
	if err != nil {
		return util.NotFoundCart(err.Error()), err
	}

	token, err_token := middleware.ReadToken()
	if err_token != nil {
		return util.UnauthorizedCart(err_token.Error()), err
	}

	input := entity.Cart{
		Quantity:  uint(req.Req.GetQuantity()),
		ProductID: uint(req.Req.GetProductId()),
		UserID:    token.ID,
	}

	err_create := c.repository.AddToCart(input)
	if err_create != nil {
		return util.BadRequestCart(err_create.Error()), err_create
	}

	return util.SuccessCart("success add to cart"), nil
}

func (c *CartServiceServer) DeleteCart(ctx context.Context, req *cart.ID) (*cart.Response, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return util.InternalServerErrorCart(err.Error()), err
	}

	err_delete := c.repository.Delete(uint(id))
	if err_delete != nil {
		return util.BadRequestCart(err_delete.Error()), err_delete
	}

	return util.DeletedCart("success delete cart by this id"), nil
}

func (c *CartServiceServer) ListCart(empty *emptypb.Empty, stream cart.CartService_ListCartServer) error {
	token, err := middleware.ReadToken()
	if err != nil {
		return constant.Unauthorized(err.Error())
	}

	data, err_list := c.repository.ListCart(token.ID)
	if err_list != nil {
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
		return util.FailureCreatedCart(err.Error()), err
	}

	product, err_find := c.product.FindID(data.ProductID)
	if err_find != nil {
		return util.NotFoundCartCreated(err_find.Error()), err_find
	}

	response := cart.Data{
		Uuid:        data.UUID,
		Quantity:    int64(data.Quantity),
		ProductName: product.Name,
		Price:       int64(product.Price),
		Total:       int64(product.Price * uint32(data.Quantity)),
	}

	return util.SuccessCreatedCart("success created for order", &response), nil
}

func (c *CartServiceServer) Update(ctx context.Context, req *cart.RequestUpdate) (*cart.Response, error) {
	request := entity.Cart{
		Quantity: uint(req.Req.GetQuantity()),
	}

	err := c.repository.Update(uint(req.Id), request)
	if err != nil {
		return util.BadRequestCart(err.Error()), err
	}

	return util.SuccessCart("success update cart"), nil
}