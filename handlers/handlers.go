package handlers

import (
	"blogo/database"
	"blogo/models"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)



func CreateArticle(c echo.Context) error {
	var article models.Article
	if err := c.Bind(&article); err != nil {
		return c.String(http.StatusBadRequest, "Invalid input")
	}

	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE id=$1)", article.AuthorID).Scan(&exists)
	if err != nil || !exists {
		return c.String(http.StatusNotFound, "User not found")
	}


	query := "INSERT INTO articles (title, content, author_id, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id"
	err = database.DB.QueryRow(query, article.Title, article.Content, article.AuthorID).Scan(&article.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create article")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Article created successfully",
	})
}

//Get all articles
func GetArticles(c echo.Context) error {
	rows, err := database.DB.Query("SELECT id, title, content, author_id, created_at, updated_at FROM articles")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retrieve articles")
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.AuthorID, &article.CreatedAt, &article.UpdatedAt); err != nil {
			return c.String(http.StatusInternalServerError, "Error scanning article")
		}
		articles = append(articles, article)
	}

	return c.JSON(http.StatusOK, articles)
}

func GetArticle(c echo.Context) error {
	id := c.Param("id")
	var article models.Article

	query := "SELECT id, title, content, author_id, created_at, updated_at FROM articles WHERE id=$1"
	err := database.DB.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Content, &article.AuthorID, &article.CreatedAt, &article.UpdatedAt)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, "Article not found")
	} else if err != nil {
		return c.String(http.StatusInternalServerError, "Error retrieving article")
	}

	return c.JSON(http.StatusOK, article)
}

func DeleteArticle(c echo.Context) error {
	id := c.Param("id")
	query := "DELETE FROM articles WHERE id=$1"

	_, err := database.DB.Exec(query, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete article")
	}

	return c.JSON(http.StatusOK, "Article deleted successfully")
}