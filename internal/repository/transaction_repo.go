package repository

import (
	"go-finance-wallet/internal/model"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *gorm.DB, transaction *model.Transaction) error
	GetByWalletID(walletID uint) ([]model.Transaction, error)
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepo{db}
}

func (r *transactionRepo) Create(tx *gorm.DB, transaction *model.Transaction) error {
	return tx.Create(transaction).Error
}

func (r *transactionRepo) GetByWalletID(walletID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Where("wallet_id = ?", walletID).Order("created_at desc").Find(&transactions).Error
	return transactions, err
}
