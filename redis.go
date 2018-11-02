package redis

import (
	"github.com/go-redis/redis"
	"log"
	"strings"
	"time"
)

var (
	redisList map[string]*redis.Client
	errs      []string
)

type Config struct {
	Server       string
	Password     string
	DB           int
	MaxRetries   int
	DialTimeout  int `json:"dial_timeout" toml:"dial_timeout"`
	ReadTimeout  int `json:"read_timeout" toml:"read_timeout"`
	WriteTimeout int `json:"write_timeout" toml:"write_timeout"`
}

func Client(name ... string) *redis.Client {
	key := "default"
	if name != nil {
		key = name[0]
	}

	client, ok := redisList[key]
	if !ok {
		log.Fatalf("[redis] the redis client `%s` is not configured", key)
	}

	return client
}

func Connect(configs map[string]Config) {
	defer func() {
		if len(errs) > 0 {
			log.Fatal("[redis] " + strings.Join(errs, "\n"))
		}
	}()

	redisList = make(map[string]*redis.Client)
	for name, conf := range configs {
		r := newRedis(&conf)
		log.Println("[redis] connect:" + conf.Server)

		_, err := r.Ping().Result()
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}

		client := newRedis(&conf)

		if r, ok := redisList[name]; ok {
			redisList[name] = client
			r.Close()
		} else {
			redisList[name] = client
		}
	}
}

// 创建 redis pool
func newRedis(conf *Config) *redis.Client {

	options := &redis.Options{
		Addr:     conf.Server,
		Password: conf.Password, // no password set
		DB:       conf.DB,       // use default DB
	}

	if conf.MaxRetries > 0 {
		options.MaxRetries = conf.MaxRetries
	}

	if conf.DialTimeout > 0 {
		options.DialTimeout = time.Duration(conf.DialTimeout) * time.Second
	}

	if conf.ReadTimeout > 0 {
		options.ReadTimeout = time.Duration(conf.ReadTimeout) * time.Second
	}

	if conf.WriteTimeout > 0 {
		options.WriteTimeout = time.Duration(conf.WriteTimeout) * time.Second
	}

	client := redis.NewClient(options)
	return client
}
