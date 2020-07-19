package redis

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"log"
	"math"
	"nokia_task/store"
	"time"
)

const (
	host           = "redis-18935.c44.us-east-1-2.ec2.cloud.redislabs.com:18935"
	password       = "S635CyTBCfYvOlcMwQspvjPX7dKRRsQ0"
	key = "set"
)

var rs *redis.Client

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:         host,
		Password:     password,
		PoolSize:     5,
		MinIdleConns: 1,
		MaxConnAge:   10 * time.Minute,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	log.Println("redis connection created")
	rs = client
}

func SetUserDataToRedis(data store.User) (err error) {
	userData, err := json.Marshal(data)
	if err != nil {
		return
	}
	p :=redis.Z{
		Member: userData,
	}
	err = rs.ZAdd(key,p).Err()
	if err != nil {
		err = errors.New("failed to set user data to redis")
		return
	}
	return
}

func GetUserDataFromRedis(start int, end int) (user []store.User, count int64, err error) {

	data, err := rs.ZRange(key,int64(start),int64(end)).Result()
	if err != nil {
		err = errors.New("failed to get user data from redis")
		return
	}
	if len(data) > 0{
		users := make([]store.User,len(data))
		for idx, d :=range data{
			err = json.Unmarshal([]byte(d), &users[idx])
			if err != nil {
				return
			}
		}
		user = users
	}
	allUsers, err := rs.ZRange(key,0, math.MaxInt64).Result()
	if err != nil {
		err = errors.New("failed to get user data from redis")
		return
	}
	count = int64(len(allUsers))

	return
}

func DeleteAllUserFromRedis() (err error) {
	err = rs.FlushAll().Err()
	return
}
