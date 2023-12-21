package handlers

import (
	"main/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateProductInput struct {
	Title      string `json:"title" binding:"required"`
	Price      int    `json:"price" binding:"required"`
	Stock      int    `json:"stock" binding:"required"`
	CategoryID uint   `json:"category_id" binding:"required"`
}

func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateProductInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the provided category ID exists in the database
		var category models.Category
		if err := db.First(&category, input.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID not found"})
			return
		}

		// Create a new product
		newProduct := models.Product{
			Title:      input.Title,
			Price:      input.Price,
			Stock:      input.Stock,
			CategoryID: input.CategoryID,
		}

		// Save the new product to the database
		if err := db.Create(&newProduct).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":          newProduct.ID,
			"title":       newProduct.Title,
			"stock":       newProduct.Stock,
			"price":       newProduct.Price,
			"category_Id": newProduct.CategoryID,
			"created_at":  newProduct.CreatedAt,
		})
	}
}

func GetAllProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []models.Product

		// Retrieve all products from the database
		if err := db.Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
			return
		}

		transformedProducts := make([]map[string]interface{}, len(products))
		for i, p := range products {
			transformedProduct := map[string]interface{}{
				"id":          p.ID,
				"title":       p.Title,
				"stock":       p.Stock,
				"price":       p.Price,
				"category_Id": p.CategoryID,
				"created_at":  p.CreatedAt,
			}
			transformedProducts[i] = transformedProduct
		}

		c.JSON(http.StatusOK, transformedProducts)
	}
}

type UpdateProductInput struct {
	Title      string `json:"title" binding:"required"`
	Price      int    `json:"price" binding:"required"`
	Stock      int    `json:"stock" binding:"required"`
	CategoryID uint   `json:"category_id" binding:"required"`
}

func UpdateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("productId")
		id, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var product models.Product
		if err := db.First(&product, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		var input UpdateProductInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update product details
		product.Title = input.Title
		product.Price = input.Price
		product.Stock = input.Stock
		product.CategoryID = input.CategoryID
		product.UpdatedAt = time.Now()

		if err := db.Save(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
			return
		}

		dataProduct := map[string]interface{}{
			"id":         product.ID,
			"title":      product.Title,
			"stock":      product.Stock,
			"price":      product.Price,
			"CategoryId": product.CategoryID,
			"createdAt":  product.CreatedAt,
			"updatedAt":  product.UpdatedAt,
		}

		c.JSON(http.StatusOK, gin.H{"product": dataProduct})
	}
}

func DeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("productId")
		id, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var product models.Product
		if err := db.First(&product, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		if err := db.Delete(&product, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product haas been successfully deleted"})
	}
}
