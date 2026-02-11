package service

import (
	"errors"
	"go-finance-wallet/internal/model"
	"go-finance-wallet/internal/repository"
	"go-finance-wallet/pkg/crypto"
	"os"
)

type AuthService interface {
	Register(username, password, pin string) error
	Login(username, password string) (string, error)
}

type authService struct {
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
}

func NewAuthService(u repository.UserRepository, w repository.WalletRepository) AuthService {
	return &authService{u, w}
}

func (s *authService) Register(username, password, pin string) error {
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}

	hashedPin, err := crypto.HashPassword(pin)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: username,
		Password: hashedPassword,
	}

	secret := os.Getenv("SECRET_KEY")
	initialSignature := crypto.GenerateSignature(1, 0, secret)

	return s.userRepo.CreateWithWallet(user, hashedPin, initialSignature)
}

func (s *authService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", errors.New("user tidak ditemukan")
	}

	if !crypto.CheckPasswordHash(password, user.Password) {
		return "", errors.New("password salah")
	}

	return crypto.GenerateJWT(user.ID)
}
