package api

import (
	"koperasi-go/helpers"
	"koperasi-go/model"
	"koperasi-go/repository"
	"koperasi-go/db"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var input struct {
		NIK string `json:"nik"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	user, err := repository.FindUserByNIK(input.NIK)
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

	// simpan token ke DB
	if err := repository.UpdateUserToken(user.ID, tokenString); err != nil {
		helpers.Error(c, http.StatusInternalServerError, "Failed to save token")
		return
	}

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

func ChangePassword(c *gin.Context) {
	var input struct {
		OldPassword     string `json:"old_password"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// ambil user_id dari middleware Auth
	userIDVal, exists := c.Get("user_id")
	if !exists {
		helpers.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDVal.(uint)

	// cek user
	var user model.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		helpers.Error(c, http.StatusUnauthorized, "User not found")
		return
	}

	// cek old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid old password ")
		return
	}

	// validasi password baru
	if len(input.Password) < 4 {
		helpers.Error(c, http.StatusBadRequest, "Minimal length password is 4")
		return
	}
	if input.Password != input.ConfirmPassword {
		helpers.Error(c, http.StatusBadRequest, "Invalid confirm password")
		return
	}

	// hash password baru
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// update password
	if err := db.DB.Model(&user).Updates(map[string]interface{}{
		"password":   string(hashedPassword),
		"updated_at": time.Now(),
	}).Error; err != nil {
		helpers.Error(c, http.StatusInternalServerError, "Update password failed")
		return
	}

	helpers.Success(c, "Update password success", nil)
}

func Logout(c *gin.Context) {
	// Ambil user_id dari context (diset di middleware JWT)
	userId, exists := c.Get("user_id")
	if !exists {
		helpers.Error(c, http.StatusForbidden, "Access denied")
		return
	}

	// Hapus token di DB
	if err := repository.ClearUserToken(userId.(uint)); err != nil {
		helpers.Error(c, http.StatusInternalServerError, "Logout failed")
		return
	}

	helpers.Success(c, "Logout success", gin.H{})
}
