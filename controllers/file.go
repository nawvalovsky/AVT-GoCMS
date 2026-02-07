package controllers

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Mana gambarnya?!"})
		return
	}

	// Validasi ekstensi file
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		c.JSON(400, gin.H{"error": "Format file tidak didukung! Pakai jpg/png/webp saja."})
		return
	}

	// Ganti nama file biar unik
	filename := fmt.Sprintf("%d%s", time.Now().Unix(), ext)
	path := "uploads/" + filename

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(500, gin.H{"error": "Gagal simpan gambar. Server lelah."})
		return
	}

	// Buat URL lengkap secara dinamis
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	
	fullURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, c.Request.Host, filename)
	
	c.JSON(200, gin.H{"url": fullURL})
}