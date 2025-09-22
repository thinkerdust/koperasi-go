package main

import (
	"fmt"
	"koperasi-go/db"
	"koperasi-go/model"
	"koperasi-go/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()
	db.DB.AutoMigrate(&model.User{}) // migrate user table

	r := gin.Default()
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(fmt.Sprintf(":%s", port))
}
