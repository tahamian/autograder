package api

import (
	"autograder/internal/config"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/ulule/limiter"
	sredis "github.com/ulule/limiter/drivers/store/redis"
)

func initializeRedis(log *logrus.Logger, cfg *config.RedisConfig) (*limiter.Store, *limiter.Rate, error) {
	rate, err := limiter.NewRateFromFormatted(cfg.RateLimiter)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid rate format %q: %w", cfg.RateLimiter, err)
	}

	option, err := redis.ParseURL(cfg.RedisServer + "/0")
	if err != nil {
		return nil, nil, fmt.Errorf("invalid redis URL %q: %w", cfg.RedisServer, err)
	}

	client := redis.NewClient(option)

	const maxAttempts = 10
	for i := 0; i < maxAttempts; i++ {
		if _, err := client.Ping().Result(); err == nil {
			break
		} else if i == maxAttempts-1 {
			return nil, nil, fmt.Errorf("redis unreachable after %d attempts: %w", maxAttempts, err)
		} else {
			log.WithError(err).Warn("redis ping failed, retrying...")
			time.Sleep(2 * time.Second)
		}
	}

	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "limiter_http",
		MaxRetry: cfg.MaxRetry,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("creating redis store: %w", err)
	}

	log.Info("connected to redis")
	return &store, &rate, nil
}
