package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/HimbeerserverDE/srp"
)

func CreateDBPassword(salt, verifier []byte) string {
	return fmt.Sprintf(
		"#1#%s#%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(verifier),
	)
}

func ParseDBPassword(password string) (salt, verifier []byte, err error) {
	parts := strings.Split(password, "#")
	if len(parts) != 4 {
		return nil, nil, errors.New("invalid delimiter count")
	}
	if parts[1] != "1" {
		return nil, nil, errors.New("invalid version: " + parts[1])
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, nil, err
	}
	verifier, err = base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return nil, nil, err
	}

	return salt, verifier, nil
}

func VerifyAuth(username, password string, salt, verifier []byte) (bool, error) {
	// client
	pubA, privA, err := srp.InitiateHandshake()
	if err != nil {
		return false, err
	}

	// server
	B, _, K, err := srp.Handshake(pubA, verifier)

	// client
	clientK, err := srp.CompleteHandshake(pubA, privA, []byte(username), []byte(password), salt, B)
	if err != nil {
		return false, err
	}

	// server
	if subtle.ConstantTimeCompare(clientK, K) != 1 {
		return false, nil
	}

	return true, nil
}

func CreateAuth(username, password string) (salt, verifier []byte, err error) {
	return srp.NewClient([]byte(username), []byte(password))
}
