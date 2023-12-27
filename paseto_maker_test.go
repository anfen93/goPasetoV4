package goPasetoV4

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMakerTokenLifecycle(t *testing.T) {
	// Testing token creation and validation lifecycle
	maker := NewPasetoMaker()
	require.NotNil(t, maker, "Maker should not be nil")

	username := "testuser"
	duration := time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err, "Token creation should not error")
	require.NotEmpty(t, token, "Token should not be empty")

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err, "Token verification should not error")
	require.NotEmpty(t, payload, "Payload should not be empty")

	require.NotZero(t, payload.ID, "Payload ID should not be zero")
	require.Equal(t, username, payload.Username, "Username should match")
	require.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second, "IssuedAt should be recent")
	require.WithinDuration(t, time.Now().Add(duration), payload.ExpiredAt, time.Second, "ExpiredAt should be correct")
}

func TestPasetoMakerExpiredToken(t *testing.T) {
	// Testing behavior with an expired token
	maker := NewPasetoMaker()
	require.NotNil(t, maker, "Maker should not be nil")

	token, err := maker.CreateToken("expireduser", -time.Minute)
	require.NoError(t, err, "Expired token creation should not error")
	require.NotEmpty(t, token, "Expired token should not be empty")

	payload, err := maker.VerifyToken(token)
	require.Error(t, err, "Expired token verification should error")
	require.EqualError(t, err, ErrExpiredToken.Error(), "Error should be ErrExpiredToken")
	require.Nil(t, payload, "Payload should be nil for an expired token")
}

func TestPasetoMakerInvalidToken(t *testing.T) {
	// Testing behavior with an invalid token
	maker := NewPasetoMaker()
	require.NotNil(t, maker, "Maker should not be nil")

	payload, err := maker.VerifyToken("invalidtoken")
	require.Error(t, err, "Invalid token verification should error")
	require.EqualError(t, err, ErrInvalidToken.Error(), "Error should be ErrInvalidToken")
	require.Nil(t, payload, "Payload should be nil for an invalid token")
}
