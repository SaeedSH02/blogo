package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)


var SecretKey = []byte(os.Getenv("JWT_SECRET"))

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Missing token"})
		}

		tokenString := strings.TrimPrefix(authHeader, "")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token claims"})
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid user_id"})
		}

		userID := int(userIDFloat)

		c.Set("user_id", userID)

		return next(c)
	}
}