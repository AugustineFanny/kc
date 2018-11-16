package utils

import (
	"time"
	"github.com/garyburd/redigo/redis"
)

type RedisCache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	password string
}

func (rc *RedisCache) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

func (rc *RedisCache) Set(key string, val interface{}, timeout ...int64) error {
	if len(timeout) > 0 {
		if _, err := rc.Do("SETEX", key, timeout[0], val); err != nil {
			return err
		}
	} else {
		if _, err := rc.Do("SET", key, val); err != nil {
			return err
		}
	}
	return nil
}

func (rc *RedisCache) Get(key string) (reply interface{},err error) {
	return rc.Do("GET", key)
}

func (rc *RedisCache) GetString(key string) (reply string, err error) {
	return redis.String(rc.Do("GET", key))
}

func (rc *RedisCache) Hmset(key string, val interface{}) error {
	if _, err := rc.Do("HMSET", redis.Args{}.Add(key).AddFlat(val)...); err != nil {
		return err
	}
	return nil
}

func (rc *RedisCache) HgetFloat64(key, field string) (reply float64, err error) {
	return redis.Float64(rc.Do("HGET", key, field))
}

func (rc *RedisCache) HgetallFloat64(key string) (map[string]float64, error) {
	values, err := redis.Values(rc.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	var k string
	var v float64
	reply := map[string]float64{}
	for len(values) > 0 {
		values, err = redis.Scan(values, &k, &v)
		if err != nil {
			return nil, err
		}
		reply[k] = v
	}
	return reply, nil
}

func (rc *RedisCache) LrangeString(key string, start, stop int) ([]string, error) {
	values, err := redis.Values(rc.Do("LRANGE", key, start, stop))
	if err != nil {
		return []string{}, err
	}
	var raw []string
	if err := redis.ScanSlice(values, &raw); err != nil {
		return []string{}, err
	}
	return raw, nil
}

func (rc *RedisCache) Delete(key string) error {
	var err error
	if _, err = rc.Do("DEL", key); err != nil {
		return err
	}
	return err
}

func (rc *RedisCache) IsExist(key string) bool {
	v, err := redis.Bool(rc.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

func (rc *RedisCache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}
