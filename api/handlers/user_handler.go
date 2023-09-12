package handlers

import (
	"errors"
	"golang-api/api/models"
	"golang-api/config"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Memeriksa apakah ada username dan password yang diberikan
	if user.Username == "" || user.Password == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email dan password harus diisi"})
		return
	}
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(HashedPassword)

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}
func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Memeriksa apakah ada username dan password yang diberikan
	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password harus diisi"})
		return
	}

	// Cari pengguna berdasarkan username dan password di database
	var userDb models.User
	if err := config.DB.Where("username = ?", user.Username).First(&userDb).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username tidak terdaftar"})
		return
	}

	// Mengambil kata sandi dari permintaan pengguna
	passwordFromUser := user.Password

	// Verifikasi kata sandi menggunakan bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(passwordFromUser)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password salah"})
		return
	}

	// Generate bearer token dengan JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := jwt.MapClaims{
		"username": userDb.Username,
	}
	token.Claims = claims

	// Sign the token dengan kunci rahasia yang aman
	tokenString, err := token.SignedString([]byte("secret_key")) // Gantilah dengan kunci rahasia yang aman
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Kirim respons dengan token dan pesan berhasil login
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil login",
		"token":   tokenString,
	})
}

func GetUser(c *gin.Context) {
	var user models.User
	userID := c.Param("id")
	result := config.DB.First(&user, userID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	var user models.User
	userID := c.Param("id")
	result := config.DB.First(&user, userID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&user)
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	userID := c.Param("id")
	result := config.DB.Delete(&user, userID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
