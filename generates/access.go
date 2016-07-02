package generates

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/LyricTian/go.uuid"
	"gopkg.in/LyricTian/lib.v2"
	"gopkg.in/oauth2.v2"
)

// NewAccessGenerate 创建访问令牌生成实例
func NewAccessGenerate() *AccessGenerate {
	return &AccessGenerate{}
}

// AccessGenerate 访问令牌生成
type AccessGenerate struct {
}

// Token 生成令牌
func (ag *AccessGenerate) Token(data *oauth2.GenerateBasic, isGenRefresh bool) (access, refresh string, err error) {
	buf := bytes.NewBufferString(data.Client.GetID())
	buf.WriteString(data.UserID)
	buf.WriteString(strconv.FormatInt(data.CreateAt.UnixNano(), 10))
	access, err = lib.NewEncryption(uuid.NewV3(uuid.NewV4(), buf.String()).Bytes()).MD5()
	if err != nil {
		return
	}
	access = strings.ToUpper(access)
	if isGenRefresh {
		refresh, err = lib.NewEncryption(uuid.NewV5(uuid.NewV4(), buf.String()).Bytes()).Sha1()
		if err != nil {
			return
		}
		refresh = strings.ToUpper(refresh)
	}

	return
}
