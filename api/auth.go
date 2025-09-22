package api

import (
	"koperasi-go/helpers"
	"koperasi-go/model"
	"koperasi-go/repository"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	user, err := repository.FindUserByUsername(input.Username)
	if err != nil {
		helpers.Error(c, http.StatusUnauthorized, "User not found")
		return
	}

	if !helpers.CheckPasswordHash(input.Password, user.Password) {
		helpers.Error(c, http.StatusUnauthorized, "Invalid password")
		return
	}

	// generate JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	helpers.Success(c, "Login success", gin.H{"token": tokenString})
}

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	hash, _ := helpers.HashPassword(input.Password)
	user := model.User{
		Username: input.Username,
		Password: hash,
		Email:    input.Email,
	}

	if err := repository.CreateUser(&user); err != nil {
		helpers.Error(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	helpers.Success(c, "User registered", user)
}
