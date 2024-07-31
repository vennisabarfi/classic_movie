package handlers

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file.", err)
	}
	fmt.Println(".env file loaded successfully!")
}
