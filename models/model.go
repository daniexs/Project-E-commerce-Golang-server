package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	FullName  string    `json:"full_name" validate:"required"`
	Email     string    `validate:"required,email" json:"email" `
	Password  string    `json:"password" validate:"required,min=6"`
	Role      string    `json:"role" validate:"required,oneof=admin customer"`
	Balance   int       `json:"balance" validate:"required,min=0,max=100000000"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	gorm.Model
	ID         uint      `gorm:"primaryKey" json:"id" validate:"-"`
	Title      string    `json:"title" validate:"required"`
	Price      int       `json:"price" validate:"required,min=0,max=50000000"`
	Stock      int       `json:"stock" validate:"required,min=5"`
	CategoryID uint      `json:"category_id" validate:"-"`
	Category   Category  `gorm:"foreignKey:CategoryID;references:ID"`
	CreatedAt  time.Time `json:"created_at" validate:"-"`
	UpdatedAt  time.Time `json:"updated_at" validate:"-"`
}

type Category struct {
	gorm.Model
	ID                uint      `gorm:"primaryKey" json:"id"`
	Type              string    `json:"type" validate:"required"`
	SoldProductAmount int       `json:"sold_product_amount" validate:"-"`
	CreatedAt         time.Time `json:"created_at" validate:"-"`
	UpdatedAt         time.Time `json:"updated_at" validate:"-"`
	Products          []Product `gorm:"foreignKey:CategoryID"`
}

type TransactionHistory struct {
	gorm.Model
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProductID  uint      `json:"product_id" validate:"required"`
	UserID     uint      `json:"user_id"`
	Quantity   int       `json:"quantity" validate:"required"`
	TotalPrice int       `json:"total_price" validate:"required"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Product    Product   `gorm:"foreignKey:ProductID;references:ID"`
	User       User      `gorm:"foreignKey:UserID;references:ID"`
}
