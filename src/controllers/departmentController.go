package controllers

import (
	"net/http"
	"university-management-api/src/models"
	"university-management-api/src/services"

	"github.com/gin-gonic/gin"
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

		err := services.DoesDepartmentExist(department.DepartmentId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "DepartmentExists", "error": "Department already exists"})
			return
		}

		err = services.CheckUserByRole(department.HeadOfDepartmentId, "HOD")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "UserError", "error": err.Error()})
			return
		}

		err = services.CheckUserByRole(department.CourseAdviserId, "COURSEADVISER")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "UserError", "error": err.Error()})
			return
		}

		err = services.CheckFacultyById(department.FacultyId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "FacultyNoFound", "error": "Faculty does not exists"})
			return
		}

	}
}
