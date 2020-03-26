package server

import (
	"time"

	"github.com/go-redis/redis"
	//log "github.com/sirupsen/logrus"
	"github.com/ulule/limiter"
	sredis "github.com/ulule/limiter/drivers/store/redis"
)

type Redis struct {
	MaxRetry    int    `yaml:"max_retry"`
	RateLimiter string `yaml:"rate_limiter"`
	RedisServer string `yaml:"redis_server"`
}

// TODO can make this a pointer
func initalize_redis(redis_config *Redis) (*limiter.Store, *limiter.Rate) {
	// create rate limiter
	rate, err := limiter.NewRateFromFormatted(redis_config.RateLimiter)
	if err != nil {
		log.Fatal(err)
	}

	option, err := redis.ParseURL(redis_config.RedisServer + "/0")
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(option)
	pong, err := client.Ping().Result()

	if err != nil {
		log.Info(err)
		for true {
			pong, err = client.Ping().Result()

			if err == nil {
				log.Info("Successful Ping", pong)
				break
			}
			time.Sleep(10 * time.Second)
			log.Info(err)
		}
	}

	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "limiter_http",
		MaxRetry: redis_config.MaxRetry,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Successfully connected to redis server")

	return &store, &rate
}
