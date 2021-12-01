package services

import (
	"context"
	"errors"
	"time"
	"university-management-api/src/data"

	"go.mongodb.org/mongo-driver/bson"
)

func CheckUser(userId string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = nil

	count, _ := data.UserCollection.CountDocuments(ctx, bson.M{"userid": userId})
	defer cancel()

	if count == 0 {
		err = errors.New("user does not exist")
	}

	return err
}
