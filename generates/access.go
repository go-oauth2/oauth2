package generates

import (
	"bytes"
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/satori/go.uuid"
	"gopkg.in/oauth2.v3"
)

// NewAccessGenerate create to generate the access token instance
func NewAccessGenerate() *AccessGenerate {
	return &AccessGenerate{}
}

// AccessGenerate generate the access token
type AccessGenerate struct {
}

// Token based on the UUID generated token
func (ag *AccessGenerate) Token(data *oauth2.GenerateBasic, isGenRefresh bool) (access, refresh string, err error) {
	buf := bytes.NewBufferString(data.Client.GetID())
	buf.WriteString(data.UserID)
	buf.WriteString(strconv.FormatInt(data.CreateAt.UnixNano(), 10))

	access = base64.URLEncoding.EncodeToString(uuid.NewV3(uuid.NewV4(), buf.String()).Bytes())
	access = strings.ToUpper(strings.TrimRight(access, "="))
	if isGenRefresh {
		refresh = base64.URLEncoding.EncodeToString(uuid.NewV5(uuid.NewV4(), buf.String()).Bytes())
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return
}
