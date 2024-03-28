package route

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/xendit/xendit-go/v4"

	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/grpc/cart"
	"github.com/rulanugrh/order/internal/grpc/order"
	"github.com/rulanugrh/order/internal/middleware"
	"github.com/rulanugrh/order/internal/repository"
	"github.com/rulanugrh/order/internal/service"
	"github.com/rulanugrh/order/pkg"
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

func restServer(crt *service.CartServiceServer, ord *service.OrderServiceServer, conf *config.App) error {
	mux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := cart.RegisterCartServiceHandlerServer(ctx, mux, crt)
	if err != nil {
		return err
	}

	err = order.RegisterOrderServiceHandlerServer(ctx, mux, ord)
	if err != nil {
		return err
	}

	log.Printf("[*] Start HTTP Server at %s:%s", conf.Server.Host, conf.Server.HTTP)
	dsnHTTP := fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.HTTP)

	serv := http.Server{
		Addr: dsnHTTP,
		Handler: mux,
	}
	
	return serv.ListenAndServe()
}

func InitServer() {
	// Main config for Xendit and Config App
	conf := config.GetConfig()
	xenditAPI := xendit.NewClient(conf.Xendit.APIKey)
	xendit := pkg.XenditPluggin(xenditAPI, conf)

	// Config for Postgres
	postgres := config.InitializeDB(conf)
	postgres.StartConnection()
	postgres.Migrate()

	// Config For MongoDB
	mongo := config.InitializeMongo(conf)
	mongo.NewMongo()

	// Config for rabbitMQ and Run Consume
	rabbitmq := config.InitializeRabbit(conf)
	rabbitmq.InitRabbit()

	productRepository := repository.ProductRepository(mongo, conf)
	orderRepository := repository.OrderRepository(postgres)
	cartRepository := repository.CartRepository(postgres)
	
	// Confsume for Product Catch and Update
	rabbimq := pkg.RabbitMQ(rabbitmq, productRepository)
	err_catch := rabbimq.CatchProduct()
	if err_catch != nil {
		log.Println("[*] Error consume catch product: ", err_catch)
	}

	err_update := rabbimq.UpdateProduct()
	if err_update != nil {
		log.Println("[*] Error consume update product: ", err_update)

	}

	// Running Services and Listener GRPC and REST
	orderService := service.OrderService(orderRepository, productRepository, xendit, rabbimq)
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

	err_http := restServer(cartService, orderService, conf)
	if err_http != nil {
		log.Println("[*] Error Running HTTP: ", err_http)

	}
}