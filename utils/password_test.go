package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T){
    password := RandomString(6)

    hashedPassword, err := HashPassword(password)
    require.NoError(t, err)
    require.NotEmpty(t, hashedPassword)

    err = CheckPassword(password, hashedPassword)
    require.NoError(t, err)

    wrongPassword := RandomString(6)
    err = CheckPassword(wrongPassword, hashedPassword)
    require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

    hashedPassword2, err := HashPassword(password)
    require.NoError(t, err)
    require.NotEmpty(t, hashedPassword2)
    require.NotEqual(t, hashedPassword, hashedPassword2)

    passwordVeryLong := RandomString(100)
    _, err = HashPassword(passwordVeryLong)
    require.EqualError(t, err, bcrypt.ErrPasswordTooLong.Error())
}