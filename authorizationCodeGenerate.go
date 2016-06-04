package oauth2

import (
	"bytes"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"github.com/LyricTian/go.uuid"

	"gopkg.in/LyricTian/lib.v2"
)

// ACGenerate 授权码生成接口(Authorization Code Generate)
type ACGenerate interface {
	// Code 根据授权码相关信息生成授权码
	Code(info *ACInfo) (string, error)

	// Parse 解析授权码，返回授权信息ID
	Parse(code string) (int64, error)

	// Verify 验证授权码的有效性
	Verify(code string, info *ACInfo) (bool, error)
}

// NewDefaultACGenerate 创建默认的授权码生成方式
func NewDefaultACGenerate() ACGenerate {
	return &ACGenerateDefault{}
}

// ACGenerateDefault 默认的授权码生成方式
type ACGenerateDefault struct{}

func (ag *ACGenerateDefault) genCode(info *ACInfo) (string, error) {
	ns, _ := uuid.FromString(info.Code)
	buf := bytes.NewBuffer(uuid.NewV3(ns, info.ClientID).Bytes())
	_, _ = buf.WriteString(info.UserID)
	_, _ = buf.WriteString(strconv.FormatInt(info.CreateAt, 10))

	md5Val, err := lib.NewEncryption(buf.Bytes()).MD5()
	if err != nil {
		return "", err
	}
	md5Val = md5Val[:15]

	return md5Val, nil
}

// Code Authorization code
func (ag *ACGenerateDefault) Code(info *ACInfo) (string, error) {
	codeVal, err := ag.genCode(info)
	if err != nil {
		return "", err
	}
	val := base64.URLEncoding.EncodeToString([]byte(codeVal + "." + strconv.FormatInt(info.ID, 10)))
	return strings.TrimRight(val, "="), nil
}

func (ag *ACGenerateDefault) parse(code string) (id int64, token string, err error) {
	codeLen := len(code) % 4
	if codeLen > 0 {
		codeLen = 4 - codeLen
	}
	code = code + strings.Repeat("=", codeLen)
	codeBV, err := base64.URLEncoding.DecodeString(code)
	if err != nil {
		return
	}
	codeVal := strings.SplitN(string(codeBV), ".", 2)
	if len(codeVal) != 2 {
		err = errors.New("Token is invalid")
		return
	}
	id, err = strconv.ParseInt(codeVal[1], 10, 64)
	if err != nil {
		return
	}
	token = codeVal[0]
	return
}

// Parse Parse authorization code
func (ag *ACGenerateDefault) Parse(code string) (id int64, err error) {
	id, _, err = ag.parse(code)
	return
}

// Verify Verify code
func (ag *ACGenerateDefault) Verify(code string, info *ACInfo) (valid bool, err error) {
	_, token, err := ag.parse(code)
	if err != nil {
		return
	}
	codeVal, err := ag.genCode(info)
	if err != nil {
		return
	}
	return token == codeVal, nil
}
