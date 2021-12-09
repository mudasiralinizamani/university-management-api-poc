package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DepartmentsStudents struct {
	ID                    primitive.ObjectID `bson:"_id"`
	DepartmentsStudentsId string             `json:"departmentsStudenetId"`
	DepartmentName        string             `json:"departmentName"`
	DepartmentId          string             `json:"DepartmentId"`
	StudentName           string             `json:"studentName"`
	StudentId             string             `json:"studentId"`
}
