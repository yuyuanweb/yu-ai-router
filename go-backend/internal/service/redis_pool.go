package service

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/yupi/airouter/go-backend/internal/config"
)

const (
	defaultRedisDialTimeout = 2 * time.Second
	defaultRedisReadTimeout = 2 * time.Second
	defaultRedisIdleTimeout = 60 * time.Second
	defaultRedisMaxIdle     = 5
	defaultRedisMaxActive   = 30
)

func NewRedisPool(cfg *config.Config) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     defaultRedisMaxIdle,
		MaxActive:   defaultRedisMaxActive,
		IdleTimeout: defaultRedisIdleTimeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			options := []redis.DialOption{
				redis.DialConnectTimeout(defaultRedisDialTimeout),
				redis.DialReadTimeout(defaultRedisReadTimeout),
				redis.DialWriteTimeout(defaultRedisReadTimeout),
			}
			if cfg.RedisUsername != "" {
				options = append(options, redis.DialUsername(cfg.RedisUsername))
			}
			if cfg.RedisPassword != "" {
				options = append(options, redis.DialPassword(cfg.RedisPassword))
			}
			if cfg.RedisDB > 0 {
				options = append(options, redis.DialDatabase(cfg.RedisDB))
			}
			return redis.Dial("tcp", cfg.RedisAddr, options...)
		},
		TestOnBorrow: func(conn redis.Conn, _ time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}
