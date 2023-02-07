package redis

import (
	"evergreen/settings"
	"fmt"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Init() error {
	poolSize := settings.Conf.RedisConfig.PoolSize
	password := settings.Conf.RedisConfig.Password
	host := settings.Conf.RedisConfig.Host
	port := settings.Conf.RedisConfig.Port
	db := settings.Conf.RedisConfig.DB
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password, // 密码
		DB:       db,       // 数据库
		PoolSize: poolSize, // 连接池大小
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
