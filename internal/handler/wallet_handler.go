package handler

import (
	"go-finance-wallet/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletService service.WalletService
}

func NewWalletHandler(s service.WalletService) *WalletHandler {
	return &WalletHandler{s}
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID := c.MustGet("user_id").(uint) 

	wallet, err := h.walletService.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": wallet.Balance})
}

func (h *WalletHandler) TopUp(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req struct {
		Amount int64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "jumlah harus lebih dari 0"})
		return
	}

	if err := h.walletService.TopUp(userID, req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "top up berhasil"})
}

func (h *WalletHandler) Withdraw(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req struct {
		Amount int64  `json:"amount" binding:"required,gt=0"`
		Pin    string `json:"pin" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.walletService.Withdraw(userID, req.Amount, req.Pin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "withdraw berhasil"})
}
