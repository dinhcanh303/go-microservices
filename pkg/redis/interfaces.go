package redis

import (
	redisV9 "github.com/redis/go-redis/v9"
)

type RedisEngine interface {
	Configure(...Option) RedisEngine
	Client() *redisV9.Client
}
