package services

import (
	"context"
	"errors"
	"time"
	"university-management-api/src/data"
	"university-management-api/src/models"

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

func GetUserById(userId string) (user models.User, err error) {
	err = nil
	var foundUser models.User

	if userExist := CheckUser(userId); userExist != nil {
		err = userExist
		return foundUser, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	foundErr := data.UserCollection.FindOne(ctx, bson.M{"userid": userId}).Decode(&foundUser)

	if foundErr != nil {
		err = foundErr
		return foundUser, err
	}

	return foundUser, err
}

func CheckUserByRole(user_id, user_role string) (err error) {
	err = nil
	var foundUser models.User

	user, err := GetUserById(user_id)

	if err != nil {
		return err
	}

	foundUser = user

	role := *&foundUser.Role

	if *role != user_role {
		err = errors.New("user is not a " + user_role)
		return err
	}
	return err
}
