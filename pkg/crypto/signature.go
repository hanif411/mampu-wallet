package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateSignature(WalletID uint, balance int64, secret string) string {
	data := fmt.Sprintf("%d:%d", WalletID, balance)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifySignature(WalletID uint, balance int64, secret, dbsignature string) bool {
	generated := GenerateSignature(WalletID, balance, secret)
	return hmac.Equal([]byte(generated), []byte(dbsignature))
}
