package main

import (
	"gocms/controllers"
	"gocms/middleware"
	"gocms/models"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 5)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			return
		}
		c.Next()
	}
}

func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default env vars")
	}

	r := gin.Default()
	models.ConnectDatabase()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4321", "https://falhafizh.vercel.app"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	r.Use(SecureHeaders())
	r.Use(RateLimitMiddleware())

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		// PUBLIC ROUTES
		api.POST("/login", controllers.Login)
		api.POST("/logout", controllers.Logout)
		api.GET("/posts", controllers.GetPosts)
		api.GET("/posts/:id", controllers.GetPost)
		api.GET("/me", middleware.AuthMiddleware(), controllers.CheckSession)

		// PROTECTED ROUTES (Butuh Login)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// CRUD POSTS
			protected.POST("/posts", controllers.CreatePost)
			protected.PUT("/posts/:id", controllers.UpdatePost)    // NEW: Edit
			protected.DELETE("/posts/:id", controllers.DeletePost) // NEW: Hapus

			protected.POST("/upload", controllers.UploadImage)

			// MANAJEMEN USER (Admin Only)
			protected.POST("/users", controllers.CreateUser)       // NEW: Buat Akun
			protected.DELETE("/users/:id", controllers.DeleteUser) // NEW: Hapus Akun
			protected.GET("/users", controllers.GetUsers) // Tambahkan baris ini!
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}