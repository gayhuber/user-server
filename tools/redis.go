package tools

import (
	"github.com/gomodule/redigo/redis"
	"time"
	"user-server/config"
)

var (
	redisClient  *redis.Pool
	REDIS_HOST   string
	REDIS_DB     int
	REDIS_AUTH   string
	MAX_ACTIVE   int
	MAX_IDLE     int
	IDLE_TIMEOUT int64
)

func init() {
	REDIS_HOST = config.Conf.Redis.Host
	REDIS_DB = 0
	REDIS_AUTH = "abc"
	MAX_ACTIVE = 10
	MAX_IDLE = 1

	/**
	*@MaxIdle 最大空闲链接
	*@MaxActive 最大活跃链接
	*@IdleTimeout 自动超时时间
	 */
	redisClient = &redis.Pool{
		MaxIdle:     MAX_IDLE,
		MaxActive:   MAX_ACTIVE,
		IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}

			if REDIS_AUTH != "" {
				c.Do("AUTH", REDIS_AUTH)
			}
			c.Do("SELECT", REDIS_DB)
			return c, nil
		},
	}
}

func GetRedis() redis.Conn {
	return redisClient.Get()
}
