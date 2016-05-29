package oauth2

import (
	"gopkg.in/LyricTian/lib.v2/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// DefaultTokenCollectionName 默认的令牌存储集合名称
	DefaultTokenCollectionName = "AuthToken"
)

// NewTokenMongoStore 创建基于MongoDB的令牌存储方式
// mongoConfig MongoDB配置参数
// cName 存储令牌的集合名称(默认为AuthToken)
func NewTokenMongoStore(mongoConfig *MongoConfig, cName string) (TokenStore, error) {
	mHandler, err := mongo.InitHandlerWithDB(mongoConfig.URL, mongoConfig.DBName)
	if err != nil {
		return nil, err
	}
	if cName == "" {
		cName = DefaultTokenCollectionName
	}
	err = mHandler.C(cName).EnsureIndexKey("AccessToken", "RefreshToken")
	if err != nil {
		return nil, err
	}
	return &TokenMongoStore{
		cName:    cName,
		mHandler: mHandler,
	}, nil
}

// TokenMongoStore 基于MongoDB的令牌存储方式
type TokenMongoStore struct {
	cName    string
	mHandler *mongo.Handler
}

// Create Add item
func (tm *TokenMongoStore) Create(item *Token) (id int64, err error) {
	tm.mHandler.CHandle(tm.cName, func(c *mgo.Collection) {
		tid, err := tm.mHandler.IncrID(tm.cName)
		if err != nil {
			return
		}
		item.ID = tid
		err = c.Insert(item)
		if err != nil {
			return
		}
		id = item.ID
	})
	return
}

// Update Modify item
func (tm *TokenMongoStore) Update(id int64, info map[string]interface{}) (err error) {
	tm.mHandler.CHandle(tm.cName, func(c *mgo.Collection) {
		err = c.UpdateId(id, bson.M{"$set": info})
		if err != nil {
			return
		}
	})
	return
}

func (tm *TokenMongoStore) findOne(query interface{}) (token *Token, err error) {
	tm.mHandler.CHandle(tm.cName, func(c *mgo.Collection) {
		var result []Token
		err = c.Find(query).Sort("-_id").Limit(1).All(&result)
		if err != nil {
			return
		}
		if len(result) > 0 {
			token = &result[0]
		}
	})
	return
}

// GetByAccessToken 根据访问令牌获取令牌信息
func (tm *TokenMongoStore) GetByAccessToken(accessToken string) (*Token, error) {
	return tm.findOne(bson.M{"AccessToken": accessToken})
}

// GetByRefreshToken 根据更新令牌获取令牌信息
func (tm *TokenMongoStore) GetByRefreshToken(refreshToken string) (*Token, error) {
	return tm.findOne(bson.M{"RefreshToken": refreshToken})
}
