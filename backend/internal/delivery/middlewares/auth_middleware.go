package middlewares

import (
	"echo-chat-app-backend/internal/models"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(mysqlDB *gorm.DB, authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header missing"})
			return
		}

		// 2. Cek format header apakah ada "Bearer"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		// 3. Ambil token dari header
		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		// 4. Verifikasi token menggunakan Firebase Auth
		token, err := authClient.VerifyIDToken(c, idToken)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		// 5. Cek apakah email ada di token
		if token.Claims["email"] == nil {
			c.JSON(401, gin.H{"error": "Email not found in token"})
			return
		}
		// cari user di database
		user := models.User{}
		err = mysqlDB.Where("email = ?", token.Claims["email"]).First(&user).Error
		if err != nil {
			c.JSON(401, gin.H{"error": "User not found"})
			return
		}

		// 6. Simpan informasi user ke context
		c.Set("id", user.ID)
		c.Set("email", token.Claims["email"])
		c.Set("firebase_uid", token.UID)
		c.Set("name", token.Claims["name"])
		c.Set("username", token.Claims["username"])
		c.Set("avatar_url", token.Claims["picture"])

		c.Next()
	}
}
