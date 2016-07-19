package token

import (
	"encoding/json"

	"gopkg.in/redis.v4"

	"strconv"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

// DefaultIncrKey store incr id
const DefaultIncrKey = "oauth2_incr"

// NewRedisStore Create a token store instance based on redis
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

// RedisStore Redis Store
type RedisStore struct {
	cli *redis.Client
}

func (rs *RedisStore) getBasicID(id int64, info oauth2.TokenInfo) string {
	return "oauth2_" + info.GetClientID() + "_" + strconv.FormatInt(id, 10)
}

// Create Create and store the new token information
func (rs *RedisStore) Create(info oauth2.TokenInfo) (err error) {
	jv, err := json.Marshal(info)
	if err != nil {
		return
	}
	id, err := rs.cli.Incr(DefaultIncrKey).Result()
	if err != nil {
		return
	}
	pipe := rs.cli.Pipeline()
	basicID := rs.getBasicID(id, info)
	aexp := info.GetAccessExpiresIn()
	rexp := aexp

	if refresh := info.GetRefresh(); refresh != "" {
		rexp = info.GetRefreshExpiresIn()
		ttl := rs.cli.TTL(refresh)
		if verr := ttl.Err(); verr != nil {
			err = verr
			return
		}
		if v := ttl.Val(); v.Seconds() > 0 {
			rexp = v
		}
		if aexp.Seconds() > rexp.Seconds() {
			aexp = rexp
		}
		pipe.Set(refresh, basicID, rexp)
	}
	pipe.Set(info.GetAccess(), basicID, aexp)
	pipe.Set(basicID, jv, rexp)

	if _, verr := pipe.Exec(); verr != nil {
		err = verr
	}
	return
}

// remove
func (rs *RedisStore) remove(key string) (err error) {
	_, verr := rs.cli.Del(key).Result()
	if verr != redis.Nil {
		err = verr
	}
	return
}

// RemoveByAccess Use the access token to delete the token information(Along with the refresh token)
func (rs *RedisStore) RemoveByAccess(access string) (err error) {
	err = rs.remove(access)
	return
}

// RemoveByRefresh Use the refresh token to delete the token information
func (rs *RedisStore) RemoveByRefresh(refresh string) (err error) {
	err = rs.remove(refresh)
	return
}

// get
func (rs *RedisStore) get(token string) (ti oauth2.TokenInfo, err error) {
	tv, verr := rs.cli.Get(token).Result()
	if verr != nil {
		if verr == redis.Nil {
			return
		}
		err = verr
		return
	}
	result := rs.cli.Get(tv)
	if verr := result.Err(); verr != nil {
		if verr == redis.Nil {
			return
		}
		err = verr
		return
	}
	iv, err := result.Bytes()
	if err != nil {
		return
	}
	var tm models.Token
	if verr := json.Unmarshal(iv, &tm); verr != nil {
		err = verr
		return
	}
	ti = &tm
	return
}

// GetByAccess Use the access token for token information data
func (rs *RedisStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	ti, err = rs.get(access)
	return
}

// GetByRefresh Use the refresh token for token information data
func (rs *RedisStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	ti, err = rs.get(refresh)
	return
}
