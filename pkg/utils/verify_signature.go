package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

func VerifyMidtransSignature(orderID, statusCode, grossAmount, signatureKey, serverKey string) bool {
	raw := orderID + statusCode + grossAmount + serverKey
	sum := sha512.Sum512([]byte(raw))
	return hex.EncodeToString(sum[:]) == signatureKey
}
