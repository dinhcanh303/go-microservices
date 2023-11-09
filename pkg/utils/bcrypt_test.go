package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBcrypt(t *testing.T) {
	password := RandomString(8)
	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	check, err := ComparePassword(hashedPassword1, password)
	require.NoError(t, err)
	require.Equal(t, check, true)

	wrongPassword := RandomString(6)
	_, err = HashPassword(wrongPassword)
	require.EqualError(t, err, "password must be less than 8 characters")

	wrongPassword2 := RandomString(8)
	check, err = ComparePassword(hashedPassword1, wrongPassword2)
	require.NoError(t, err)
	require.Equal(t, check, false)

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
