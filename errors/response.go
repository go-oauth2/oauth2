package errors

import (
	"errors"
	"net/http"
)

// Response error response
// 错误响应
type Response struct {
	Error       error
	ErrorCode   int
	Description string
	URI         string
	StatusCode  int
	Header      http.Header
}

// NewResponse create the response pointer
// 创建响应指针
func NewResponse(err error, statusCode int) *Response {
	return &Response{
		Error:      err,
		StatusCode: statusCode,
	}
}

// SetHeader sets the header entries associated with key to
// the single element value.
// 将与key关联的标题条目设置为
// 单个元素的值。
func (r *Response) SetHeader(key, value string) {
	if r.Header == nil {
		r.Header = make(http.Header)
	}
	r.Header.Set(key, value)
}

// https://tools.ietf.org/html/rfc6749#section-5.2
var (
	ErrInvalidRequest          = errors.New("invalid_request")
	ErrUnauthorizedClient      = errors.New("unauthorized_client")
	ErrAccessDenied            = errors.New("access_denied")
	ErrUnsupportedResponseType = errors.New("unsupported_response_type")
	ErrInvalidScope            = errors.New("invalid_scope")
	ErrServerError             = errors.New("server_error")
	ErrTemporarilyUnavailable  = errors.New("temporarily_unavailable")
	ErrInvalidClient           = errors.New("invalid_client")
	ErrInvalidGrant            = errors.New("invalid_grant")
	ErrUnsupportedGrantType    = errors.New("unsupported_grant_type")
)

// Descriptions error description
// 错误说明
var Descriptions = map[error]string{
	ErrInvalidRequest:          "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed",
	ErrUnauthorizedClient:      "The client is not authorized to request an authorization code using this method",
	ErrAccessDenied:            "The resource owner or authorization server denied the request",
	ErrUnsupportedResponseType: "The authorization server does not support obtaining an authorization code using this method",
	ErrInvalidScope:            "The requested scope is invalid, unknown, or malformed",
	ErrServerError:             "The authorization server encountered an unexpected condition that prevented it from fulfilling the request",
	ErrTemporarilyUnavailable:  "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server",
	ErrInvalidClient:           "Client authentication failed",
	ErrInvalidGrant:            "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client",
	ErrUnsupportedGrantType:    "The authorization grant type is not supported by the authorization server",
}

// StatusCodes response error HTTP status code
// 响应错误HTTP状态代码
var StatusCodes = map[error]int{
	ErrInvalidRequest:          400,
	ErrUnauthorizedClient:      401,
	ErrAccessDenied:            403,
	ErrUnsupportedResponseType: 401,
	ErrInvalidScope:            400,
	ErrServerError:             500,
	ErrTemporarilyUnavailable:  503,
	ErrInvalidClient:           401,
	ErrInvalidGrant:            401,
	ErrUnsupportedGrantType:    401,
}
