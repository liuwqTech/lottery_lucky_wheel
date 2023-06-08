package datasource

import (
	"Lucky_Wheel/conf"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"sync"
	"time"
)

var rdsLock sync.Mutex

var cacheInstance *RedisConn

type RedisConn struct {
	pool      *redis.Pool
	showDebug bool
}

func (rds *RedisConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := rds.pool.Get()
	defer conn.Close()
	t1 := time.Now().UnixNano()
	reply, err = conn.Do(commandName, args...)
	if err != nil {
		e := conn.Err()
		if e != nil {
			log.Println("rdsHelper.Do", err, e)
		}
	}
	t2 := time.Now().UnixNano()
	if rds.showDebug {
		fmt.Printf("[redis][info][%dus]cmd=%s,err=%s,arg=%s,reply=%s\n",
			(t2-t1)/1000,
			commandName,
			err,
			args,
			reply,
		)
	}
	return reply, err
}

func (rds *RedisConn) ShowDebug(b bool) {
	rds.showDebug = b
}
func InstanceCache() *RedisConn {
	if cacheInstance != nil {
		return cacheInstance
	}
	rdsLock.Lock()
	defer rdsLock.Unlock()
	if cacheInstance != nil {
		return cacheInstance
	}
	return NewCache()
}

func NewCache() *RedisConn {
	pool := redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d",
				conf.RedisCache.Host,
				conf.RedisCache.Port,
			))
			if err != nil {
				log.Fatal("rdsHelper.NewCache Dial error=", err)
				return nil, err
			}
			return c, nil
		},
		DialContext: nil,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         10000,
		MaxActive:       10000,
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}
	instance := &RedisConn{
		pool: &pool,
	}
	instance.ShowDebug(true)
	cacheInstance = instance
	return instance
}
