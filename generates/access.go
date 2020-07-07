package generates

import (
	"bytes"
	"context"
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/google/uuid"
)

// NewAccessGenerate create to generate the access token instance
// 创建以生成访问令牌实例
func NewAccessGenerate() *AccessGenerate {
	return &AccessGenerate{}
}

// AccessGenerate generate the access token
// 生成访问令牌
type AccessGenerate struct {
}

// Token based on the UUID generated token
// 基于UUID生成的令牌的令牌
func (ag *AccessGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	buf := bytes.NewBufferString(data.Client.GetID())
	buf.WriteString(data.UserID)
	buf.WriteString(strconv.FormatInt(data.CreateAt.UnixNano(), 10))

	access := base64.URLEncoding.EncodeToString([]byte(uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes()).String()))
	access = strings.ToUpper(strings.TrimRight(access, "="))
	refresh := ""
	if isGenRefresh {
		refresh = base64.URLEncoding.EncodeToString([]byte(uuid.NewSHA1(uuid.Must(uuid.NewRandom()), buf.Bytes()).String()))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return access, refresh, nil
}
