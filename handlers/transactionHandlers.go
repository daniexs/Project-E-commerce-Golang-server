package handlers

import (
	"main/helper"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTransactionHistoriesForUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDParam, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User detail not found"})
			return
		}

		userID, ok := userIDParam.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			return
		}

		var transactionHistories []models.TransactionHistory
		if err := db.Joins("Product").Where("user_id = ?", userID).Find(&transactionHistories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction histories"})
			return
		}
		transformedTransactionHistories := make([]map[string]interface{}, len(transactionHistories))
		for i, t := range transactionHistories {
			transformedTransaction := map[string]interface{}{
				"id":          t.ID,
				"product_id":  t.ProductID,
				"user_id":     t.UserID,
				"quantity":    t.Quantity,
				"total_price": t.TotalPrice,
				"Product": map[string]interface{}{
					"id":          t.Product.ID,
					"title":       t.Product.Title,
					"price":       t.Product.Price,
					"stock":       t.Product.Stock,
					"category_id": t.Product.CategoryID,
					"created_at":  t.Product.CreatedAt,
					"updated_at":  t.Product.UpdatedAt,
				},
			}
			transformedTransactionHistories[i] = transformedTransaction
		}
		c.JSON(http.StatusOK, transformedTransactionHistories)
	}
}

func GetAllTransactionHistories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transactionHistories []models.TransactionHistory
		if err := db.Joins("Product").Joins("User").Find(&transactionHistories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction histories"})
			return
		}

		transformedTransactionHistories := make([]map[string]interface{}, len(transactionHistories))
		for i, t := range transactionHistories {
			transformedTransaction := map[string]interface{}{
				"id":          t.ID,
				"product_id":  t.ProductID,
				"user_id":     t.UserID,
				"quantity":    t.Quantity,
				"total_price": t.TotalPrice,
				"Product": map[string]interface{}{
					"id":          t.Product.ID,
					"title":       t.Product.Title,
					"price":       t.Product.Price,
					"stock":       t.Product.Stock,
					"category_id": t.Product.CategoryID,
					"created_at":  t.Product.CreatedAt,
					"updated_at":  t.Product.UpdatedAt,
				},
				"User": map[string]interface{}{
					"id":         t.User.ID,
					"email":      t.User.Email,
					"full_name":  t.User.FullName,
					"balance":    t.User.Balance,
					"created_at": t.User.CreatedAt,
					"updated_at": t.User.UpdatedAt,
				},
			}
			transformedTransactionHistories[i] = transformedTransaction
		}

		c.JSON(http.StatusOK, gin.H{"transaction_histories": transformedTransactionHistories})
	}
}

type CreateTransactionInput struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required"`
}

func CreateTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDParam, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User detail not found"})
			return
		}

		userID, ok := userIDParam.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			return
		}
		var input CreateTransactionInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := helper.Validate(input); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		var product models.Product
		if err := db.First(&product, input.ProductID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		// Check stock availability
		if input.Quantity > product.Stock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		totalPrice := input.Quantity * product.Price

		// Check user balance
		if totalPrice > user.Balance {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
			return
		}

		// Deduct stock from product and balance from user
		product.Stock -= input.Quantity
		user.Balance -= totalPrice

		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
			return
		}

		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user balance"})
			return
		}

		// Increment sold_product_amount in category
		var category models.Category
		if err := tx.First(&category, product.CategoryID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find category"})
			return
		}
		category.SoldProductAmount += input.Quantity
		if err := tx.Save(&category).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
			return
		}

		// Create transaction history
		transaction := models.TransactionHistory{
			UserID:     userID,
			ProductID:  input.ProductID,
			Quantity:   input.Quantity,
			TotalPrice: totalPrice,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction history"})
			return
		}

		tx.Commit()

		c.JSON(http.StatusCreated, gin.H{
			"message": "You have successfully purchased the product",
			"transaction_bill": gin.H{
				"total_price":   transaction.TotalPrice,
				"quantity":      transaction.Quantity,
				"product_title": product.Title,
			},
		})
	}
}
