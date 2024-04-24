package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRandomBytes(t *testing.T) {
	key, err := GenerateRandomBytes(64)

	require.NoError(t, err)
	require.NotEmpty(t, key)
	require.Equal(t, len(key), 64)
}
func TestGenerateRandomHexBytes(t *testing.T) {
	key, err := GenerateRandomHexBytes(64)
	fmt.Println(key)
	require.NoError(t, err)
	require.NotEmpty(t, key)
	require.Equal(t, len(key), 128)
}
