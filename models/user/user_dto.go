package user

import (
	"os"
	"time"
)

var (
	TableName = os.Getenv("USERS_TABLE")
)

type User struct {
	ID        string    `dynamodbav:"id,omitempty" json:"id"`
	FirstName string    `dynamodbav:"firstName,omitempty" json:"firstName"`
	LastName  string    `dynamodbav:"lastName,omitempty" json:"lastName"`
	Email     string    `dynamodbav:"email,omitempty" json:"email"`
	CreatedAt time.Time `dynamodbav:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time `dynamodbav:"updatedAt,omitempty" json:"updatedAt"`
}
