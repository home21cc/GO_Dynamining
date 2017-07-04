package dJwt

import "github.com/garyburd/redigo/redis"

type RedisClient  struct {
	conn redis.Conn
}

var instanceRedisClient *RedisClient = nil

func Connect() (conn *RedisClient) {
	if instanceRedisClient == nil {
		instanceRedisClient = new(RedisClient)
		var err error
		instanceRedisClient.conn, err = redis.Dial("tcp", ":6379")
		if err != nil {
			panic(err)
		}

		if _, err := instanceRedisClient.conn.Do("AUTH", "Brainattica"); err != nil {
			instanceRedisClient.conn.Close()
			panic(err)
		}
	}
	return instanceRedisClient
}

func(redisClient *RedisClient) SetValue(key, value string, expiration ...interface{}) error {
	_, err := redisClient.conn.Do("SET", key, value)
	if err == nil && expiration != nil {
		redisClient.conn.Do("EXPIRE", key, expiration[0])
	}
	return err
}

func(redisClient *RedisClient) GetValue(key string) (interface{}, error) {
	return redisClient.conn.Do("GET", key)
}
