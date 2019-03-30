package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/wbsifan/devgo/errors"

	"github.com/wbsifan/devgo/helper"

	"github.com/go-redis/redis"
)

type (
	Cache interface {
		Get(key string, v interface{}) error
		Set(key string, v interface{}, exp time.Duration) error
		Del(keys ...string) error
	}

	RedisCache struct {
		Disable bool
		Debug   bool
		rc      *redis.Client
		prefix  string
	}
)

var DisableError = errors.New("Cache disabled")

func NewRedisCache(rc *redis.Client, prefix string) *RedisCache {
	return &RedisCache{
		rc:     rc,
		prefix: prefix,
	}
}

func (this *RedisCache) Get(key string, v interface{}) error {
	if this.Disable {
		return DisableError
	}
	rkey := this.RealKey(key)
	val, err := this.rc.Get(rkey).Result()
	if this.Debug {
		log.Println("[redisCache]", "Get", key, val, err)
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), v)
	return err
}

func (this *RedisCache) Set(key string, v interface{}, exp time.Duration) error {
	if this.Disable {
		return nil
	}
	val, err := json.Marshal(v)
	if err != nil {
		return err
	}
	rkey := this.RealKey(key)
	err = this.rc.Set(rkey, string(val), exp*time.Second).Err()
	if this.Debug {
		log.Println("[redisCache]", "Set", key, string(val), exp, err)
	}
	return err
}

func (this *RedisCache) Del(keys ...string) error {
	if this.Disable {
		return nil
	}
	err := this.rc.Del(keys...).Err()
	return err
}

func (this *RedisCache) RealKey(keys ...interface{}) string {
	str := this.prefix + ":" + helper.Md5(fmt.Sprintln(keys...))
	return str
}

func (this *RedisCache) GetGroup(group interface{}) ([]string, error) {
	if this.Disable {
		return nil, DisableError
	}
	rkey := this.prefix + ":group:" + fmt.Sprint(group)
	keys, err := this.rc.SMembers(rkey).Result()
	if this.Debug {
		log.Println("[redisCache]", "GetGroup", rkey, keys, err)
	}
	return keys, err
}

func (this *RedisCache) SetGroup(group interface{}, keys ...interface{}) error {
	if this.Disable {
		return nil
	}
	rkey := this.prefix + ":group:" + fmt.Sprint(group)
	err := this.rc.SAdd(rkey, keys...).Err()
	if this.Debug {
		log.Println("[redisCache]", "SetGroup", rkey, keys, err)
	}
	return err
}

func (this *RedisCache) Clear(group interface{}) error {
	if this.Disable {
		return nil
	}
	keys, err := this.GetGroup(group)
	if err != nil {
		return nil
	}
	err = this.Del(keys...)
	if err != nil {
		return err
	}
	rkey := this.prefix + ":group:" + fmt.Sprint(group)
	err = this.rc.Del(rkey).Err()
	return err
}
