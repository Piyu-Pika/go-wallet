package routers

import (
	handlers "github.com/Piyu-Pika/go-wallet/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// User routes
	router.POST("/api/users", handlers.CreateUser)
	router.GET("/api/users/:id", handlers.GetUserInfo)

	// Wallet routes
	router.POST("/api/wallets", handlers.CreateWallet)
	router.GET("/api/wallets/:id", handlers.GetWalletInfo)
	router.POST("/api/wallets/:id/deposit", handlers.Deposit)
	router.GET("/api/wallets/:id/balance", handlers.GetBalance)

	// Transaction routes
	router.POST("/api/transactions", handlers.CreateTransaction)
	router.GET("/api/transactions/wallet/:wallet_id", handlers.ListTransactions)

	return router
}
