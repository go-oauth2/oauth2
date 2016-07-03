package token

import (
	"time"

	"gopkg.in/LyricTian/lib.v2/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/oauth2.v2"
	"gopkg.in/oauth2.v2/models"
)

// MongoConfig MongoDB Configuration
type MongoConfig struct {
	// Connection String
	URL string
	// DB Name(default oauth2)
	DB string
	// Collection Name(default tokens)
	C string
}

// NewMongoStore 创建MongoDB的令牌存储
func NewMongoStore(cfg *MongoConfig) (store oauth2.TokenStore, err error) {
	if cfg.DB == "" {
		cfg.DB = "oauth2"
	}
	if cfg.C == "" {
		cfg.C = "tokens"
	}
	handler, err := mongo.InitHandlerWithDB(cfg.URL, cfg.DB)
	if err != nil {
		return
	}
	// 创建自动过期索引
	err = handler.C(cfg.C).EnsureIndex(mgo.Index{
		Key:         []string{"ExpiredAt"},
		ExpireAfter: time.Second,
	})
	if err != nil {
		return
	}
	err = handler.C(cfg.C).EnsureIndexKey("Access")
	if err != nil {
		return
	}
	err = handler.C(cfg.C).EnsureIndexKey("Refresh")
	if err != nil {
		return
	}
	store = &MongoStore{
		handler: handler,
		cfg:     cfg,
	}
	return
}

// MongoStore MongoDB Store
type MongoStore struct {
	cfg     *MongoConfig
	handler *mongo.Handler
}

// Create 存储令牌信息
func (ms *MongoStore) Create(info oauth2.TokenInfo) (err error) {
	tm := info.(*models.Token)
	var expiredAt time.Time
	if refresh := tm.Refresh; refresh != "" {
		expiredAt = tm.RefreshCreateAt.Add(tm.RefreshExpiresIn)
		rinfo, rerr := ms.GetByRefresh(refresh)
		if rerr != nil {
			err = rerr
			return
		}
		if rinfo != nil {
			expiredAt = rinfo.GetRefreshCreateAt().Add(rinfo.GetRefreshExpiresIn())
		}
	}
	if expiredAt.IsZero() {
		expiredAt = tm.AccessCreateAt.Add(tm.AccessExpiresIn)
	}
	doc := map[string]interface{}{
		"ExpiredAt":        expiredAt,
		"ClientID":         tm.ClientID,
		"UserID":           tm.UserID,
		"RedirectURI":      tm.RedirectURI,
		"Scope":            tm.Scope,
		"AuthType":         tm.AuthType,
		"Access":           tm.Access,
		"AccessCreateAt":   tm.AccessCreateAt,
		"AccessExpiresIn":  tm.AccessExpiresIn,
		"Refresh":          tm.Refresh,
		"RefreshCreateAt":  tm.RefreshCreateAt,
		"RefreshExpiresIn": tm.RefreshExpiresIn,
	}

	ms.handler.CHandle(ms.cfg.C, func(c *mgo.Collection) {
		err = c.Insert(doc)
	})
	return
}

func (ms *MongoStore) remove(selector interface{}) (err error) {
	ms.handler.CHandle(ms.cfg.C, func(c *mgo.Collection) {
		err = c.Remove(selector)
	})
	return
}

// RemoveByAccess 移除令牌
func (ms *MongoStore) RemoveByAccess(access string) (err error) {
	err = ms.remove(bson.M{"Access": access})
	return
}

// RemoveByRefresh 移除令牌
func (ms *MongoStore) RemoveByRefresh(refresh string) (err error) {
	err = ms.remove(bson.M{"Refresh": refresh})
	return
}

func (ms *MongoStore) get(find interface{}) (info oauth2.TokenInfo, err error) {
	ms.handler.CHandle(ms.cfg.C, func(c *mgo.Collection) {
		var tm models.Token
		aerr := c.Find(find).Select(bson.M{"_id": 0}).One(&tm)
		if aerr != nil {
			if aerr == mgo.ErrNotFound {
				return
			}
			err = aerr
			return
		}
		info = &tm
	})
	return
}

// GetByAccess 获取令牌数据
func (ms *MongoStore) GetByAccess(access string) (info oauth2.TokenInfo, err error) {
	info, err = ms.get(bson.M{"Access": access})
	return
}

// GetByRefresh 获取令牌数据
func (ms *MongoStore) GetByRefresh(refresh string) (info oauth2.TokenInfo, err error) {
	info, err = ms.get(bson.M{"Refresh": refresh})
	return
}
