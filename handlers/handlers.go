package handlers

import (
	"blogo/database"
	"blogo/models"
	"net/http"

	"github.com/labstack/echo/v4"
)



func RegisterUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "Invalid input")
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Failed to register user" )
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "User registered successfully",
	})
}

func CreateArticle(c echo.Context) error {
	var article models.Article
	if err := c.Bind(&article); err != nil {
		return c.String(http.StatusBadRequest, "Invalid input")
	}
	var user models.User
    if err := database.DB.First(&user, article.AuthorID).Error; err != nil {
        return c.String(http.StatusNotFound, "User not found")
    }
	if err := database.DB.Create(&article).Error; err != nil {
        return c.String(http.StatusInternalServerError, "Failed to create article")
    }
	return c.JSON(http.StatusOK, map[string]string{
        "message": "Article created successfully",
    })

}

//Get all articles
func GetArticles(c echo.Context) error {
	var articles []models.Article
	database.DB.Find(&articles)
	return c.JSON(http.StatusOK, articles)
}

func GetArticle(c echo.Context) error {
	id := c.Param("id")
	var article models.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Article not found....")
	}
	return c.JSON(http.StatusOK, article)
}

func DeleteArticle(c echo.Context) error {
	id := c.Param("id")
	database.DB.Delete(&models.Article{}, id)
	return c.JSON(http.StatusOK, "article deleted successfully...")
}