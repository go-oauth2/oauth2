package token

import (
	"encoding/json"

	"gopkg.in/oauth2.v2"
	"gopkg.in/oauth2.v2/models"
	"gopkg.in/redis.v4"
)

// NewRedisStore 创建redis存储的实例
func NewRedisStore(cfg *RedisConfig) (store oauth2.TokenStore, err error) {
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
	if verr := cli.Ping().Err(); verr != nil {
		err = verr
		return
	}
	store = &RedisStore{cli: cli}
	return
}

// RedisStore 令牌的redis存储
type RedisStore struct {
	cli *redis.Client
}

// Create 存储令牌信息
func (rs *RedisStore) Create(info oauth2.TokenInfo) (err error) {
	jv, err := json.Marshal(info)
	if err != nil {
		return
	}
	pipe := rs.cli.Pipeline()

	aexp := info.GetAccessExpiresIn()
	if refresh := info.GetRefresh(); refresh != "" {
		exp := info.GetRefreshExpiresIn()
		ttl := rs.cli.TTL(refresh)
		if verr := ttl.Err(); verr != nil {
			err = verr
			return
		}
		if v := ttl.Val(); v.Seconds() > 0 {
			exp = v
		}
		if aexp.Seconds() > exp.Seconds() {
			aexp = exp
		}
		pipe.Set(refresh, jv, exp)
	}
	pipe.Set(info.GetAccess(), jv, aexp)

	if _, verr := pipe.Exec(); verr != nil {
		err = verr
	}
	return
}

// remove
func (rs *RedisStore) remove(key string) (err error) {
	del := rs.cli.Del(key)
	if verr := del.Err(); verr != nil {
		err = verr
	}
	return
}

// RemoveByAccess 移除令牌
func (rs *RedisStore) RemoveByAccess(access string) (err error) {
	err = rs.remove(access)
	return
}

// RemoveByRefresh 移除令牌
func (rs *RedisStore) RemoveByRefresh(refresh string) (err error) {
	err = rs.remove(refresh)
	return
}

func (rs *RedisStore) get(key string) (ti oauth2.TokenInfo, err error) {
	gv, gerr := rs.cli.Get(key).Result()
	if gerr != nil {
		if gerr == redis.Nil {
			return
		}
		err = gerr
		return
	}
	var tm models.Token
	if verr := json.Unmarshal([]byte(gv), &tm); verr != nil {
		err = verr
		return
	}
	ti = &tm
	return
}

// GetByAccess 获取令牌数据
func (rs *RedisStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	ti, err = rs.get(access)
	return
}

// GetByRefresh 获取令牌数据
func (rs *RedisStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	ti, err = rs.get(refresh)
	return
}
