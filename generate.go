package oauth2

import "time"

type (
	// GenerateBasic Provide the basis of the generated token data
	GenerateBasic struct {
		Client   ClientInfo // The client information
		UserID   string     // The user id
		CreateAt time.Time  // Creation time
	}

	// AuthorizeGenerate Generate the authorization code interface
	AuthorizeGenerate interface {
		Token(data *GenerateBasic) (code string, err error)
	}

	// AccessGenerate Generate the access and refresh tokens interface
	AccessGenerate interface {
		Token(data *GenerateBasic, isGenRefresh bool) (access, refresh string, err error)
	}
)
