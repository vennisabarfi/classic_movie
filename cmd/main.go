package main

import (
	"classic_movies/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // postgres driver
	"github.com/lopezator/migrator"
)

var pool *sql.DB

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

	// Configure migrations
	m, err := migrator.New(
		migrator.Migrations(
			&migrator.Migration{
				Name: "Create table foo",
				Func: func(tx *sql.Tx) error {
					if _, err := tx.Exec("CREATE TABLE foo (id INT PRIMARY KEY)"); err != nil {
						return err
					}
					return nil
				},
			},
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// pool, err := sql.Open("postgres", "user=postgres password=PNeumono38%21 dbname=classic_movies_db sslmode=disable")
	if err != nil {
		log.Fatal("Error opening database connection", err)
	} else {
		fmt.Println("Database connection successful")
	}

	defer pool.Close()

	// ensure the database can be contacted
	if err = pool.Ping(); err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	} else {
		fmt.Println("Database successfully pinged")
	}

	// Migrate up
	if err := m.Migrate(pool); err != nil {
		log.Fatal(err)
	}

	// Security headers
	r.Use(securityHeaders)

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to Classic Movie Streaming Platform")
	})

	r.GET("/stream", handlers.StreamVideo)

	r.Run()

}
