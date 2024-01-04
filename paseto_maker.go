// Package goPasetoV4 provides an implementation of the Maker interface
// using Paseto (Platform-Agnostic Security Tokens). It includes functionality
// to create and verify tokens securely.
package goPasetoV4

import (
	"aidanwoods.dev/go-paseto"
	"crypto/rand"
	"github.com/anfen93/goPasetoV4/util"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
)

// PasetoMaker is a struct that implements the Maker interface. It provides
// methods for creating and verifying Paseto tokens.
type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey // symmetricKey is used for token encryption and decryption.
	implicit     []byte                // implicit is an additional layer of security for token operations.
}

// NewPasetoMaker initializes and returns a new PasetoMaker. It attempts to
// acquire a nonce from the PASETO_NONCE environment variable. If not available,
// it generates a random nonce. This function returns a Maker interface.
func NewPasetoMaker() Maker {
	var nonce []byte
	if envNonce := os.Getenv("PASETO_NONCE"); envNonce != "" {
		nonce = []byte(envNonce)
	} else {
		// Fallback to generate a random nonce
		nonce = make([]byte, 24)
		if _, err := rand.Read(nonce); err != nil {
			return nil
		}
	}

	return &PasetoMaker{paseto.NewV4SymmetricKey(), nonce}
}

// validateDuration validates the duration of a token. It returns true if the
// duration is valid or false if it is not. It returns an error if the duration
// is invalid.
func validateDuration(duration time.Duration) (bool, error) {
	if duration == 0 {
		return false, util.ErrDurationNotSet()
	}
	if duration < 0 {
		return false, util.ErrDurationNegative()
	}
	return true, nil
}

// CreateToken generates a new Paseto token for a given username and duration.
// It returns an encrypted token string and the Payload struct, or an error if the token generation fails.
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {

	// Validate the duration
	_, err := validateDuration(duration)
	if err != nil {
		return "", nil, err
	}
	token := paseto.NewToken() // Initializes a new Paseto token

	tokenID, err := uuid.NewRandom() // Generates a unique identifier for the token
	if err != nil {
		return "", nil, err // Returns an error if the token generation fails
	}

	// Adding necessary data to the token
	err = token.Set("id", tokenID.String())
	if err != nil {
		return "", nil, err
	}
	err = token.Set("username", username)
	if err != nil {
		return "", nil, err
	}
	// Setting issued at and expiration times for the token
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(duration))
	// Encrypts and returns the token

	encryptedToken := token.V4Encrypt(maker.symmetricKey, maker.implicit)
	payload, err := getPayloadFromToken(&token)
	if err != nil {
		return "", nil, err // Returns an error if the token generation fails
	}
	return encryptedToken, payload, nil
}

// VerifyToken checks the validity of a provided Paseto token. It returns a Payload
// struct if the token is valid or an error if it is invalid or expired.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, token, maker.implicit)
	if err != nil {
		// Handling different types of token parsing errors
		if strings.Contains(err.Error(), "expired") {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	// construct payload from token
	payload, err := getPayloadFromToken(parsedToken)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return payload, nil

}

// getPayloadFromToken extracts and returns payload data from a parsed Paseto token.
// It returns an error if the token data extraction fails.
func getPayloadFromToken(t *paseto.Token) (*Payload, error) {
	id, err := t.GetString("id")
	if err != nil {
		return nil, ErrInvalidToken
	}
	username, err := t.GetString("username")
	if err != nil {
		return nil, ErrInvalidToken
	}
	issuedAt, err := t.GetIssuedAt()
	if err != nil {
		return nil, ErrInvalidToken
	}
	expiredAt, err := t.GetExpiration()
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Payload{
		ID:        uuid.MustParse(id),
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}

// Ensuring PasetoMaker correctly implements the Maker interface
var _ Maker = (*PasetoMaker)(nil)
