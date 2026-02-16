package controllers

import (
	"net/http"
	"os"
	"time"

	"clinic-backend/internal/config"
	"clinic-backend/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hash)

	config.DB.Create(&user)
	c.JSON(200, gin.H{"message": "User registered"})
}

func Login(c *gin.Context) {
	var body models.User
	var user models.User

	c.BindJSON(&body)
	config.DB.Where("email = ?", body.Email).First(&user)

	if user.ID == 0 {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Wrong password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
