package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/VILJkid/golang-jwt-project/database"
	"github.com/VILJkid/golang-jwt-project/helpers"
	"github.com/VILJkid/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var userModel = database.ModelForDbOperations(database.DB, models.User{})
var validate = validator.New()

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword), err
}

func VerifyPassword(userPassword string, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
}

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

		*user.User_id = uuid.New().String()
		token, refreshToken, err := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *user.User_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		user.Token = &token
		user.Refresh_token = &refreshToken

		password, err := HashPassword(*user.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		user.Password = &password

		if err := userModel.WithContext(c).Create(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var user, foundUser models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := userModel.WithContext(c).First(&foundUser, &models.User{Email: user.Email}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := VerifyPassword(*user.Password, *foundUser.Password); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if foundUser.Email == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "user not found",
			})
			return
		}

		token, refreshToken, err := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *foundUser.User_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if foundUser, err = helpers.UpdateAllTokens(token, refreshToken, *foundUser.User_id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := helpers.CheckUserType(ctx, "ADMIN"); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		page, _ := strconv.Atoi(ctx.Query("page"))
		limit, _ := strconv.Atoi(ctx.Query("limit"))
		sort_by_column := ctx.Query("sort_by_column")
		sort_direction := ctx.Query("sort_direction")

		p := &database.Pagination{
			Limit:         limit,
			Page:          page,
			SortByColumn:  sort_by_column,
			SortDirection: sort_direction,
		}

		var users []models.User
		if err := userModel.WithContext(c).Scopes(database.Paginate(&models.User{}, p, userModel)).Find(&users).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		p.Rows = users
		ctx.JSON(http.StatusOK, users)
	}
}

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
		if err := userModel.WithContext(c).First(&user, &models.User{User_id: &userId}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}
