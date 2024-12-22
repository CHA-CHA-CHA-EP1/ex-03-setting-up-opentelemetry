package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var (
	meter       = otel.GetMeterProvider().Meter("ping_service")
	pingCounter metric.Int64Counter
)

func initProvider() {
	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
	)
	otel.SetMeterProvider(provider)
}

func main() {
	initProvider()

	var err error
	pingCounter, err = meter.Int64Counter("ping_total",
		metric.WithDescription("Total number of ping requests"),
	)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/ping", handlePing)
	r.Run(":8080")
}

func handlePing(c *gin.Context) {
	ctx := c.Request.Context()
	pingCounter.Add(ctx, 1)
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
