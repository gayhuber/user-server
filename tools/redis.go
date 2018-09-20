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

// GetRedis 从连接池获取一个 redis 链接
func GetRedis() RedisBody {
	return RedisBody{
		Conn: redisClient.Get(),
	}
}

// RedisBody 一个单独的 redis 包
type RedisBody struct {
	Conn   redis.Conn
	value  interface{}
	key    interface{}
	expire int
	Error  error
	resp   interface{}
}

// Set 添加 key-value
func (rb *RedisBody) Set(key interface{}, value interface{}, expire int) *RedisBody {
	resp, err := rb.Conn.Do("SET", key, value, "EX", expire)
	rb.resp = resp
	rb.Error = err
	return rb
}

// Get 获取 key
func (rb *RedisBody) Get(key interface{}) (resp interface{}, err error) {
	resp, err = rb.Conn.Do("GET", key)
	return
}

// Del 删除
func (rb *RedisBody) Del(key ...interface{}) *RedisBody {
	if len(key) == 0 {
		return rb
	}
	resp, err := rb.Conn.Do("Del", key)
	rb.resp = resp
	rb.Error = err
	return rb
}

// Expire 设置时间
func (rb *RedisBody) Expire(key interface{}, expire int) *RedisBody {
	resp, err := rb.Conn.Do("EXPIRE", key, expire)
	rb.resp = resp
	rb.Error = err
	return rb
}

// Exist 验证是否存在
func (rb *RedisBody) Exist(key interface{}) bool {
	resp, err := rb.Conn.Do("EXISTS", key)
	rb.resp = resp
	rb.Error = err
	if err != nil {
		return false
	}
	return resp.(int64) == 1
}

// Close 关闭链接
func (rb *RedisBody) Close() {
	rb.Conn.Close()
}
