package main

import (
	"log"
	"net/http"
	"time"

	"ihandify-go/api" // Replace with your actual Go module path

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type GenerateKeyRequest struct {
	Engines          []string `json:"engines"`
	ExpiresInSeconds *int     `json:"expiresInSeconds"`
}

// Simple CORS Middleware mimicking FastAPI's behavior
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Api-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Initialize upstream service configuration
	api.InitConfig()

	r := gin.Default()
	r.Use(CORSMiddleware())

	// Note: FastAPI had two "/" routes. The first one handles the demo file return.
	r.GET("/", func(c *gin.Context) {
		c.File("static/demo.html")
	})

	// Serve static files
	r.Static("/static", "./static")

	// API Route
	r.POST("/api/generate-scoped-public-key", func(c *gin.Context) {
		var req GenerateKeyRequest

		// binding:"required" ensures fields are present
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid request payload structure",
			})
			return
		}

		result, err := api.GenerateScopedPublicKey(req.Engines, req.ExpiresInSeconds)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"detail": "Upstream service error",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	// Server configurations
	srv := &http.Server{
		Addr:         "0.0.0.0:3000",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 35 * time.Second,
	}

	log.Println("Starting server on port 3000...")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
