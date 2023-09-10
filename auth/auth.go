package auth

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
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
	if err != nil {
		return false, err
	}

	// client
	clientK, err := srp.CompleteHandshake(pubA, privA, []byte(strings.ToLower(username)), []byte(password), salt, B)
	if err != nil {
		return false, err
	}

	// server
	if subtle.ConstantTimeCompare(clientK, K) != 1 {
		return false, nil
	}

	return true, nil
}

func VerifyLegacyPassword(username, password, b64_verifier string) bool {
	digest := sha1.Sum([]byte(username + password))
	b64 := base64.RawStdEncoding.EncodeToString(digest[:])
	return b64 == b64_verifier
}

func CreateAuth(username, password string) (salt, verifier []byte, err error) {
	return srp.NewClient([]byte(strings.ToLower(username)), []byte(password))
}

var ValidPlayernameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_]*$`)

func ValidateUsername(username string) error {
	if len(username) == 0 {
		return fmt.Errorf("playername empty")
	}
	if len(username) > 20 {
		return fmt.Errorf("playername too long")
	}
	if !ValidPlayernameRegex.Match([]byte(username)) {
		return fmt.Errorf("playername can only contain chars a to z, A to Z, 0 to 9 and -, _")
	}
	return nil
}
