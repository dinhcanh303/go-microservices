package redis

import (
	"testing"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestRedisClient(t *testing.T) {
	redisEngine, err := connectRedis()
	require.NoError(t, err)
	require.NotEmpty(t, redisEngine)

}
func connectRedis() (RedisEngine, error) {
	err := utils.LoadFileEnvOnLocal()
	if err != nil {
		return nil, err
	}
	cfg, err := configs.NewConfigRedis()
	if err != nil {
		return nil, err
	}
	redisEngine, err := NewRedisClient(cfg)
	if err != nil {
		return nil, err
	}
	return redisEngine, nil
}
