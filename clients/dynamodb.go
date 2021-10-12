package clients

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
)

var Dynamo = &dynamoDBClient{}

type DatabaseKeys struct {
	ID        string `dynamodbav:"id"`
	Namespace string `dynamodbav:"namespace"`
}

type dynamoDBClient struct{}

func (dynamoDBClient) GetClient() (*dynamodb.Client, error) {
	var cfg aws.Config
	var err error

	local := os.Getenv("AWS_SAM_LOCAL")
	logrus.Debug(local)

	cfg, err = config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = os.Getenv("REGION")
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if local != "true" {
		svc := dynamodb.NewFromConfig(cfg)
		return svc, nil
	} else {
		svc := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.EndpointResolver = dynamodb.EndpointResolverFromURL("http://dynamodb.local:8000")
		})
		return svc, nil
	}
}

func (d dynamoDBClient) GetPaginator(tableName string, limit int32, lastEvaluatedKey map[string]types.AttributeValue) (*dynamodb.ScanPaginator, error) {
	svc, err := d.GetClient()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	paginator := dynamodb.NewScanPaginator(svc, &dynamodb.ScanInput{
		TableName:         aws.String(tableName),
		Limit:             aws.Int32(limit),
		ExclusiveStartKey: lastEvaluatedKey,
	})
	return paginator, nil
}

func (d dynamoDBClient) GetItems(tableName string) (*dynamodb.ScanOutput, error) {
	svc, err := d.GetClient()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	out, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return out, nil
}

func (d dynamoDBClient) GetItem(tableName string, key map[string]types.AttributeValue) (*dynamodb.GetItemOutput, error) {
	svc, err := d.GetClient()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	out, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}

	return out, err
}

func (d dynamoDBClient) PutItem(tableName string, item map[string]types.AttributeValue) (*dynamodb.PutItemOutput, error) {
	svc, err := d.GetClient()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	out, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return out, nil
}

func (d dynamoDBClient) UpdateItem(tableName string, keys map[string]types.AttributeValue, expr expression.Expression) (*dynamodb.UpdateItemOutput, error) {
	svc, err := d.GetClient()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if expr.Update() == nil {
		logrus.Error("No update expression provided")
		return nil, errors.New("update expression required")
	}

	out, err := svc.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       keys,
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return out, nil
}

func (d dynamoDBClient) DeleteItem(tableName string, key map[string]types.AttributeValue) (*dynamodb.DeleteItemOutput, error) {
	svc, err := d.GetClient()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	out, err := svc.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
