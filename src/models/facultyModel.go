package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Faculty struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `json:"name" validate:"required"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	DeanId    string             `json:"deanId" validate:"required"`
	DeanName  string             `json:"deanName"`
	FacultyId string             `json:"facultyId"`
}
