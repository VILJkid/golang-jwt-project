package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/VILJkid/golang-jwt-project/database"
	helper "github.com/VILJkid/golang-jwt-project/helpers"
	"github.com/VILJkid/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var userModel = database.ModelForDbOperations(database.DB, models.User{})
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func Signup()

func Login()

func GetUsers()

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")

		if err := helper.MatchUserTypeToUid(ctx, userId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c, cancel := context.WithTimeout(context.Background(), time.Second*100)

		var user models.User
		err := userModel.WithContext(c).First(&user, "user_id = ?", userId).Error
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, user)
	}
}
