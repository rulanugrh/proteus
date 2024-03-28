package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/grpc/cart"
	"github.com/rulanugrh/order/internal/grpc/order"
	"github.com/rulanugrh/order/internal/middleware"
	"github.com/rulanugrh/order/internal/repository"
	"github.com/rulanugrh/order/internal/service"
	"github.com/rulanugrh/order/pkg"
	"github.com/xendit/xendit-go/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func grpcServer(crt *service.CartServiceServer, ord *service.OrderServiceServer, conf *config.App, nets net.Listener) error {
	interceptors := middleware.NewGRPCInterceptors(conf)
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptors.Unary()),
		grpc.StreamInterceptor(interceptors.Stream()),
	}

	serve := grpc.NewServer(serverOptions...)
	cart.RegisterCartServiceServer(serve, crt)
	order.RegisterOrderServiceServer(serve, ord)

	reflection.Register(serve)

	log.Printf("[*] Start GRPC Server at %s", nets.Addr().String())
	return serve.Serve(nets)
}

func main() {
	conf := config.GetConfig()

	xenditAPI := xendit.NewClient(conf.Xendit.APIKey)
	postgres := config.InitializeDB(conf)
	postgres.StartConnection()
	postgres.Migrate()

	mongo := config.InitializeMongo(conf)
	mongo.NewMongo()

	xendit := pkg.XenditPluggin(xenditAPI, conf)

	rabbitmq := config.InitializeRabbit(conf)
	rabbitmq.InitRabbit()

	productRepository := repository.ProductRepository(mongo, conf)
	orderRepository := repository.OrderRepository(postgres)
	cartRepository := repository.CartRepository(postgres)

	orderService := service.OrderService(orderRepository, productRepository, xendit)
	cartService := service.CartService(cartRepository, productRepository)

	dsnGRPC := fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.GRPC)
	listener, err := net.Listen("tcp", dsnGRPC)
	if err != nil {
		log.Println("[*] Error binding network", err)
	}

	err_grpc := grpcServer(cartService, orderService, conf, listener)
	if err_grpc != nil {
		log.Println("[*] Error Running GRPC: ", err_grpc)
	}
}
