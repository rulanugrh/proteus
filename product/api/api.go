package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rulanugrh/tokoku/product/api/handler"
	"github.com/rulanugrh/tokoku/product/internal/config"
	"github.com/rulanugrh/tokoku/product/internal/middleware"
	"github.com/rulanugrh/tokoku/product/internal/repository"
	"github.com/rulanugrh/tokoku/product/internal/service"
	"github.com/rulanugrh/tokoku/product/pkg"
)

type API struct {
	product  handler.ProductInterface
	comment  handler.CommentInterface
	category handler.CategoryInterface
}

func (a *API) ProductRoute(r *mux.Router) {
	app := r.PathPrefix("/api/product").Subrouter()
	app.HandleFunc("/create", a.product.Create).Methods("POST")
	app.HandleFunc("/get", a.product.FindAll).Methods("GET")
	app.HandleFunc("/find/{id}", a.product.FindID).Methods("GET")
	app.HandleFunc("/update/{id}", a.product.Update).Methods("PUT")
}

func (a *API) CommentRoute(r *mux.Router) {
	app := r.PathPrefix("/api/comment").Subrouter()
	app.HandleFunc("/create", a.comment.Create).Methods("POST")
	app.HandleFunc("/get", a.comment.FindUID).Methods("GET")
	app.HandleFunc("/product/{id}", a.comment.FindPID).Methods("GET")
}

func (a *API) CategoryRoute(r *mux.Router) {
	app := r.PathPrefix("/api/category").Subrouter()
	app.HandleFunc("/create", a.category.Create).Methods("POST")
	app.HandleFunc("/get", a.category.FindAll).Methods("GET")
	app.HandleFunc("/find/{id}", a.category.FindID).Methods("GET")
}

func main() {
	conf := config.GetConfig()
	db := config.InitializeDB(conf)
	db.ConnectionDB()
	db.Migration()

	rabbit := config.InitRabbit(conf)
	rabbit.InitRabbit()

	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())
	metric := pkg.NewPrometheus(reg)

	productRepo := repository.ProductRepository(db)
	commentRepo := repository.CommentRepository(db)
	categoryRepo := repository.CategoryRepository(db)

	productService := service.ProductService(productRepo)
	commentService := service.CommentService(commentRepo)
	categoryService := service.CategoryService(categoryRepo)

	rabbitInterface := pkg.RabbitMQ(*rabbit)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})

	api := API{
		product:  handler.ProductHandler(productService, rabbitInterface, metric),
		comment:  handler.CommentHandler(commentService, metric),
		category: handler.CategoryHandler(categoryService, metric),
	}

	route := mux.NewRouter()
	// middleware impelment
	route.Use(middleware.CORS)
	route.Use(middleware.ValidateToken)
	route.Handle("/metric", promHandler).Methods("GET")

	// routes app
	api.ProductRoute(route)
	api.CategoryRoute(route)
	api.CommentRoute(route)

	dsn := fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port)
	serve := http.Server{
		Addr:    dsn,
		Handler: route,
	}

	err := serve.ListenAndServe()
	if err != nil {
		log.Println("error, cant running http service")
	}

	log.Printf("running at %s", dsn)
}
