package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/opentelemetry-go-contrib/launcher"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("godotenv.Load", err)
	}

	shutdown, err := launcher.ConfigureOpenTelemetry(
		honeycomb.WithApiKey(os.Getenv("HONEYKOMB_API_KEY")),
		launcher.WithServiceName(os.Getenv("OTEL_SERVICE_NAME")),
	)
	if err != nil {
		log.Fatal("launcher.ConfigureOpenTelemetry", err)
	}
	defer shutdown()

	server := gin.New()

	server.Any(
		"/my-endpoint",
		otelgin.Middleware("my-operation"),
		func(c *gin.Context) {
			fmt.Println("request received")
			c.String(http.StatusAccepted, "Hello World from Gin!\n")
		},
	)

	if err := http.ListenAndServe(":3030", server); err != nil {
		log.Fatal("http.ListenAndServe", err)
	}
}
