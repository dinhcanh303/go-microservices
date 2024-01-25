package utils

import (
	"encoding/json"

	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/pkg/errors"
)

func HandleHitCache(model interface{}, redis redis.RedisEngine, key string) error {
	byteData, exists, err := redis.Get(key)
	if exists && err == nil {
		err = json.Unmarshal(byteData, model)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal group")
		}
		return nil
	}
	return errors.Wrap(err, "miss cache")
}
