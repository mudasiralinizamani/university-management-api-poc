package controllers

import (
	"context"
	"net/http"
	"time"
	dtos "university-management-api/src/Dtos"
	"university-management-api/src/data"
	"university-management-api/src/models"
	"university-management-api/src/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "ValidationError", "error": validationErr.Error()})
			defer cancel()
			return
		}

		count, err := data.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "error occurred while counting user email"})
			defer cancel()
			return
		} else if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": "EmailAlreadyExist", "error": "email address already exists"})
			defer cancel()
			return
		}

		hashedPassword := services.HashPassword(*user.Password)

		user.Password = &hashedPassword

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()

		token, refreshToken, err := services.GenerateTokens(*user.Email, *user.FirstName, *user.LastName, *user.Role, user.UserId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "error occurred while generating token"})
			defer cancel()
			return
		}

		user.Token = &token
		user.RefreshToken = &refreshToken

		resultInsertionNumber, insertErr := data.UserCollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "error occurred while creating user"})
			defer cancel()
			return
		}

		c.JSON(http.StatusOK, gin.H{"succeeded": true, "insertionNumber": resultInsertionNumber})
	}
}

func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var dto dtos.Signin
		var foundUser models.User

		if err := c.BindJSON(&dto); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(dto)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "ValidationError", "error": validationErr.Error()})
			defer cancel()
			return
		}

		err := data.UserCollection.FindOne(ctx, bson.M{"email": dto.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "EmailNotFound", "error": "email address does not exist"})
			defer cancel()
			return
		}

		isPasswordValid, msg := services.CheckPassword(*foundUser.Password, *dto.Password)

		if !isPasswordValid {
			c.JSON(http.StatusBadRequest, gin.H{"code": "IncorrectPassword", "error": msg})
			defer cancel()
			return
		}

		token, refreshToken, err := services.GenerateTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.Role, foundUser.UserId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": " Error occurred while creating token and refresh token"})
			return
		}

		services.UpdateTokens(token, refreshToken, foundUser.UserId)

		err = data.UserCollection.FindOne(ctx, bson.M{"userid": foundUser.UserId}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}
