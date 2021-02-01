package redis

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v7"
)

const (
	CACHE_WEEK_TTL = 604800 * time.Second // 7天
	CACHE_DAY_TTL  = 86400 * time.Second  // 1天
	CACHE_HOUR_TTL = 3600 * time.Second   // 1小时
	CACHE_MIN_TTL  = 60 * time.Second     // 1分钟
)

var ErrNoData = errors.New("data is empty")
var ErrClosed = redis.ErrClosed

type Redis struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{
		client: client,
	}
}

func NewRedisWithName(name string) *Redis {
	return &Redis{
		client: Client(name),
	}
}

func IsNil(err error) bool {
	return errors.Is(err, redis.Nil)
}

func (b *Redis) Client() *redis.Client {
	return b.client
}

func (b *Redis) HGetAll(key string, m interface{}) error {
	info, err := b.HGetAllMap(key)
	if err != nil {
		return err
	}
	return scanStruct(info, m)
}

func (b *Redis) HGetAllMap(key string) (map[string]string, error) {
	data, err := b.client.HGetAll(key).Result()
	if len(data) == 0 {
		return nil, ErrNoData
	}
	return data, err
}

// HMSet support gosql.Model and map[string]interface{}
func (b *Redis) HMSet(key string, m interface{}) error {
	var mm map[string]interface{}
	var ok bool
	ref := reflect.TypeOf(m)
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	switch ref.Kind() {
	case reflect.Struct:
		mm = structToMapInterface(m)
	case reflect.Map:
		mm, ok = m.(map[string]interface{})
		if !ok {
			ms, ok := m.(map[string]string)
			if !ok {
				return errors.New("value must is map[string]interafce{}")
			}
			mm = make(map[string]interface{})
			for key, value := range ms {
				mm[key] = value
			}
		}
	default:
		return errors.New(fmt.Sprintf("cannot convert from %s", ref))
	}

	return b.hMSet(key, mm)
}

func (b *Redis) hMSet(key string, m map[string]interface{}) error {
	err := b.client.HMSet(key, m).Err()
	if err != nil {
		return err
	}
	return b.client.Expire(key, CACHE_DAY_TTL).Err()
}

func (b *Redis) HGet(key string, field string) (string, error) {
	return b.client.HGet(key, field).Result()
}

func (b *Redis) HSet(key string, field string, value interface{}) error {
	return b.client.HSet(key, field, value).Err()
}

func (b *Redis) HDel(key string, field string) error {
	return b.client.HDel(key, field).Err()
}

func (b *Redis) HExists(key string, field string) (bool, error) {
	return b.client.HExists(key, field).Result()
}

func (b *Redis) HIncrBy(key string, field string, incr int64) (int64, error) {
	return b.client.HIncrBy(key, field, incr).Result()
}

func (b *Redis) Exists(key string) bool {
	n, err := b.client.Exists(key).Result()

	if err != nil || n != 1 {
		return false
	}

	return true
}

func (b *Redis) LPop(key string) (string, error) {
	return b.client.LPop(key).Result()
}

func (b *Redis) RPop(key string) (string, error) {
	return b.client.RPop(key).Result()
}

func (b *Redis) LPush(key string, values ...interface{}) (int64, error) {
	return b.client.LPush(key, values...).Result()
}

func (b *Redis) RPush(key string, values ...interface{}) (int64, error) {
	return b.client.RPush(key, values...).Result()
}

func (b *Redis) Del(key string) error {
	return b.client.Del(key).Err()
}

// Set Our specification is that all keys must have an expiration time, so the Set key will expire in one day by default
func (b *Redis) Set(key, val string) error {
	return b.client.Set(key, val, CACHE_DAY_TTL).Err()
}

func (b *Redis) SetEX(key, val string, expiration time.Duration) error {
	return b.client.Set(key, val, expiration).Err()
}

func (b *Redis) SetNX(key, val string, expiration time.Duration) (bool, error) {
	return b.client.SetNX(key, val, expiration).Result()
}

func (b *Redis) Get(key string) (string, error) {
	return b.client.Get(key).Result()
}

func (b *Redis) Incr(key string) (int64, error) {
	return b.client.Incr(key).Result()
}

func (b *Redis) Decr(key string) (int64, error) {
	return b.client.Decr(key).Result()
}

func (b *Redis) IncrBy(key string, value int64) (int64, error) {
	return b.client.IncrBy(key, value).Result()
}

func (b *Redis) DecrBy(key string, value int64) (int64, error) {
	return b.client.DecrBy(key, value).Result()
}

func (b *Redis) Expire(key string, expiration time.Duration) (bool, error) {
	return b.client.Expire(key, expiration).Result()
}

func (b *Redis) SAdd(key string, members ...interface{}) (int64, error) {
	return b.client.SAdd(key, members...).Result()
}

func (b *Redis) TTL(key string) time.Duration {
	return b.client.TTL(key).Val()
}
