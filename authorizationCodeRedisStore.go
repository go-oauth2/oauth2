package oauth2

import (
	"encoding/json"
	"fmt"

	"gopkg.in/redis.v3"
)

const (
	// DefaultACRedisIDKey Redis存储授权码唯一标识的键
	DefaultACRedisIDKey = "ACID"
)

// NewACRedisStore 创建Redis存储的实例
// config Redis配置参数
// key Redis存储授权码唯一标识的键(默认为ACID)
func NewACRedisStore(cfg *RedisConfig, key string) (*ACRedisStore, error) {
	opt := &redis.Options{
		Network:      cfg.Network,
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  cfg.PoolTimeout,
	}
	cli := redis.NewClient(opt)
	err := cli.Ping().Err()
	if err != nil {
		return nil, err
	}
	if key == "" {
		key = DefaultACRedisIDKey
	}
	return &ACRedisStore{
		cli: cli,
		key: key,
	}, nil
}

// ACRedisStore 提供授权码的redis存储
type ACRedisStore struct {
	cli *redis.Client
	key string
}

// Put 存储授权码
func (ar *ACRedisStore) Put(item ACInfo) (id int64, err error) {
	n, err := ar.cli.Incr(ar.key).Result()
	if err != nil {
		return
	}
	item.ID = n
	jv, err := json.Marshal(item)
	if err != nil {
		return
	}
	key := fmt.Sprintf("%s_%d", ar.key, n)
	err = ar.cli.Set(key, string(jv), item.ExpiresIn).Err()
	if err != nil {
		return
	}
	id = item.ID
	return
}

// TakeByID 取出授权码
func (ar *ACRedisStore) TakeByID(id int64) (info *ACInfo, err error) {
	key := fmt.Sprintf("%s_%d", ar.key, id)
	data, err := ar.cli.Get(key).Result()
	if err != nil {
		return
	}
	var v ACInfo
	err = json.Unmarshal([]byte(data), &v)
	if err != nil {
		return
	}
	err = ar.cli.Del(key).Err()
	if err != nil {
		return
	}
	info = &v
	return
}
