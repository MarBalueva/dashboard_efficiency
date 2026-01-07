package services

import (
	"net/http"
	"os"
	"time"

	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	return string(bytes), err
}

func CheckPassword(pass, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}

func GenerateJWT(userID uint, groups []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"groups":  groups,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func RequireGroup(allowedGroups ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id_not_found"})
			c.Abort()
			return
		}

		userID, ok := userIDVal.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_user_id"})
			c.Abort()
			return
		}

		var user models.User
		if err := db.DB.
			Preload("AccessGroups.AccessGroup").
			First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_not_found"})
			c.Abort()
			return
		}

		for _, ag := range user.AccessGroups {
			for _, group := range allowedGroups {
				if ag.AccessGroup.Code == group {
					c.Next()
					return
				}
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "Доступ запрещён",
		})
		c.Abort()
	}
}

func HasAnyGroup(user models.User, groups ...string) bool {
	for _, ug := range user.AccessGroups {
		for _, g := range groups {
			if ug.AccessGroup.Code == g {
				return true
			}
		}
	}
	return false
}
