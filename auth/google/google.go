package google

import (
	"github.com/influx6/faux/auth"
	"golang.org/x/oauth2/google" // https://github.com/golang/oauth2/google
)

// contains sets of user related scopes.
var (
	EmailScope    = "https://www.googleapis.com/auth/userinfo.email"
	UserInfoScope = "https://www.googleapis.com/auth/userinfo.profile"
)

// New returns a new instance of auth.Auth for use with the google OAuth2 API.
func New(cred auth.Credential, redirectURL string) *auth.Auth {
	var hasEmail, hasUserInfo bool

	for _, scope := range cred.Scopes {
		if scope == EmailScope {
			hasEmail = true
		}

		if scope == UserInfoScope {
			hasUserInfo = true
		}
	}

	if !hasEmail {
		cred.Scopes = append(cred.Scopes, EmailScope)
	}

	if !hasUserInfo {
		cred.Scopes = append(cred.Scopes, UserInfoScope)
	}

	return auth.New(cred, google.Endpoint, redirectURL)
}
