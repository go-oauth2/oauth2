package generates

import (
	"bytes"
	"context"
	"encoding/base64"
	"strings"

	"gopkg.in/oauth2.v4"
	"gopkg.in/oauth2.v4/utils/uuid"
)

// NewAuthorizeGenerate create to generate the authorize code instance
func NewAuthorizeGenerate() *AuthorizeGenerate {
	return &AuthorizeGenerate{}
}

// AuthorizeGenerate generate the authorize code
type AuthorizeGenerate struct{}

// Token based on the UUID generated token
func (ag *AuthorizeGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic) (string, error) {
	buf := bytes.NewBufferString(data.Client.GetID())
	buf.WriteString(data.UserID)
	token := uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes())
	code := base64.URLEncoding.EncodeToString(token.Bytes())
	code = strings.ToUpper(strings.TrimRight(code, "="))

	return code, nil
}
