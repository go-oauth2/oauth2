package oauth2

import (
	"bytes"
	"strconv"

	"github.com/LyricTian/go.uuid"

	"gopkg.in/LyricTian/lib.v2"
)

// NewTokenBasicInfo 创建用于生成令牌的基础信息
// cli 客户端信息
// userID 用户标识
// createAt 创建令牌的时间戳
func NewTokenBasicInfo(cli Client, tokenID, userID string, createAt int64) *TokenBasicInfo {
	return &TokenBasicInfo{
		Client:   cli,
		UserID:   userID,
		TokenID:  tokenID,
		CreateAt: createAt,
	}
}

// TokenBasicInfo 用于生成令牌的基础信息
type TokenBasicInfo struct {
	Client   Client // 客户端信息
	UserID   string // 用户标识
	TokenID  string // 令牌标识
	CreateAt int64  // 创建令牌的时间戳
}

// TokenGenerate 令牌生成接口
type TokenGenerate interface {
	// AccessToken 生成访问令牌
	AccessToken(basicInfo *TokenBasicInfo) (string, error)

	// RefreshToken 生成刷新令牌
	RefreshToken(basicInfo *TokenBasicInfo) (string, error)
}

// NewDefaultTokenGenerate 创建默认的访问令牌生成方式
func NewDefaultTokenGenerate() TokenGenerate {
	return &TokenGenerateDefault{}
}

// TokenGenerateDefault 提供访问令牌、更新令牌的默认生成函数
type TokenGenerateDefault struct{}

// AccessToken 生成访问令牌(md5)
// basicInfo 生成访问令牌的基础参数
func (tg *TokenGenerateDefault) AccessToken(basicInfo *TokenBasicInfo) (token string, err error) {
	ns, _ := uuid.FromString(basicInfo.TokenID)
	buf := bytes.NewBuffer(uuid.NewV3(ns, basicInfo.Client.ID()).Bytes())
	if basicInfo.UserID != "" {
		_, _ = buf.WriteString(basicInfo.UserID)
	}
	_, _ = buf.WriteString(strconv.FormatInt(basicInfo.CreateAt, 10))

	return lib.NewEncryption(buf.Bytes()).MD5()
}

// RefreshToken 生成刷新令牌(sha1)
// basicInfo 生成刷新令牌的基础参数
func (tg *TokenGenerateDefault) RefreshToken(basicInfo *TokenBasicInfo) (string, error) {
	ns, _ := uuid.FromString(basicInfo.TokenID)
	buf := bytes.NewBuffer(uuid.NewV5(ns, basicInfo.Client.ID()).Bytes())
	if basicInfo.UserID != "" {
		_, _ = buf.WriteString(basicInfo.UserID)
	}
	_, _ = buf.WriteString(strconv.FormatInt(basicInfo.CreateAt, 10))

	return lib.NewEncryption(buf.Bytes()).Sha1()
}
