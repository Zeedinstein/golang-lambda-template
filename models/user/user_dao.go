package user

import (
	"context"
	"errors"
	"time"

	"github.com/GotBot-AI/users-api/clients"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func GetPagination(lastEvaluatedKey map[string]types.AttributeValue) (users []User, nextStartKey map[string]types.AttributeValue, err error) {
	logrus.Debug(TableName)
	paginator, err := clients.Dynamo.GetPaginator(TableName, 2, lastEvaluatedKey)
	if err != nil {
		logrus.Error(err)
		return nil, nil, err
	}

	if paginator.HasMorePages() {
		data, err := paginator.NextPage(context.TODO())
		if err != nil {
			logrus.Error(err)
			return nil, nil, err
		}

		err = attributevalue.UnmarshalListOfMaps(data.Items, &users)
		if err != nil {
			logrus.Error(err)
			return nil, nil, err
		}

		return users, data.LastEvaluatedKey, nil
	}
	logrus.Debug("Not more pages")
	return nil, nil, errors.New("no more pages")
}

func GetAll() (users []User, err error) {
	out, err := clients.Dynamo.GetItems(TableName)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &users)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return users, nil
}

func GetByID(userID string) (user User, err error) {
	out, err := clients.Dynamo.GetItem(TableName, map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: userID},
	})
	if err != nil {
		logrus.Error(err)
		return User{}, err
	}

	err = attributevalue.UnmarshalMap(out.Item, &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func Create(user *User) (*dynamodb.PutItemOutput, error) {
	// Create UNIQUE ID
	uid, err := uuid.NewUUID()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	user.ID = uid.String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	data, err := attributevalue.MarshalMap(user)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	out, err := clients.Dynamo.PutItem(TableName, data)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return out, nil
}

func Update(user *User) (*dynamodb.UpdateItemOutput, error) {
	user.UpdatedAt = time.Now()

	expr, err := expression.NewBuilder().WithUpdate(
		expression.Set(
			expression.Name("firstName"),
			expression.Value(user.FirstName),
		).Set(
			expression.Name("lastName"),
			expression.Value(user.LastName),
		).Set(
			expression.Name("email"),
			expression.Value(user.Email),
		).Set(
			expression.Name("updatedAt"),
			expression.Value(user.UpdatedAt),
		),
	).Build()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: user.ID},
	}

	out, err := clients.Dynamo.UpdateItem(TableName, key, expr)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return out, nil
}

func Delete(userID string) error {
	_, err := clients.Dynamo.DeleteItem(TableName, map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: userID},
	})
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func Validate(user User) error {
	if user.FirstName == "" {
		return errors.New("firstName required in body")
	}
	if user.LastName == "" {
		return errors.New("lastName required in body")
	}
	if user.Email == "" {
		return errors.New("email required in body")
	}
	return nil
}
