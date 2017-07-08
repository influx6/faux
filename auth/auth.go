package auth

import (
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"strings"
	"errors"
)

// Credential defines a struct which holds clientID and clientSecret which
// are used by oauths.
type Credential struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
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
	Type        string    `json:"type"`
	AccessToken string    `json:"access_token"`
	Expires     time.Time `json:"expires"`
}

// Fields returns the given fields as a map.
func (t Token) Fields() map[string]interface{} {
	return map[string]interface{}{
		"type":         t.Type,
		"access_token": t.AccessToken,
		"expires":      t.Expires,
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
		Type:        token.Type(),
		AccessToken: token.AccessToken,
		Expires:     token.Expiry,
	}, nil
}
