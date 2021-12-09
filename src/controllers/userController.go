package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"university-management-api/src/data"
	"university-management-api/src/models"
	"university-management-api/src/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))

		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{
			{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}},
				{"total_count", bson.D{{"$sum", 1}}},
				{"data", bson.D{{"$push", "$$ROOT"}}}}}}

		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}

		result, err := data.UserCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "error occurred while listing users"})
		}

		var allUsers []bson.M

		if err = result.All(ctx, &allUsers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "error occurred while fetching result"})
		}

		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")

		err := services.CheckUser(user_id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "UserNotFound", "error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User

		err = data.UserCollection.FindOne(ctx, bson.M{"userid": user_id}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "Something went wrong in server"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func GetUsersByRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.Param("user_role")

		if role != "ADMIN" && role != "DEAN" && role != "HOD" && role != "COURSEADVISER" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "RoleNotExist", "error": "Role does not exist"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := data.UserCollection.Find(ctx, bson.M{"role": role})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "ServerError", "error": "error occurred while fetching result"})
			return
		}

		var users []models.User

		if err := result.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "error occurres while getting result"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}
