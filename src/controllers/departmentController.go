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

func CreateDepartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var department models.Department

		if err := c.BindJSON(&department); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "error occurred while binding json data"})
			return
		}

		validationErr := validate.Struct(department)

		if validationErr != nil {
			c.JSON(http.StatusOK, gin.H{"code": "ValidationError", "error": validationErr.Error()})
			return
		}

		err := services.DoesDepartmentExist(department.Name)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "DepartmentExists", "error": "Department already exists"})
			return
		}

		hod, err := services.GetUserById(department.HeadOfDepartmentId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "UserNotFound", "error": "HOD does not exist"})
			return
		}

		courseAdviser, err := services.GetUserById(department.CourseAdviserId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "UserNotFound", "error": "CourseAdviser does not exist"})
			return
		}

		err = services.CheckUserByRole(department.HeadOfDepartmentId, "HOD")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "UserNotHod", "error": "This user is not a HOD"})
			return
		}

		err = services.CheckUserByRole(department.CourseAdviserId, "COURSEADVISER")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "UserNotCourseAdviser", "error": "This user is not a CourseAdvier"})
			return
		}

		err = services.CheckFacultyById(department.FacultyId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "FacultyNoFound", "error": "Faculty does not exists"})
			return
		}

		faculty, _ := services.GetFacultyById(department.FacultyId)

		department.HeadOfDepartmentId = hod.UserId
		department.HeadOfDepartmentName = string(*hod.FirstName + " " + *hod.LastName)
		department.CourseAdviserId = courseAdviser.UserId
		department.CourseAdviserName = string(*courseAdviser.FirstName + " " + *courseAdviser.LastName)
		department.FacultyId = faculty.FacultyId
		department.FacultyName = faculty.Name

		department.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		department.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		department.ID = primitive.NewObjectID()
		department.DepartmentId = department.ID.Hex()

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		departmentInsertNumber, insertErr := data.DepartmentCollection.InsertOne(ctx, department)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "error occurred while inserting department"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"succeeded": true, "insertedId": departmentInsertNumber.InsertedID})
	}
}

func GetDepartments() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := data.DepartmentCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "Error occurred while fetching result"})
			return
		}

		var departments []models.Department

		if err := result.All(ctx, &departments); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "Error occurred while fetching result"})
			return
		}

		c.JSON(http.StatusOK, departments)
	}
}

func GetDepartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		department_id := c.Param("department_id")

		err := services.CheckDepartmentById(department_id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "DepartmentNotFound", "error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var department models.Department

		err = data.DepartmentCollection.FindOne(ctx, bson.M{"departmentid": department_id}).Decode(&department)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "ServerError", "error": "Error occurred while fetching Department"})
			return
		}

		c.JSON(http.StatusOK, department)
	}
}

func GetFacultyDepartments() gin.HandlerFunc {
	return func(c *gin.Context) {
		faculty_id := c.Param("faculty_id")

		err := services.CheckFacultyById(faculty_id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "FacultyNotFound", "error": "Faculty Does not exist"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := data.DepartmentCollection.Find(ctx, bson.M{"facultyid": faculty_id})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "Error occurred while getting result"})
			return
		}

		var departments []models.Department

		if err = result.All(ctx, &departments); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "ServerError", "error": "Error occurred while fetching result"})
			return
		}

		c.JSON(http.StatusOK, departments)
	}
}
