package handlers

import (
	"blogo/database"
	"blogo/models"
	"database/sql"
	"net/http"
	"strconv"

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




//------------Comments



func CreateComment(c echo.Context) error {
	var comment models.Comment


	userID := c.Get("user_id").(int)


	articleID := c.Param("articleID")
	articleIDInt, err := strconv.Atoi(articleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid article ID"})
	}


	if err := c.Bind(&comment); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	_, err = database.DB.Exec(
		"INSERT INTO comments (content, user_id, article_id) VALUES ($1, $2, $3)",
		comment.Content, userID, articleIDInt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create comment"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "Comment added successfully"})
}



func GetCommentsByArticle(c echo.Context) error {
	articleID := c.Param("articleID")
	articleIDInt, err := strconv.Atoi(articleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid article ID"})
	}

	rows, err := database.DB.Query(
		"SELECT id, content, user_id, article_id FROM comments WHERE article_id=$1", articleIDInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch comments"})
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.ArticleID); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to read comments"})
		}
		comments = append(comments, comment)
	}

	return c.JSON(http.StatusOK, comments)
}


func GetCommentByID(c echo.Context) error {
	articleID := c.Param("articleID")
	commentID := c.Param("commentID")
	commentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid comment ID"})
	}

	articleIDInt, err := strconv.Atoi(articleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid article ID"})
	}

	var comment models.Comment
	// جستجو برای کامنت با هر دو articleID و commentID
	err = database.DB.QueryRow(
		"SELECT id, content, user_id, article_id FROM comments WHERE id=$1 AND article_id=$2",
		commentIDInt, articleIDInt).
		Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.ArticleID)

	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Comment not found"})
	}

	return c.JSON(http.StatusOK, comment)
}