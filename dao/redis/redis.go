package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

// 初始化连接
func Init() (err error) {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port"),
		),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
		//Addr:     "localhost:8088",
		//Password: "",
		//DB:       0,
		//PoolSize: 100,
	})

	_, err = client.Ping().Result()
	return
}

func Close() {
	_ = client.Close()
}
