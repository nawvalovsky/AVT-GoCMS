package controllers

import (
	"gocms/middleware"
	"gocms/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input salah!"})
		return
	}

	if err := models.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username tidak ditemukan!"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password salah!"})
		return
	}

	token, _ := middleware.GenerateToken(user.ID)

	// Set HttpOnly Cookie
	// Secure: false (localhost), true (production/https)
	c.SetCookie("Authorization", token, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login sukses!"})
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout sukses!"})
}

func CheckSession(c *gin.Context) {
	userId, _ := c.Get("currentUser")
	c.JSON(http.StatusOK, gin.H{"user_id": userId, "status": "authenticated"})
}