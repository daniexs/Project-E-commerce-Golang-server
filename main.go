package main

import (
	"log"
	"main/config"
	"main/handlers"
	"main/middleware"
	"main/models"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbConfig := config.GetDbConfig()
	dsn := dbConfig.GetDBURL()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.TransactionHistory{})
	if err != nil {
		log.Fatal("Failed to auto migrate", err)
	}

	r := gin.Default()

	r.POST("/users/register", handlers.CreateUser(db))
	r.POST("/users/login", handlers.UserLogin(db))
	r.Use(middleware.TokenAuthMiddleware(db))
	r.PATCH("/users/topup", handlers.UpdateBalance(db))
	r.POST("/categories", middleware.AdminAuthMiddleware(), handlers.CreateCategory(db))
	r.GET("/categories", middleware.AdminAuthMiddleware(), handlers.GetCategories(db))
	r.PATCH("/categories/:categoryId", middleware.AdminAuthMiddleware(), handlers.UpdateCategory(db))
	r.DELETE("/categories/:categoryId", middleware.AdminAuthMiddleware(), handlers.DeleteCategory(db))
	r.POST("/products", middleware.AdminAuthMiddleware(), handlers.CreateProduct(db))
	r.GET("/products", handlers.GetAllProducts(db))
	r.PUT("/products/:productId", middleware.AdminAuthMiddleware(), handlers.UpdateProduct(db))
	r.DELETE("/products/:productId", middleware.AdminAuthMiddleware(), handlers.DeleteProduct(db))
	r.POST("/transactions", handlers.CreateTransaction(db))
	r.GET("/transactions/my-transactions", handlers.GetTransactionHistoriesForUser(db))
	r.GET("/transactions/user-transactions", middleware.AdminAuthMiddleware(), handlers.GetAllTransactionHistories(db))

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
