package controllers

import (
	"context"
	"net/http"
	"time"
	"university-management-api/src/data"
	"university-management-api/src/models"
	"university-management-api/src/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateFaculty() gin.HandlerFunc {
	return func(c *gin.Context) {
		var faculty models.Faculty

		if err := c.BindJSON(&faculty); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": err.Error()})
			return
		}

		validationErr := validate.Struct(faculty)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "ValidationError", "error": validationErr.Error()})
			return
		}

		err := services.CheckFaculty(faculty.Name)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "FacultyExist", "error": "Faculty already exist"})
			return
		}

		err = services.CheckUser(faculty.DeanId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "UserNotFound", "error": "Dean does not exist"})
			return
		}

		err = services.CheckUserByRole(faculty.DeanId, "DEAN")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		deanUser, _ := services.GetUserById(faculty.DeanId)

		faculty.DeanName = string(*deanUser.FirstName + " " + *deanUser.LastName)
		faculty.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		faculty.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		faculty.ID = primitive.NewObjectID()
		faculty.FacultyId = faculty.ID.Hex()

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		facultyInsetionNumber, insertErr := data.FacultyCollection.InsertOne(ctx, faculty)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "Something went wrong server"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"succeeded": true, "insertedId": facultyInsetionNumber.InsertedID})
	}
}

func GetFaculties() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := data.FacultyCollection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "error occurred while fetching faculties"})
			return
		}

		var faculties []models.Faculty

		if err := result.All(ctx, &faculties); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "error occurred while fetching result"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"succeeded": true, "faculties": faculties})
	}
}

func GetFaculty() gin.HandlerFunc {
	return func(c *gin.Context) {
		facultyId := c.Param("faculty_id")

		err := services.CheckFacultyById(facultyId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "FacultyNotFound", "error": "Faculty does not exist"})
			return
		}

		var faculty models.Faculty
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = data.FacultyCollection.FindOne(ctx, bson.M{"facultyid": facultyId}).Decode(&faculty)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "Something went wrong in server"})
			return
		}
		c.JSON(http.StatusOK, faculty)
	}
}
