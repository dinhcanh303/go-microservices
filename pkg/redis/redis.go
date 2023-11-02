package redis

import (
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	redisV9 "github.com/redis/go-redis/v9"
)

const (
	_maxRetries      = 5
	_minRetryBackoff = 300 * time.Millisecond
	_maxRetryBackoff = 500 * time.Millisecond
	_dialTimeout     = 5 * time.Second
	_readTimeout     = 5 * time.Second
	_writeTimeout    = 5 * time.Second
	_minIdleConns    = 20
	_poolTimeout     = 6 * time.Second
	_poolSize        = 300
	_database        = 0
	_password        = ""
)

var ctx = context.Background()

type RedisConnString string

type redis struct {
	password string
	database int
	poolSize int

	client *redisV9.Client
}

// Client implements RedisEngine.
func (r *redis) Client() *redisV9.Client {
	return r.client
}

// Configure implements RedisEngine.
func (r *redis) Configure(opts ...Option) RedisEngine {
	for _, opt := range opts {
		opt(r)
	}
	return r
}

var _ RedisEngine = (*redis)(nil)

func NewRedisClient(url RedisConnString) (RedisEngine, error) {
	slog.Info("CONNECT_STRING", "connect string", url)
	redis := &redis{
		poolSize: _poolSize,
		database: _database,
		password: _password,
	}
	redis.client = redisV9.NewClient(
		&redisV9.Options{
			Addr:            string(url),
			Password:        redis.password,
			DB:              redis.database,
			MaxRetries:      _maxRetries,
			MinRetryBackoff: _minRetryBackoff,
			MaxRetryBackoff: _maxRetryBackoff,
			DialTimeout:     _dialTimeout,
			ReadTimeout:     _readTimeout,
			WriteTimeout:    _writeTimeout,
			MinIdleConns:    _minIdleConns,
			PoolTimeout:     _poolTimeout,
			PoolSize:        redis.poolSize,
		},
	)
	_, err := redis.client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect redis")
	}
	slog.Info("📫 connected to redis 🎉")
	return redis, nil
}
