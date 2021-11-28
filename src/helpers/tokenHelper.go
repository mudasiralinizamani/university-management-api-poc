package helpers

import (
	"context"
	"fmt"
	"os"
	"time"
	"university-management-api/src/data"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	Uid       string
	jwt.StandardClaims
}

func GenerateTokens(email string, firstName string, lastName string, role string, uid string) (signedToken string, signedRefreshToken string, err error) {
	err = nil

	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(170)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		fmt.Println("error occurred while creating token")
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		fmt.Println("error occurred while creating refresh token")
		return
	}

	return token, refreshToken, err

}

func UpdateTokens(signedToken string, signedRefreshToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refreshtoken", signedRefreshToken})

	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updatedat", UpdatedAt})

	upsert := false
	filter := bson.M{"userid": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := data.UserCollection.UpdateOne(ctx, filter, bson.D{
		{"$set", updateObj},
	}, &opt)

	defer cancel()

	if err != nil {
		return
	}

}

// This is function will validate the Token, and it is used in the Authenticate Middleware - Mudasir Ali
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		// This msg will be shown in Output - Mudasir Ali
		msg = fmt.Sprintf("Token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		// This msg will be shown in Output - Mudasir Ali
		msg = fmt.Sprintf("Token is expired")
		msg = err.Error()
		return
	}

	return claims, msg
}
