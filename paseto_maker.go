package goPasetoV4

import (
	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"strings"
	"time"
)

// PasetoMaker is a struct that implements Maker interface
type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
	implicit     []byte
}

func NewPasetoMaker() Maker {
	return &PasetoMaker{paseto.NewV4SymmetricKey(), []byte("my implicit nonce")}
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// create paseto token
	token := paseto.NewToken()
	// Create uuid for token id
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	// add data to the token.
	err = token.Set("id", tokenID.String())
	if err != nil {
		return "", err
	}
	err = token.Set("username", username)
	if err != nil {
		return "", err
	}
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(duration))
	return token.V4Encrypt(maker.symmetricKey, maker.implicit), nil
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, token, maker.implicit)
	if err != nil {
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

var _ Maker = (*PasetoMaker)(nil)
