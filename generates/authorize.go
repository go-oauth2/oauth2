package generates

import (
	"bytes"
	"strings"

	"github.com/LyricTian/go.uuid"
	"gopkg.in/LyricTian/lib.v2"
	"gopkg.in/oauth2.v2"
)

// NewAuthorizeGenerate 创建授权令牌生成实例
func NewAuthorizeGenerate() *AuthorizeGenerate {
	return &AuthorizeGenerate{}
}

// AuthorizeGenerate 授权令牌生成
type AuthorizeGenerate struct{}

// Token 生成令牌
func (ag *AuthorizeGenerate) Token(data *oauth2.GenerateBasic) (code string, err error) {
	buf := bytes.NewBuffer(uuid.NewV1().Bytes())
	buf.WriteString(data.UserID)
	buf.WriteString(data.Client.GetID())
	code, err = lib.NewEncryption(buf.Bytes()).MD5()
	if err != nil {
		return
	}
	code = strings.ToUpper(code)
	return
}
