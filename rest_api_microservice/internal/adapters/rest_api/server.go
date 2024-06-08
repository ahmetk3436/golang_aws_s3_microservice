package rest_api

import (
	_ "bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/docs"
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/internal/ports"
	"fmt"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
)

type Adapter struct {
	api    ports.Api
	port   int
	server fiber.App
}

func NewAdapter(api ports.Api, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a *Adapter) Start() error {
	a.server = *fiber.New()
	a.server.Use(otelfiber.Middleware())
	a.server.Get("/product", a.GetProductById)

	a.server.Get("/swagger/*", swagger.HandlerDefault)

	log.Printf("serving metrics at localhost:%v/metrics", a.port)
	a.server.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	return a.server.Listen(fmt.Sprintf(":%d", a.port))
}

func (a *Adapter) Stop() error {
	return a.server.Shutdown()
}
