package generates

import (
	"bytes"
	"context"
	"encoding/base64"
	"strings"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/google/uuid"
)

// NewAuthorizeGenerate create to generate the authorize code instance
// 创建以生成授权代码实例
func NewAuthorizeGenerate() *AuthorizeGenerate {
	return &AuthorizeGenerate{}
}

// AuthorizeGenerate generate the authorize code
// 生成授权代码
type AuthorizeGenerate struct{}

// Token based on the UUID generated token
// 基于UUID生成的令牌的令牌
func (ag *AuthorizeGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic) (string, error) {
	buf := bytes.NewBufferString(data.Client.GetID())
	buf.WriteString(data.UserID)
	token := uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes())
	code := base64.URLEncoding.EncodeToString([]byte(token.String()))
	code = strings.ToUpper(strings.TrimRight(code, "="))

	return code, nil
}
