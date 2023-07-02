package token

import (
	"testing"
	"time"

	"github.com/morka17/shiny_bank/v1/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	assert.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, _, err := maker.CreateToken(username, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, payload)

	assert.NotEmpty(t,payload.ID)
	assert.Equal(t, username, payload.Username)
	assert.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
	assert.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	assert.NoError(t, err)

	username := utils.RandomOwner()
	duration := -time.Minute

	token, _, err := maker.CreateToken(username, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrExpiredToken.Error())
	assert.Nil(t, payload)

}

