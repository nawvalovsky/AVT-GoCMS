package controllers

import (
	"gocms/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

// READ: Ambil semua post
func GetPosts(c *gin.Context) {
	var posts []models.Post
	models.DB.Find(&posts)
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// READ: Ambil satu post
func GetPost(c *gin.Context) {
	var post models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Postingan hilang ditelan bumi!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": post})
}

// CREATE: Buat post baru
func CreatePost(c *gin.Context) {
	var input models.Post
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Sanitasi HTML
	p := bluemonday.UGCPolicy()
	input.Content = p.Sanitize(input.Content)

	input.UserID = userId.(uint)
	if input.Slug == "" {
		input.Slug = strings.ReplaceAll(strings.ToLower(input.Title), " ", "-")
	}

	models.DB.Create(&input)
	c.JSON(http.StatusOK, gin.H{"data": input})
}

// UPDATE: Edit post
func UpdatePost(c *gin.Context) {
	var post models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post tidak ditemukan!"})
		return
	}

	// Cek apakah user pemilik post ini?
	userId, _ := c.Get("currentUser")
	if post.UserID != userId.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Jangan sentuh post orang lain, baka!"})
		return
	}

	var input models.Post
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update data yang diperbolehkan
	if input.Title != "" {
		post.Title = input.Title
		post.Slug = strings.ReplaceAll(strings.ToLower(input.Title), " ", "-")
	}
	if input.Content != "" {
		p := bluemonday.UGCPolicy()
		post.Content = p.Sanitize(input.Content)
	}
	if input.ImageURL != "" {
		post.ImageURL = input.ImageURL
	}

	models.DB.Save(&post)
	c.JSON(http.StatusOK, gin.H{"data": post})
}

// DELETE: Hapus post
func DeletePost(c *gin.Context) {
	var post models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post tidak ditemukan!"})
		return
	}

	// Cek kepemilikan
	userId, _ := c.Get("currentUser")
	if post.UserID != userId.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Dilarang menghapus kenangan orang lain!"})
		return
	}

	models.DB.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Postingan berhasil dimusnahkan!"})
}