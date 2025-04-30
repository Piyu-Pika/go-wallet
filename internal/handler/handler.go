package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Piyu-Pika/go-wallet/internal/database"
	"github.com/Piyu-Pika/go-wallet/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var userInput struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var existingUser models.User
	if err := db.Where("email = ?", userInput.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	newUser := models.User{
		Name:      userInput.Name,
		Email:     userInput.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Wallets: []models.Wallet{
			{
				Balance:   0,
				Currency:  "USD",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func GetUserInfo(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db := database.GetDB()
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
func Deposit(c *gin.Context) {
	var addBalanceInput struct {
		WalletID uint    `json:"wallet_id" binding:"required"`
		Amount   float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&addBalanceInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// db := database.GetDB()
	if err := database.AddBalance(addBalanceInput.WalletID, addBalanceInput.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add balance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Balance added successfully"})
}

func CreateWallet(c *gin.Context) {
	var walletInput struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&walletInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var user models.User
	if err := db.First(&user, walletInput.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	newWallet := models.Wallet{
		UserID:    walletInput.UserID,
		Balance:   0,
		Currency:  "USD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&newWallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}
	c.JSON(http.StatusCreated, newWallet)
}

func GetWalletInfo(c *gin.Context) {
	walletID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	db := database.GetDB()
	var wallet models.Wallet
	if err := db.First(&wallet, walletID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

func GetBalance(c *gin.Context) {
	walletID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	balance, err := database.GetWalletBalance(uint(walletID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func CreateTransaction(c *gin.Context) {
	var txRequest struct {
		SourceWalletID uint    `json:"source_wallet_id" binding:"required"`
		DestWalletID   uint    `json:"dest_wallet_id" binding:"required"`
		Amount         float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&txRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.TransferFunds(txRequest.SourceWalletID, txRequest.DestWalletID, txRequest.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Transaction processed"})
}

func ListTransactions(c *gin.Context) {
	walletID, err := strconv.Atoi(c.Param("wallet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	transactions, err := database.GetWalletTransactions(uint(walletID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, transactions)
}
