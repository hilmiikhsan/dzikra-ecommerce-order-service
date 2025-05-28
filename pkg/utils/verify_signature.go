package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

func VerifyMidtransSignature(orderID, statusCode, grossAmount, signatureKey, serverKey string) bool {
	raw := fmt.Sprintf("%s%s%s%s", orderID, statusCode, grossAmount, serverKey)
	sum := sha512.Sum512([]byte(raw))
	return hex.EncodeToString(sum[:]) == signatureKey
}
