package handlers

import (
	"blogo/database"
	"blogo/models"
	"context"
	"database/sql"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)






func Login(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error":"Invalid request"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
    defer cancel()

	var storedUser models.User
	err := database.DB.QueryRowContext(ctx ,"SELECT id, username, password FROM users WHERE username=$1", user.Username).
	Scan(&storedUser.ID, &storedUser.Username, &storedUser.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not create token"})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": tokenString})


}


func Register(c echo.Context) error {
    var user models.User
    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
    }

    if user.Username == "" || user.Password == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Username and password are required"})
    }
    if len(user.Password) < 8 {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Password must be at least 8 characters long"})
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not hash password"})
    }


    ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
    defer cancel()

    var existingUser string
    err = database.DB.QueryRowContext(ctx, "SELECT username FROM users WHERE username = $1", user.Username).Scan(&existingUser)
    if err != sql.ErrNoRows {
        return c.JSON(http.StatusConflict, echo.Map{"error": "Username already exists"})
    }


    _, err = database.DB.ExecContext(ctx, `
        INSERT INTO users (username, password, created_at) 
        VALUES ($1, $2, NOW())
    `, user.Username, hashedPassword)
    if err != nil {
        if strings.Contains(err.Error(), "unique_violation") { 
            return c.JSON(http.StatusConflict, echo.Map{"error": "Username already exists"})
        }
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
    }

    return c.JSON(http.StatusCreated, echo.Map{"message": "User created successfully"})
}