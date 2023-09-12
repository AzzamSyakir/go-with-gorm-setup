package routes

import (
	"golang-api/api/handlers"
	"golang-api/config"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	config.InitDB()
	// Create a Gin router
	router := gin.Default()

	// Define your routes and associate them with handlers
	router.POST("/api/user/Register", handlers.CreateUser)
	router.POST("/api/user/LoginUser", handlers.LoginUser)
	router.POST("/api/user/UpdateUser", handlers.UpdateUser)
	router.POST("/api/user/DeleteUser", handlers.DeleteUser)
	router.GET("/api/user/GetUser", handlers.GetUser)

	// Mulai server HTTP
	router.Run(":9000")
}
