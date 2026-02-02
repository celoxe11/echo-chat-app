package middlewares

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthMiddleware(authClient *auth.Client) gin.HandlerFunc {
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

		// 5. Simpan informasi user ke context
		c.Set("email", token.Claims["email"])
		c.Set("firebase_uid", token.UID)
		// tambahkan informasi lain sesuai kebutuhan

		c.Next()
	}
}
