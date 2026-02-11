package repository

import (
	"go-finance-wallet/internal/model"

	"gorm.io/gorm"
)

type WalletRepository interface {
	GetByUserID(userID uint) (*model.Wallet, error)
	UpdateBalanceWithLock(tx *gorm.DB, walletID uint, newBalance int64, newSignature string) error
}

type walletRepo struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepo{db}
}

func (r *walletRepo) GetByUserID(userID uint) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	return &wallet, err
}

func (r *walletRepo) UpdateBalanceWithLock(tx *gorm.DB, walletID uint, newBalance int64, newSignature string) error {
	return tx.Model(&model.Wallet{}).Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?", walletID).
		Updates(map[string]interface{}{
			"balance":   newBalance,
			"signature": newSignature,
		}).Error
}
