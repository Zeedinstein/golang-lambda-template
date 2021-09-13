package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	fiberadaptor "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type Channel struct {
	ChannelID string `dynamodbav:"channelID" json:"channelID,omitempty"`
	Thing     string `dynamodbav:"thing" json:"thing,omitempty"`
}

var adapter *fiberadaptor.FiberLambda

func init() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("yes")
	})
	app.Get("/ping", func(c *fiber.Ctx) error {
		log.Println("Handler!!")
		return c.SendString("pong")
	})

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	app.Get("/tables", func(c *fiber.Ctx) error {
		// Using the SDK's default configuration, loading additional config
		// and credentials values from the environment variables, shared
		// credentials, and shared configuration files
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("REGION")))
		if err != nil {
			return err
		}

		var channels []Channel

		// Using the Config value, create the DynamoDB client
		svc := dynamodb.NewFromConfig(cfg)

		channel := Channel{
			ChannelID: "HASDASDASD",
			Thing:     "thisthing",
		}

		channelMarshal, err := attributevalue.MarshalMap(channel)
		if err != nil {
			panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
		}

		_, err = svc.PutItem(context.Background(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("CHANNELS_TABLE")),
			Item:      channelMarshal,
		})
		if err != nil {
			return err
		}

		// Build the request with its input parameters
		resp, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName: aws.String(os.Getenv("CHANNELS_TABLE")),
			Limit:     aws.Int32(5),
		})
		if err != nil {
			return err
		}

		err = attributevalue.UnmarshalListOfMaps(resp.Items, &channels)
		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", err))
		}
		return c.JSON(channels)
	})

	adapter = fiberadaptor.New(app)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return adapter.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
