package oauth2

import (
	"strconv"

	"gopkg.in/LyricTian/lib.v2"

	"bytes"
)

// NewTokenBasicInfo 创建用于生成令牌的基础信息
// cli 客户端信息
// userID 用户标识
// createAt 创建令牌的时间戳
func NewTokenBasicInfo(cli Client, userID string, createAt int64) TokenBasicInfo {
	return TokenBasicInfo{
		Client:   cli,
		UserID:   userID,
		CreateAt: createAt,
	}
}

// TokenBasicInfo 用于生成令牌的基础信息
type TokenBasicInfo struct {
	Client   Client // 客户端信息
	UserID   string // 用户标识
	CreateAt int64  // 创建令牌的时间戳
}

// TokenGenerate 令牌生成接口
type TokenGenerate interface {
	// AccessToken 根据客户端信息生成访问令牌
	AccessToken(basicInfo TokenBasicInfo) (string, error)

	// RefreshToken 根据客户端信息生成更新令牌
	RefreshToken(basicInfo TokenBasicInfo) (string, error)
}

// NewDefaultTokenGenerate 创建默认的访问令牌生成方式
func NewDefaultTokenGenerate() TokenGenerate {
	return &TokenGenerateDefault{}
}

// TokenGenerateDefault 提供默认的令牌生成
// 采用MD5(ClientID+ClientSecret+RandomCode+Nanosecond Timestamp)的生成方式
type TokenGenerateDefault struct{}

func (tg *TokenGenerateDefault) generate(basicInfo TokenBasicInfo) (string, error) {
	var buf bytes.Buffer
	_, _ = buf.WriteString(basicInfo.Client.ID())
	if basicInfo.UserID != "" {
		_ = buf.WriteByte('_')
		_, _ = buf.WriteString(basicInfo.UserID)
	}
	_ = buf.WriteByte('\n')
	_, _ = buf.WriteString(basicInfo.Client.Secret())
	_ = buf.WriteByte('\n')
	_, _ = buf.WriteString(lib.NewRandom(6).NumberAndLetter())
	_ = buf.WriteByte('\n')
	_, _ = buf.WriteString(strconv.FormatInt(basicInfo.CreateAt, 10))
	val, err := lib.NewEncryption(buf.Bytes()).MD5()
	buf.Reset()
	return val, err
}

// AccessToken Generate access token
func (tg *TokenGenerateDefault) AccessToken(basicInfo TokenBasicInfo) (string, error) {
	return tg.generate(basicInfo)
}

// RefreshToken Generate refresh token
func (tg *TokenGenerateDefault) RefreshToken(basicInfo TokenBasicInfo) (string, error) {
	return tg.generate(basicInfo)
}
