package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title    string `json:"title" binding:"required"`
	Slug     string `json:"slug" gorm:"uniqueIndex"`
	Content  string `json:"content" gorm:"type:text"` // Isi blog (bisa HTML/Markdown)
	ImageURL string `json:"image_url"`                // Path gambar
	UserID   uint   `json:"user_id"`                  // Relasi ke User
}