package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Department struct {
	ID                   primitive.ObjectID `bson:"_id"`
	Name                 string             `json:"name" validate:"required"`
	CreatedAt            time.Time          `json:"createdAt"`
	UpdatedAt            time.Time          `json:"updatedAt"`
	DepartmentId         string             `json:"departmentId"`
	HeadOfDepartmentId   string             `json:"headOfDepartmentId"`
	HeadOfDepartmentName string             `json:"headOfDepartmentName"`
	CourseAdviserId      string             `json:"courseAdviserId"`
	CourseAdviserName    string             `json:"courseAdviserName"`
	FacultyId            string             `json:"facultyId"`
	FacultyName          string             `json:"facultyName"`
}
