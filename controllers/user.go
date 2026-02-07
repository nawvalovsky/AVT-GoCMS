package controllers

import (
	"gocms/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CREATE: Tambah Admin Baru
func CreateUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid!"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{Username: input.Username, Password: string(hashedPassword)}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username sudah ada atau database error!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pasukan bertambah! Admin baru siap bertugas."})
}

// READ: Daftar Semua Admin
func GetUsers(c *gin.Context) {
	var users []models.User
	// Ambil ID dan Username saja, buang password!
	if err := models.DB.Select("id", "username", "created_at").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melihat daftar pasukan!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// DELETE: Pecat Admin
func DeleteUser(c *gin.Context) {
	currentUserId, _ := c.Get("currentUser")
	
	var user models.User
	if err := models.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan!"})
		return
	}

	// Jangan biarkan Master menghapus diri sendiri!
	if user.ID == currentUserId.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Master tidak boleh memecat diri sendiri! Siapa nanti yang pimpin?"})
		return 
	}

	models.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Admin berhasil ditendang dari markas!"})
}