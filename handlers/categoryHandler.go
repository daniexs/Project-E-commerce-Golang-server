package handlers

import (
	"main/helper"
	"main/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateCategoryInput struct {
	Type string `json:"type" validate:"required"`
}

func CreateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateCategoryInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := helper.Validate(input); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		// Create a new category
		newCategory := models.Category{
			Type:              input.Type,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			SoldProductAmount: 0,
		}

		// Save the new category to the database
		if err := db.Create(&newCategory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"id":                  newCategory.ID,
			"type":                newCategory.Type,
			"sold_product_amount": newCategory.SoldProductAmount,
			"created_at":          newCategory.CreatedAt,
		})
	}
}

func GetCategories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []models.Category

		// Retrieve all categories from the database
		if err := db.Preload("Products").Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}
		transformedCategories := make([]map[string]interface{}, len(categories))
		for i, t := range categories {
			products := make([]map[string]interface{}, len(t.Products))
			for j, p := range t.Products {
				product := map[string]interface{}{
					"id":          p.ID,
					"title":       p.Title,
					"price":       p.Price,
					"stock":       p.Stock,
					"category_id": p.CategoryID,
					"created_at":  p.CreatedAt,
					"updated_at":  p.UpdatedAt,
				}
				products[j] = product
			}

			transformedCategory := map[string]interface{}{
				"id":                  t.ID,
				"type":                t.Type,
				"sold_product_amount": t.SoldProductAmount,
				"created_at":          t.CreatedAt,
				"updated_at":          t.UpdatedAt,
				"Products":            products,
			}
			transformedCategories[i] = transformedCategory
		}

		c.JSON(http.StatusOK, transformedCategories)
	}
}

type UpdateCategoryInput struct {
	Type string `json:"type" validate:"required"`
}

func UpdateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("categoryId")
		id, err := strconv.ParseUint(categoryID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		var category models.Category
		if err := db.First(&category, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		var input UpdateCategoryInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := helper.Validate(input); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		// Update category details
		category.Type = input.Type
		category.UpdatedAt = time.Now()

		if err := db.Save(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":                  category.ID,
			"type":                category.Type,
			"sold_product_amount": category.SoldProductAmount,
			"updated_at":          category.UpdatedAt,
		})
	}
}

func DeleteCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("categoryId")
		id, err := strconv.ParseUint(categoryID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		var category models.Category
		if err := db.First(&category, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		// Delete the category
		if err := db.Delete(&category, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category has been successfully deleted"})
	}
}
