package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/VILJkid/golang-jwt-project/database"
	"github.com/VILJkid/golang-jwt-project/helpers"
	"github.com/VILJkid/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var userModel = database.ModelForDbOperations(database.DB, models.User{})
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := validate.Struct(user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var count int64
		if err := userModel.WithContext(c).Where(&models.User{Email: user.Email}).Count(&count).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Email already exists",
			})
			return
		}

		if err := userModel.WithContext(c).Where(&models.User{Phone: user.Phone}).Count(&count).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Phone already exists",
			})
			return
		}
	}
}

func Login()

func GetUsers()

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")

		if err := helpers.MatchUserTypeToUid(ctx, userId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var user models.User
		err := userModel.WithContext(c).First(&user, &models.User{User_id: &userId}).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}
