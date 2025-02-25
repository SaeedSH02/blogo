package main

import (
	"blogo/database"
	"blogo/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	database.ConnectDB()
	e := echo.New()


	
	
	api := e.Group("/api/v1")
	api.POST("/register", handlers.RegisterUser)
	//Articles
	api.POST("/articles", handlers.CreateArticle)
	api.GET("/articles", handlers.GetArticles)
	api.GET("/article/:id", handlers.GetArticle)
	api.DELETE("/article/:id", handlers.DeleteArticle)



	e.Start(":8080")
}