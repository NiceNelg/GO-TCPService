package redispool

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

/**
 * @Function 建立redis连接池
 * @Auther Nelg
 * @Date 2019.05.30
 */
func NewPool(server string, password string, db string) *redis.Pool {
	return &redis.Pool{
		/*最大的空闲连接数*/
		MaxIdle:     15,
		/*最大的激活连接数*/
		MaxActive:   500,
		/*空闲连接超时关闭秒数*/
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
