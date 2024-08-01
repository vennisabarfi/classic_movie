package main

import (
	"classic_movies/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func securityHeaders(c *gin.Context) {

	port := os.Getenv("PORT")
	expectedHost := "localhost:" + port

	if c.Request.Host != expectedHost {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
		return
	}
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Header("Referrer-Policy", "strict-origin")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
	c.Next()
}

func main() {
	// load .env
	handlers.LoadEnv()

	r := gin.Default()

	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error opening database connection")
	} else {
		fmt.Println("Database connection successful")
	}

	defer pool.Close()

	// Security headers
	r.Use(securityHeaders)

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to Classic Movie Streaming Platform")
	})

	r.GET("/stream", handlers.StreamVideo)

	r.Run()

}
