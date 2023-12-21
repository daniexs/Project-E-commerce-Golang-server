package handlers

import (
	"fmt"
	"main/helper"
	"main/middleware"

	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser models.User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		existingUser := models.User{}
		if err := db.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
		newUser.Balance = 0
		newUser.Role = "customer"
		hashed, err := helper.HashPassword(newUser.Password)
		if err != nil {
			fmt.Println("Error hashing password:", err)
			return
		}
		newUser.Password = hashed
		if err := db.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"id":         newUser.ID,
			"full_name":  newUser.FullName,
			"email":      newUser.Email,
			"password":   newUser.Password,
			"balance":    0,
			"created_at": newUser.CreatedAt,
		})
	}
}

func UserLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var existingUser models.User
		err := db.Where("email = ?", user.Email).First(&existingUser).Error
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Email or password is wrong"})
			return
		}

		err = helper.VerifyPassword(existingUser.Password, user.Password)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Email or password is wrong"})
			return
		}

		fmt.Print(user)
		token, err := middleware.CreateToken(user.Email, existingUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token", "error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func UpdateBalance(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmail, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User detail not found"})
			return
		}

		var user models.User
		if err := db.Where("email = ?", userEmail).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Get the new balance value from the request body
		var updateData struct {
			Balance int `json:"balance"`
		}
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the user's balance
		user.Balance = updateData.Balance + user.Balance
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Your balance has been successfully updated to Rp %d", user.Balance)})
	}
}
