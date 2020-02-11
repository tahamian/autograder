package server

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ulule/limiter"
	sredis "github.com/ulule/limiter/drivers/store/redis"
)

type Redis struct {
	Redis_server string
	Rate_limit   string
	Max_retry 	 int
}

// TODO can make this a pointer
func initalize_redis(redis Redis) (limiter.Store, limiter.Rate) {
	// create rate limiter
	rate, err := limiter.NewRateFromFormatted(redis.Rate_limit)
	if err != nil {
		log.Fatal(err)
	}

	// Create a redis client.
	option, err := redis.ParseURL(redis.Redis_server + "/0")
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(option)

	pong, err := client.Ping().Result()

	// redis_server := strings.Replace(c.Redis_server, "redis", "http", 1)
	if err != nil {
		log.Info(err)
		for true {
			pong, err = client.Ping().Result()

			if err == nil {
				log.Info("Successful Ping", pong)
				break
			}
			time.Sleep(10 * time.Second)
			fmt.Println(err)
		}
	}

	// Create a store with the redis client.
	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "limiter_http",
		MaxRetry: redis.Max_retry,
	})
	if err != nil {
		log.Fatal(err)
	}

	return store, rate
}
