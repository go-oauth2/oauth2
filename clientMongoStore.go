package oauth2

import (
	"gopkg.in/LyricTian/lib.v2/mongo"
	"gopkg.in/mgo.v2"
)

const (
	// DefaultClientCollectionName 默认的客户端存储集合名称
	DefaultClientCollectionName = "ClientInfo"
)

// NewClientMongoStore 创建基于MongoDB的客户端存储方式
// mongoConfig MongoDB配置参数
// cName 存储客户端的集合名称(默认为ClientInfo)
func NewClientMongoStore(mongoConfig *MongoConfig, cName string) (ClientStore, error) {
	mHandler, err := mongo.InitHandlerWithDB(mongoConfig.URL, mongoConfig.DBName)
	if err != nil {
		return nil, err
	}
	if cName == "" {
		cName = DefaultClientCollectionName
	}
	return &ClientMongoStore{
		cName:    cName,
		mHandler: mHandler,
	}, nil
}

// ClientMongoStore 基于MongoDB的默认客户端信息存储
type ClientMongoStore struct {
	cName    string
	mHandler *mongo.Handler
}

// GetByID 根据ID获取客户端信息
func (dcm *ClientMongoStore) GetByID(id string) (client Client, err error) {
	dcm.mHandler.CHandle(dcm.cName, func(c *mgo.Collection) {
		var result []DefaultClient
		err = dcm.mHandler.C(dcm.cName).FindId(id).Limit(1).All(&result)
		if err != nil {
			return
		}
		if len(result) > 0 {
			client = result[0]
		}
	})
	return
}
