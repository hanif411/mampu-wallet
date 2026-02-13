package main

import (
	"go-finance-wallet/internal/handler"
	"go-finance-wallet/internal/middleware"
	"go-finance-wallet/internal/model"
	"go-finance-wallet/internal/repository"
	"go-finance-wallet/internal/service"
	"go-finance-wallet/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db := database.InitDB()
	db.AutoMigrate(&model.User{}, &model.Wallet{}, &model.Transaction{})
	log.Println("migration succes")

	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	trxRepo := repository.NewTransactionRepository(db)

	authService := service.NewAuthService(userRepo, walletRepo)
	walletService := service.NewWalletService(walletRepo, trxRepo, db)

	authHandler := handler.NewAuthHandler(authService)
	walletHandler := handler.NewWalletHandler(walletService)

	r := gin.Default()
	r.POST("/api/v1/register", authHandler.Register)
	r.POST("/api/v1/login", authHandler.Login)
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/balance", walletHandler.GetBalance)
		protected.POST("/topup", walletHandler.TopUp)
		protected.POST("/withdraw", walletHandler.Withdraw)
	}

	log.Println("Server running on :5000")
	r.Run(":5000")
}
