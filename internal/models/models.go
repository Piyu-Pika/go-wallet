package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Wallets   []Wallet `gorm:"foreignKey:UserID"`
}

type Wallet struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"index;not null"`
	Balance   float64 `gorm:"not null;default:0.0"`
	Currency  string  `gorm:"not null;default:'USD'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transaction struct {
	ID             uint    `gorm:"primaryKey"`
	SourceWalletID uint    `gorm:"index;not null"`
	DestWalletID   uint    `gorm:"index;not null"`
	Amount         float64 `gorm:"not null"`
	Status         string  `gorm:"not null;default:'PENDING'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
