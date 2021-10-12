package main

import (
	"context"
	"os"

	"github.com/GotBot-AI/users-api/api/routes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadaptor "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var adapter *fiberadaptor.FiberLambda

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetReportCaller(true)
	// Only log the warning severity or above.
	if os.Getenv("DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Info("Debug enabled")
		// Log as JSON instead of the default ASCII formatter.
		logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: false,
		})
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.Info("Debug disabled")
		// Log as JSON instead of the default ASCII formatter.
		logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: false,
		})
	}

	app := fiber.New()
	router := app.Group("/api/resources")
	routes.CreateRoutes(router)
	adapter = fiberadaptor.New(app)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return adapter.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
