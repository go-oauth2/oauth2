package oauth2

import (
	"gopkg.in/LyricTian/lib.v2"
	"gopkg.in/LyricTian/lib.v2/mongo"
	"gopkg.in/mgo.v2/bson"
)

const (
	// MongoURL MongoDB连接字符串
	MongoURL = "mongodb://admin:123456@192.168.33.70:27017"
	// DBName 数据库名称
	DBName = "test"
)

var (
	oManager *OAuthManager
)

// ClientHandle 执行客户端处理
func ClientHandle(handle func(cli Client)) {
	info := DefaultClient{
		ClientID:     bson.NewObjectId().Hex(),
		ClientDomain: "http://www.example.com",
	}
	info.ClientSecret, _ = lib.NewEncryption([]byte(info.ClientID)).MD5()
	mHandler, err := mongo.InitHandlerWithDB(MongoURL, DBName)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = mHandler.C(DefaultClientCollectionName).RemoveId(info.ClientID)
		if err != nil {
			panic(err)
		}
		mHandler.Session().Close()
	}()
	err = mHandler.C(DefaultClientCollectionName).Insert(info)
	if err != nil {
		panic(err)
	}
	handle(info)
}
