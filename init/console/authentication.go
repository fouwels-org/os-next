package console

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"os-next/init/external/otp"
	"os-next/init/external/otp/totp"
	"os-next/init/journal"

	"golang.org/x/crypto/bcrypt"
)

func generateAuthenticator(text string) (string, error) {
	const _bcryptCost = 10

	bytes, err := bcrypt.GenerateFromPassword([]byte(text), _bcryptCost)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func checkAuthenticator(hash string, text string) bool {

	if hash == "" || text == "" {
		return false
	}

	bts, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		journal.Logfln("failed to decode base64 authenticator hash: %v", err)
		return false
	}

	if len(bts) == 0 {
		journal.Logfln("authenticator hash decoded to 0 length?")
		return false
	}

	err = bcrypt.CompareHashAndPassword(bts, []byte(text))

	//lint:ignore S1008 clearer flow being verbose
	if err != nil {
		return false
	}

	return true
}

func generateTotp() (string, string, error) {

	opts := totp.GenerateOpts{
		Issuer:      "os-next",
		AccountName: "root",
		Period:      30,
		SecretSize:  32,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA512,
		Rand:        rand.Reader,
	}

	key, err := totp.Generate(opts)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	return key.Secret(), key.URL(), nil
}

func checkTotp(secret string, passcode string) (bool, error) {

	opts := totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	}

	res, err := totp.Validate(passcode, secret, time.Now(), opts)

	if err != nil {
		return false, fmt.Errorf("failed to validate OTP: %w", err)
	}

	return res, nil
}
