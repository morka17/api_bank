package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err, fmt.Sprint(err))
	assert.NotEmpty(t, hashedPassword)

	err = CheckPassword(password, hashedPassword)
	assert.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword)
	assert.Error(t, err, fmt.Sprint(err))


	hashedPassword2, err := HashPassword(password)
	assert.NoError(t, err, fmt.Sprint(err))
	assert.NotEmpty(t, hashedPassword2)
	assert.NotEqual(t,hashedPassword, hashedPassword2)
}