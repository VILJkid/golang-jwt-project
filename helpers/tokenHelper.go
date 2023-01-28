package helpers

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/VILJkid/golang-jwt-project/database"
	"github.com/VILJkid/golang-jwt-project/models"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	User_id    string
	User_type  string
	jwt.StandardClaims
}

var userModel *gorm.DB = database.ModelForDbOperations(database.DB, models.User{})
var secretKey string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, userType, userId string) (signedToken, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_id:    userId,
		User_type:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return
	}

	signedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))
	if err != nil {
		return
	}

	return
}

func UpdateAllTokens(signedToken, signedRefreshToken, userId string) (updatedUser models.User, err error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	err = userModel.WithContext(c).First(&updatedUser, &models.User{User_id: &userId}).Updates(&models.User{
		Token:         &signedToken,
		Refresh_token: &signedRefreshToken,
	}).First(&updatedUser, &models.User{User_id: &userId}).Error

	return
}

func ValidateToken(signedToken string) (claims *SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		err = errors.New("the token is invalid")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("the token is expired")
		return
	}

	return
}
