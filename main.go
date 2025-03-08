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
	api.POST("/articles", handlers.CreateArticle, middleware.JWTMiddleware)
	api.GET("/articles", handlers.GetArticles)
	api.GET("/article/:id", handlers.GetArticle)
	api.DELETE("/article/:id", handlers.DeleteArticle, middleware.JWTMiddleware)

	comment :=  e.Group("/api/v1/articles")
	api.POST("/articles/:articleID/comments", handlers.CreateComment, middleware.JWTMiddleware)
	comment.GET("/:articleID/comments", handlers.GetCommentsByArticle, middleware.JWTMiddleware)
	comment.GET("/:articleID/comments/:commentID", handlers.GetCommentsByArticle, middleware.JWTMiddleware)


	e.Start(":8080")
}