package service

import (
	"errors"
	"go-finance-wallet/internal/model"
	"go-finance-wallet/internal/repository"
	"go-finance-wallet/pkg/crypto"
	"os"

	"gorm.io/gorm"
)

type WalletService interface {
	GetBalance(userID uint) (*model.Wallet, error)
	TopUp(userID uint, amount int64) error
	Withdraw(userID uint, amount int64, pin string) error
}

type walletService struct {
	walletRepo repository.WalletRepository
	trxRepo    repository.TransactionRepository
	db         *gorm.DB
}

func NewWalletService(w repository.WalletRepository, t repository.TransactionRepository, db *gorm.DB) WalletService {
	return &walletService{w, t, db}
}

func (s *walletService) GetBalance(userID uint) (*model.Wallet, error) {
	wallet, err := s.walletRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("wallet tidak ditemukan")
	}

	secret := os.Getenv("SECRET_KEY")
	if !crypto.VerifySignature(wallet.ID, wallet.Balance, secret, wallet.Signature) {
		return nil, errors.New("data saldo tidak valid (manipulasi terdeteksi!)")
	}

	return wallet, nil
}

func (s *walletService) TopUp(userID uint, amount int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		wallet, _ := s.walletRepo.GetByUserID(userID)

		newBalance := wallet.Balance + amount
		secret := os.Getenv("SECRET_KEY")
		newSig := crypto.GenerateSignature(wallet.ID, newBalance, secret)

		if err := s.walletRepo.UpdateBalanceWithLock(tx, wallet.ID, newBalance, newSig); err != nil {
			return err
		}

		trx := &model.Transaction{
			WalletID: wallet.ID,
			Amount:   amount,
			Type:     "CREDIT",
		}
		return s.trxRepo.Create(tx, trx)
	})
}

func (s *walletService) Withdraw(userID uint, amount int64, pin string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		wallet, err := s.walletRepo.GetByUserID(userID)
		if err != nil {
			return err
		}

		if !crypto.CheckPasswordHash(pin, wallet.Pin) {
			return errors.New("PIN salah")
		}

		if wallet.Balance < amount {
			return errors.New("saldo tidak mencukupi")
		}

		newBalance := wallet.Balance - amount
		secret := os.Getenv("SECRET_KEY")
		newSig := crypto.GenerateSignature(wallet.ID, newBalance, secret)

		if err := s.walletRepo.UpdateBalanceWithLock(tx, wallet.ID, newBalance, newSig); err != nil {
			return err
		}

		trx := &model.Transaction{
			WalletID: wallet.ID,
			Amount:   amount,
			Type:     "DEBIT",
		}
		return s.trxRepo.Create(tx, trx)
	})
}
