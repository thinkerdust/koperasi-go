package middleware

import (
	"koperasi-go/helpers"
	"koperasi-go/repository"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helpers.Error(c, http.StatusUnauthorized, "Missing token")
			c.Abort()
			return
		}

		// Parse JWT
		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			helpers.Error(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		// Ambil claims user_id
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			helpers.Error(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			helpers.Error(c, http.StatusUnauthorized, "Invalid token payload")
			c.Abort()
			return
		}
		userID := uint(userIDFloat)

		// 🔹 Validasi token di DB
		valid, err := repository.CheckUserToken(userID, authHeader)
		if err != nil || !valid {
			helpers.Error(c, http.StatusUnauthorized, "Invalid token or session timeout")
			c.Abort()
			return
		}

		// Simpan user_id ke context
		c.Set("user_id", userID)

		c.Next()
	}
}
