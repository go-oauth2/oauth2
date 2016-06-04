package oauth2

// MongoConfig MongoDB配置参数
type MongoConfig struct {
	URL    string // MongoDB连接字符串
	DBName string // 数据库名称
}

// NewMongoConfig 创建MongoDB配置参数的实例
func NewMongoConfig(url, dbName string) *MongoConfig {
	return &MongoConfig{
		URL:    url,
		DBName: dbName,
	}
}

// ACConfig 授权码模式配置参数(Authorization Code Config)
type ACConfig struct {
	ACExpiresIn int64 // 授权码有效期(单位秒)
	ATExpiresIn int64 // 访问令牌有效期(单位秒)
	RTExpiresIn int64 // 更新令牌有效期(单位秒)
}

// ImplicitConfig 简化模式配置参数
type ImplicitConfig struct {
	ATExpiresIn int64 // 访问令牌有效期(单位秒)
}

// PasswordConfig 密码模式配置参数
type PasswordConfig struct {
	ATExpiresIn int64 // 访问令牌有效期(单位秒)
	RTExpiresIn int64 // 更新令牌有效期(单位秒)
}

// CCConfig 客户端模式配置参数(Client Credentials Config)
type CCConfig struct {
	ATExpiresIn int64 // 访问令牌有效期(单位秒)
}

// OAuthConfig OAuth授权配置参数
type OAuthConfig struct {
	ACConfig       *ACConfig       // 授权码模式配置参数
	ImplicitConfig *ImplicitConfig // 简化模式配置参数
	PasswordConfig *PasswordConfig // 密码模式配置参数
	CCConfig       *CCConfig       // 客户端模式配置参数
}
