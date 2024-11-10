package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"upbot-server-go/database"
	"upbot-server-go/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserInfo struct {
	Email string `json:"email"`
}

// Define a struct for JWT claims
type CustomClaims struct {
	UserID uint   `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func Google(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header missing"})
		c.Abort()
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization format"})
		c.Abort()
		return
	}
	accessToken := parts[1]

	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=" + accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user info"})
		c.Abort()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid access token"})
		c.Abort()
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read response body"})
		c.Abort()
		return
	}
	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse response body"})
		c.Abort()
		return
	}

	var user models.User
	if err := database.DB.First(&user, "email = ?", userInfo.Email).Error; err != nil {
		user = models.User{Email: userInfo.Email}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
			return
		}
	}

	signedToken, err := GetSignedToken(userInfo.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to sign token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Auth successfull",
		"token":   signedToken,
	})
}

func GetSignedToken(email string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	claims := CustomClaims{
		Email: email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
