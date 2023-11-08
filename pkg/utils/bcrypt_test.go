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

	check, err := ComparePassword(password, hashedPassword1)
	require.NoError(t, err)
	require.Equal(t, check, true)

	wrongPassword := RandomString(6)
	check, err = ComparePassword(wrongPassword, hashedPassword1)
	require.NoError(t, err)
	require.Equal(t, check, false)

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
