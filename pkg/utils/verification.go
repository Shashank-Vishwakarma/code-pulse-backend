package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
)

func GenerateVerificationCode() string {
	verificationCode, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		logrus.Errorf("Error generating the verification code: %v", err)
		return ""
	}
	return fmt.Sprintf("%06d", verificationCode.Int64())
}
