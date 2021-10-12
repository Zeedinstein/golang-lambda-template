package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

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
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	if len(sqsEvent.Records) == 0 {
		return errors.New("no SQS message passed to function")
	}

	for _, msg := range sqsEvent.Records {
		logrus.Infof("Got SQS message %q with body %q\n", msg.MessageId, msg.Body)
	}

	return nil
}
