package auth

import (
	"encoding/base64"
	"net/http"
	"time"

	"errors"
	"strings"

	"golang.org/x/oauth2"
)

// Credential defines a struct which holds clientID and clientSecret which
// are used by oauths.
type Credential struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
}

//====================================================================================================

// Auth defines a structure which allows us properly retrieve oauth
// authentication data from the OAuth2 api.
type Auth struct {
	config oauth2.Config
}

// New returns a new instance of OAuth.
func New(cred Credential, endpoints oauth2.Endpoint, redirectURL string) *Auth {
	return &Auth{
		config: oauth2.Config{
			Endpoint:     endpoints,
			RedirectURL:  redirectURL,
			Scopes:       cred.Scopes,
			ClientID:     cred.ClientID,
			ClientSecret: cred.ClientSecret,
		},
	}
}

// LoginURL returns the login URL for redirect users to login into acct
// to provide access token for API requests.
func (a *Auth) LoginURL(state string, xs ...oauth2.AuthCodeOption) string {
	return a.config.AuthCodeURL(state, xs...)
}

// Token defines the data returned from a OAuth op.
type Token struct {
	Type         string    `json:"type"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expires      time.Time `json:"expires"`
}

// Fields returns the given fields as a map.
func (t Token) Fields() map[string]interface{} {
	return map[string]interface{}{
		"type":          t.Type,
		"access_token":  t.AccessToken,
		"refresh_token": t.RefreshToken,
		"expires":       t.Expires,
	}
}

// AuthorizeFromUser takes the code retrieved from the users login process
// and attempts to retrieve a access token from the configuration.
func (a *Auth) AuthorizeFromUser(code string) (*http.Client, Token, error) {
	token, err := a.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, Token{}, err
	}

	return a.config.Client(oauth2.NoContext, token), Token{
		Type:         token.Type(),
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expires:      token.Expiry,
	}, nil
}

//===================================================================================================

// ParseAuthorization returns the scheme and token of the Authorization string
// if it's valid.
func ParseAuthorization(val string) (authType string, token string, err error) {
	authSplit := strings.SplitN(val, " ", 2)
	if len(authSplit) != 2 {
		err = errors.New("Invalid Authorization: Expected content: `AuthType Token`")
		return
	}

	authType = strings.TrimSpace(authSplit[0])
	token = strings.TrimSpace(authSplit[1])

	return
}

// ParseToken parses the base64 encoded token, which it returns the
// associated userID and session token.
func ParseToken(val string) (userID string, token string, err error) {
	var decoded []byte

	decoded, err = base64.StdEncoding.DecodeString(val)
	if err != nil {
		return
	}

	// Attempt to get the session token split which has the userid:session_token.
	sessionToken := strings.Split(string(decoded), ":")
	if len(sessionToken) != 2 {
		err = errors.New("Invalid Token: Token must be UserID:Token  format")
		return
	}

	userID = sessionToken[0]
	token = sessionToken[1]

	return
}
