package services

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("role")
	err = nil

	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("role")
	uid := c.GetString("uid")
	err = nil

	if userType != "ADMIN" && uid != userId {
		err = errors.New("unauthorized to acccess this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}

func CheckPassword(hashedPassword, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	var check bool = true
	var msg string = ""

	if err != nil {
		msg = "Password is incorrect"
		check = false
	}
	return check, msg
}
