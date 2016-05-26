package oauth2_test

import (
	"gopkg.in/LyricTian/lib.v2"
	"gopkg.in/LyricTian/lib.v2/mongo"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/oauth2.v1"
)

const (
	// MongoURL MongoDB连接字符串
	MongoURL = "mongodb://admin:123456@45.78.35.157:37017"
	// DBName 数据库名称
	DBName = "test"
)

// ClientHandle 执行客户端处理
func ClientHandle(handle func(cli oauth2.Client)) {
	info := oauth2.DefaultClient{
		ClientID:     bson.NewObjectId().Hex(),
		ClientDomain: "http://www.example.com",
	}
	info.ClientSecret, _ = lib.NewEncryption([]byte(info.ClientID)).MD5()
	mHandler, err := mongo.InitHandlerWithDB(MongoURL, DBName)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = mHandler.C(oauth2.DefaultClientCollectionName).RemoveId(info.ClientID)
		if err != nil {
			panic(err)
		}
		mHandler.Session().Close()
	}()
	err = mHandler.C(oauth2.DefaultClientCollectionName).Insert(info)
	if err != nil {
		panic(err)
	}
	handle(info)
}
