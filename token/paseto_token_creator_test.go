package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wenealves10/gobank/utils"
)

func TestPasetoTokenCreator(t *testing.T) {
    tokenCreator, err := NewPasetoTokenCreator(utils.RandomString(32))
    require.NoError(t, err)

    username := utils.RandomOwner()
    duration := time.Minute

    issuedAt := time.Now()
    expiresAt := issuedAt.Add(duration)

    token, err := tokenCreator.CreateToken(username, duration)
    require.NoError(t, err)
    require.NotEmpty(t, token)

    payload, err := tokenCreator.VerifyToken(token)
    require.NoError(t, err)
    require.NotNil(t, payload)

    require.Equal(t, username, payload.Username)
    require.NotZero(t, payload.ID)
    require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
    require.WithinDuration(t, expiresAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoTokenCreator(t *testing.T) {
    tokenCreator, err := NewPasetoTokenCreator(utils.RandomString(32))
    require.NoError(t, err)

    token, err := tokenCreator.CreateToken(utils.RandomOwner(), -time.Minute)
    require.NoError(t, err)
    require.NotEmpty(t, token)

    payload, err := tokenCreator.VerifyToken(token)
    require.Error(t, err)
    require.EqualError(t, err, ErrExpiredToken.Error())
    require.Nil(t, payload)
}