package helpers

import (
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
	Uid        string
	User_type  string
	jwt.StandardClaims
}

var userModel *gorm.DB = database.ModelForDbOperations(database.DB, models.User{})
var secretKey string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, userType, uid string) (signedToken, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
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
