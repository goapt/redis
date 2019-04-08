package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/ilibs/gosql"
	"github.com/verystar/golib/convert"
)

const (
	CACHE_WEEK_TTL = 604800 * time.Second //7天
	CACHE_DAY_TTL  = 86400 * time.Second  //1天
	CACHE_HOUR_TTL = 3600 * time.Second   //1小时
	CACHE_MIN_TTL  = 60 * time.Second     //1分钟
)

type BaseRedis struct {
	Client func() *redis.Client
}

func NewBaseRedis(name string) *BaseRedis {
	return &BaseRedis{
		Client: func() *redis.Client {
			return Client(name)
		},
	}
}

func (b *BaseRedis) HGetAll(key string, m gosql.IModel) error {
	info, err := b.Client().HGetAll(key).Result()
	if err != nil {
		return err
	}
	return ScanStruct(info, m)
}

func (b *BaseRedis) HGetAllMap(key string) (map[string]string, error) {
	return b.Client().HGetAll(key).Result()
}

func (b *BaseRedis) HMSet(key string, m gosql.IModel) error {
	var fn = func() map[string]interface{} {
		return convert.StructToMapInterface(m)
	}
	return b.hMSet(key, fn)
}

func (b *BaseRedis) HMSetMap(key string, m map[string]interface{}) error {
	var fn = func() map[string]interface{} {
		return m
	}
	return b.hMSet(key, fn)
}

func (b *BaseRedis) hMSet(key string, fn func() map[string]interface{}) error {
	err := b.Client().HMSet(key, fn()).Err()
	if err != nil {
		return err
	}
	return b.Client().Expire(key, CACHE_DAY_TTL).Err()
}

func (b *BaseRedis) HGet(key string, field string) (string, error) {
	return b.Client().HGet(key, field).Result()
}

func (b *BaseRedis) HSet(key string, field string, value interface{}) error {
	return b.Client().HSet(key, field, value).Err()
}

func (b *BaseRedis) HDel(key string, field string) error {
	return b.Client().HDel(key, field).Err()
}

func (b *BaseRedis) HExists(key string, field string) (bool, error) {
	return b.Client().HExists(key, field).Result()
}

func (b *BaseRedis) HIncrBy(key string, field string, incr int64) (int64, error) {
	return b.Client().HIncrBy(key, field, incr).Result()
}

func (b *BaseRedis) Exists(key string) bool {
	n, err := b.Client().Exists(key).Result()

	if err != nil {
		return false
	}

	if n != 1 {
		return false
	}

	return true
}

func (b *BaseRedis) LPop(key string) (string, error) {
	return b.Client().LPop(key).Result()
}

func (b *BaseRedis) RPop(key string) (string, error) {
	return b.Client().RPop(key).Result()
}

func (b *BaseRedis) LPush(key string, values ...interface{}) (int64, error) {
	return b.Client().LPush(key, values...).Result()
}

func (b *BaseRedis) RPush(key string, values ...interface{}) (int64, error) {
	return b.Client().RPush(key, values...).Result()
}

func (b *BaseRedis) Del(key string) error {
	return b.Client().Del(key).Err()
}

func (b *BaseRedis) Set(key, val string) error {
	return b.Client().Set(key, val, CACHE_DAY_TTL).Err()
}

func (b *BaseRedis) SetEX(key, val string, expiration time.Duration) error {
	return b.Client().SetXX(key, val, expiration).Err()
}

func (b *BaseRedis) Get(key string) (string, error) {
	return b.Client().Get(key).Result()
}

func (b *BaseRedis) Incr(key string) (int64, error) {
	return b.Client().Incr(key).Result()
}

func (b *BaseRedis) Decr(key string) (int64, error) {
	return b.Client().Decr(key).Result()
}

func (b *BaseRedis) IncrBy(key string, value int64) (int64, error) {
	return b.Client().IncrBy(key, value).Result()
}

func (b *BaseRedis) DecrBy(key string, value int64) (int64, error) {
	return b.Client().DecrBy(key, value).Result()
}

func (b *BaseRedis) Expire(key string, expiration time.Duration) (bool, error) {
	return b.Client().Expire(key, expiration).Result()
}

func (b *BaseRedis) SAdd(key string, members ...interface{}) (int64, error) {
	return b.Client().SAdd(key, members...).Result()
}

func (b *BaseRedis) TTL(key string) time.Duration {
	return b.Client().TTL(key).Val()
}
