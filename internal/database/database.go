package database

import (
	"errors"
	"log"
	"os"

	"github.com/Piyu-Pika/go-wallet/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file, will use environment variables directly")
	}

	dsn := os.Getenv("COCKROACH_DSN")
	if dsn == "" {
		log.Fatal("Environment variable COCKROACH_DSN is not set")
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to CockroachDB:", err)
	}
	log.Println("Successfully connected to database")

	err = DB.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
	log.Println("Database schema migration successful")
}

func GetDB() *gorm.DB {
	return DB
}

func CreateUser(user *models.User) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		wallet := &models.Wallet{
			UserID:   user.ID,
			Balance:  0,
			Currency: "USD",
		}
		return tx.Create(wallet).Error
	})
}

func AddBalance(walletID uint, amount float64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.First(&wallet, walletID).Error; err != nil {
			return err
		}

		wallet.Balance += amount
		return tx.Save(&wallet).Error
	})
}

func GetWalletBalance(walletID uint) (float64, error) {
	var wallet models.Wallet
	if err := DB.First(&wallet, walletID).Error; err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func TransferFunds(fromWalletID, toWalletID uint, amount float64) error {
	if fromWalletID == toWalletID {
		return errors.New("cannot transfer to the same wallet")
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		var fromWallet models.Wallet
		if err := tx.First(&fromWallet, fromWalletID).Error; err != nil {
			return err
		}

		if fromWallet.Balance < amount {
			return errors.New("insufficient funds")
		}

		var toWallet models.Wallet
		if err := tx.First(&toWallet, toWalletID).Error; err != nil {
			return err
		}

		if err := tx.Model(&fromWallet).Update("balance", fromWallet.Balance-amount).Error; err != nil {
			return err
		}

		if err := tx.Model(&toWallet).Update("balance", toWallet.Balance+amount).Error; err != nil {
			return err
		}

		transaction := &models.Transaction{
			SourceWalletID: fromWalletID,
			DestWalletID:   toWalletID,
			Amount:         amount,
			Status:         "COMPLETED",
		}
		return tx.Create(transaction).Error
	})
}

func GetWalletTransactions(walletID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := DB.Where("source_wallet_id = ? OR dest_wallet_id = ?", walletID, walletID).Find(&transactions).Error
	return transactions, err
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Error getting underlying DB instance:", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Println("Error closing database connection:", err)
	} else {
		log.Println("Database connection closed successfully")
	}
}
