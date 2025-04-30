package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Piyu-Pika/go-wallet/internal/database"
	"github.com/Piyu-Pika/go-wallet/internal/models"
	"github.com/Piyu-Pika/go-wallet/internal/routers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB initializes an in-memory SQLite database for testing
func SetupTestDB() {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	database.DB = db
	db.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{})
}

// SetupRouterForTest sets up the router for testing
func SetupRouterForTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return routers.SetupRouter()
}

// TestCreateUser tests the CreateUser handler
func TestCreateUser(t *testing.T) {
	SetupTestDB()
	router := SetupRouterForTest()

	userJSON := `{"name": "John Doe", "email": "john@example.com"}`
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(userJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["ID"])
	assert.Equal(t, "John Doe", response["Name"])
	assert.Equal(t, "john@example.com", response["Email"])
	assert.Len(t, response["Wallets"], 1)
}

// TestGetUserInfo tests the GetUserInfo handler
func TestGetUserInfo(t *testing.T) {
	SetupTestDB()
	router := SetupRouterForTest()

	// Create a user first
	user := models.User{
		Name:      "Jane Doe",
		Email:     "jane@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&user)

	req, _ := http.NewRequest("GET", "/api/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.User
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, user.Name, response.Name)
	assert.Equal(t, user.Email, response.Email)
}

// TestCreateWallet tests the CreateWallet handler
func TestCreateWallet(t *testing.T) {
	SetupTestDB()
	router := SetupRouterForTest()

	// Create a user first
	user := models.User{
		Name:      "Alice",
		Email:     "alice@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&user)

	walletJSON := `{"user_id": 1}`
	req, _ := http.NewRequest("POST", "/api/wallets", bytes.NewBufferString(walletJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Wallet
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, uint(1), response.UserID)
	assert.Equal(t, 0.0, response.Balance)
	assert.Equal(t, "USD", response.Currency)
}

// TestGetWalletInfo tests the GetWalletInfo handler
func TestGetWalletInfo(t *testing.T) {
	SetupTestDB()
	router := SetupRouterForTest()

	// Create a user and wallet
	user := models.User{
		Name:      "Bob",
		Email:     "bob@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&user)

	wallet := models.Wallet{
		UserID:    user.ID,
		Balance:   100.0,
		Currency:  "USD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&wallet)

	req, _ := http.NewRequest("GET", "/api/wallets/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Wallet
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, wallet.UserID, response.UserID)
	assert.Equal(t, wallet.Balance, response.Balance)
	assert.Equal(t, wallet.Currency, response.Currency)
}

// TestGetBalance tests the GetBalance handler
func TestGetBalance(t *testing.T) {
	SetupTestDB()
	router := SetupRouterForTest()

	// Create a user and wallet
	user := models.User{
		Name:      "Charlie",
		Email:     "charlie@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&user)

	wallet := models.Wallet{
		UserID:    user.ID,
		Balance:   200.0,
		Currency:  "USD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&wallet)

	req, _ := http.NewRequest("GET", "/api/wallets/1/balance", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]float64
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200.0, response["balance"])
}

// TestCreateTransaction tests the CreateTransaction handler
func TestCreateTransaction(t *testing.T) {
	SetupTestDB()
	router := SetupRouterForTest()

	// Create users and wallets
	user1 := models.User{
		Name:      "David",
		Email:     "david@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&user1)

	user2 := models.User{
		Name:      "Eva",
		Email:     "eva@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&user2)

	wallet1 := models.Wallet{
		UserID:    user1.ID,
		Balance:   500.0,
		Currency:  "USD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&wallet1)

	wallet2 := models.Wallet{
		UserID:    user2.ID,
		Balance:   300.0,
		Currency:  "USD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&wallet2)

	transactionJSON := `{"source_wallet_id": 1, "dest_wallet_id": 2, "amount": 100.0}`
	req, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBufferString(transactionJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Check updated balances
	var updatedWallet1, updatedWallet2 models.Wallet
	database.DB.First(&updatedWallet1, wallet1.ID)
	database.DB.First(&updatedWallet2, wallet2.ID)
	assert.Equal(t, 400.0, updatedWallet1.Balance)
	assert.Equal(t, 400.0, updatedWallet2.Balance)
}

// TestListTransactions tests the ListTransactions handler
func TestListTransactions(t *testing.T) {
	SetupTestDB()
	router := SetupRouterForTest()

	// Create users, wallets, and transactions
	user := models.User{
		Name:      "Frank",
		Email:     "frank@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&user)

	wallet := models.Wallet{
		UserID:    user.ID,
		Balance:   1000.0,
		Currency:  "USD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&wallet)

	transaction := models.Transaction{
		SourceWalletID: wallet.ID,
		DestWalletID:   wallet.ID, // Self-transfer for testing
		Amount:         50.0,
		Status:         "COMPLETED",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	database.DB.Create(&transaction)

	req, _ := http.NewRequest("GET", "/api/transactions/wallet/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Transaction
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response, 1)
	assert.Equal(t, transaction.Amount, response[0].Amount)
	assert.Equal(t, transaction.Status, response[0].Status)
}
