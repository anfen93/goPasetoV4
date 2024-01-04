package goPasetoV4

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker := NewPasetoMaker()

	t.Run("CreateToken with Valid Inputs", func(t *testing.T) {
		username := "testuser"
		duration := time.Minute

		token, payload, err := maker.CreateToken(username, duration)
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.NotNil(t, payload)
		require.Equal(t, username, payload.Username)
	})

	t.Run("VerifyToken with Valid Token", func(t *testing.T) {
		username := "testuser"
		duration := time.Minute

		token, _, err := maker.CreateToken(username, duration)
		require.NoError(t, err)

		payload, err := maker.VerifyToken(token)
		require.NoError(t, err)
		require.NotNil(t, payload)
		require.Equal(t, username, payload.Username)
	})

	t.Run("VerifyToken with Valid Token and Expiration in 24h", func(t *testing.T) {
		username := "testuser"
		//duration is expressed as time.Duration of 24h
		duration := time.Duration(24) * time.Hour

		token, _, err := maker.CreateToken(username, duration)
		require.NoError(t, err)

		payload, err := maker.VerifyToken(token)
		require.NoError(t, err)
		require.NotNil(t, payload)
		require.Equal(t, username, payload.Username)
	})

	t.Run("VerifyToken with Expired Token", func(t *testing.T) {
		username := "testuser"
		//two minutes ago
		expiredDuration := -time.Minute * 2

		token, _, err := maker.CreateToken(username, expiredDuration)

		payload, err := maker.VerifyToken(token)
		require.Error(t, err)
		require.Nil(t, payload)
	})

	t.Run("CreateToken with Zero Duration", func(t *testing.T) {
		username := "testuser"

		token, payload, err := maker.CreateToken(username, 0)
		require.Error(t, err)
		require.Empty(t, token)
		require.Nil(t, payload)
	})

	t.Run("CreateToken with Negative Duration", func(t *testing.T) {
		username := "testuser"

		token, payload, err := maker.CreateToken(username, -time.Minute)
		require.Error(t, err)
		require.Empty(t, token)
		require.Nil(t, payload)
	})

	t.Run("VerifyToken with Altered Token", func(t *testing.T) {
		username := "testuser"
		duration := time.Minute

		token, _, err := maker.CreateToken(username, duration)
		require.NoError(t, err)

		alteredToken := token + "something-extra"
		payload, err := maker.VerifyToken(alteredToken)
		require.Error(t, err)
		require.Nil(t, payload)
	})
}
