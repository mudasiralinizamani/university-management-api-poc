package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    *string            `json:"firstName" validate:"required,min=2,max=12"`
	LastName     *string            `json:"lastName" validate:"required,min=2,max=12"`
	Email        *string            `json:"email" validate:"required,email"`
	Password     *string            `json:"password" validate:"required,min=6"`
	Role         *string            `json:"role" validate:"required,eq=ADMIN|eq=STUDENT|eq=HOD|eq=DEAN|eq=COURSEADVISER"`
	ProfilePic   *string            `json:"profilePic" validate:"required"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshToken"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	UserId       string             `json:"userId"`
}
