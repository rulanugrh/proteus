package route

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	// Logger initialize
	logger := pkg.Logrus()

	productRepository := repository.ProductRepository(mongo, conf)
	orderRepository := repository.OrderRepository(postgres)
	cartRepository := repository.CartRepository(postgres)
	
	// Confsume for Product Catch and Update
	rabbimq := pkg.RabbitMQ(rabbitmq, productRepository, orderRepository)
	err_catch := rabbimq.CatchProduct()
	if err_catch != nil {
		log.Println("[*] Error consume catch product: ", err_catch)
	}

	err_update := rabbimq.UpdateProduct()
	if err_update != nil {
		log.Println("[*] Error consume update product: ", err_update)

	}

	err_notified := rabbimq.NotifierPayment()
	if err_notified != nil {
		log.Println("[*] Error consume notifier webhook: ", err_notified)
	}


	// prometheus initalize
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector())
	metric := pkg.NewPrometheus(registry)
	metric.SetTotalCPU()
	metric.SetTotalMemory()

	// Running Services and Listener GRPC
	orderService := service.OrderService(orderRepository, productRepository, xendit, rabbimq, metric, logger)
	cartService := service.CartService(cartRepository, productRepository, metric, logger)

	dsnGRPC := fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.GRPC)
	listener, err := net.Listen("tcp", dsnGRPC)
	if err != nil {
		log.Println("[*] Error binding network", err)
	}

	err_grpc := grpcServer(cartService, orderService, conf, listener)
	if err_grpc != nil {
		log.Println("[*] Error Running GRPC: ", err_grpc)
	}

	// Running HTTP Service for Metric
	server := http.NewServeMux()
	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry})
	server.Handle("/metric", promHandler)
	
	dsnHTTP := fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.HTTP)

	serv := http.Server{
		Addr: dsnHTTP,
		Handler: server,
	}

	err_http := serv.ListenAndServe()
	if err_http != nil {
		log.Println("[*] Error Running HTTP: ", err_http)
	}

}