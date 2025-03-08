package main

import (
	"blogo/database"
	"blogo/handlers"
	"blogo/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	database.ConnectDB()
	e := echo.New()


	
	
	api := e.Group("/api/v1")
	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)
	//Articles
	api.POST("/articles", handlers.CreateArticle)
	api.GET("/articles", handlers.GetArticles)
	api.GET("/article/:id", handlers.GetArticle)
	api.DELETE("/article/:id", handlers.DeleteArticle, middleware.JWTMiddleware)



	e.Start(":8080")
}