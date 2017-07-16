package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"encoding/base64"

	"github.com/influx6/faux/auth"
	"github.com/influx6/faux/httputil"
)

// IdentityStatus defines the status int type which specifies the current state of a giving identity.
type IdentityStatus int

// Defines series of IdentityStatus types.
const (
	Pending IdentityStatus = iota + 1
	Resolved
	Expired
	Revoked
)

// IdentityPath represent the response delivered when a giving oauth relay is called for
// a new user login.
type IdentityPath struct {
	Identity string `json:"identity"`
	Login    string `json:"login"`
}

// IdentityResponse defines a type that contains the initial response received to process a
// authentication/authorization login.
type IdentityResponse struct {
	Code     string                 `json:"code"`
	Identity string                 `json:"identity"`
	Data     map[string]interface{} `json:"data"`
}

// Identity defines the response delivered to all request for the retrieval of a giving identity token
// details.
type Identity struct {
	Identity string                 `json:"identity"`
	Token    auth.Token             `json:"token"`
	Status   IdentityStatus         `json:"status"`
	Data     map[string]interface{} `json:"data"`
}

// OAuthService defines a API which exposes the consistent operations needed to both
// manage and deploy a oauth service, which will manage both OAuth authentication and
// retireve authorization from such service.
type OAuthService interface {
	Revoke(identity string) error
	Get(identity string) (Identity, error)
	New(identity string, secret string) (string, error)
	Process(identity string, response IdentityResponse) error
	Authenticate(identity string, bearerType string, token string) error
}

// AuthAPI defines a core which exposes
type AuthAPI struct {
	Service     OAuthService
	ServiceName string
}

// Approve defines a function for approving a access token received
// from a request.
func (au AuthAPI) Approve(c *httputil.Context) error {
	// identity, ok := c.GetString("identity")
	// if !ok {
	// 	// c.NoContent(http.StatusBadRequest)
	// 	return errors.New("identity param not found")
	// }

	if stateError, ok := c.GetString("error"); ok {
		return fmt.Errorf("Error occured from OAUTH service: %q", stateError)
	}

	secret, ok := c.GetString("state")
	if !ok {
		return errors.New("State value not received")
	}

	code, ok := c.GetString("code")
	if !ok {
		return errors.New("Code value not received")
	}

	// We issue a state secret where its a combination in this format: SERVICENAME:REQUEST-URI:IDENTITY.
	// If the secret does not match this format then it didnt come from us.
	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return err
	}

	sections := strings.Split(string(decodedSecret), ":")
	if len(sections) != 3 {
		return errors.New("Failed to decode request secret")
	}

	identity := sections[2]
	requestURI := sections[1]
	serviceName := sections[0]

	var response IdentityResponse
	response.Code = code
	response.Identity = identity

	if err := au.Service.Process(identity, response); err != nil {
		return err
	}

	_ = requestURI
	_ = serviceName

	return nil
}

// Register defines a function to create a new oauth request for the underline
// OAuthService.
func (au AuthAPI) Register(c *httputil.Context) error {
	url := c.Request().URL
	identity, ok := c.GetString("identity")
	if !ok {
		// c.NoContent(http.StatusBadRequest)
		return errors.New("identity param not found")
	}

	inibase := fmt.Sprintf("%s:%s:%s:%s", au.ServiceName, url.RequestURI(), randString(15), identity)
	secret := base64.StdEncoding.EncodeToString([]byte(inibase))

	redirectURL, err := au.Service.New(identity, secret)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, IdentityPath{
		Identity: identity,
		Login:    redirectURL,
	})
}

// Revoke defines a function to revoke a existing oauth access for the underline
// OAuthService.
func (au AuthAPI) Revoke(c *httputil.Context) error {
	identity, ok := c.GetString("identity")
	if !ok {
		// c.NoContent(http.StatusBadRequest)
		return errors.New("identity param not found")
	}

	if err := au.Service.Revoke(identity); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// Retreive defines a function to return a existing oauth access record through the underline
// OAuthService.
func (au AuthAPI) Retrieve(c *httputil.Context) error {
	identity, ok := c.GetString("identity")
	if !ok {
		// c.NoContent(http.StatusBadRequest)
		return errors.New("identity param not found")
	}

	identityInfo, err := au.Service.Get(identity)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, identityInfo)
}

// Authenticate defines a function to validate a received token against a
// existing oauth access record through the underline OAuthService.
func (au AuthAPI) Authenticate(c *httputil.Context) error {
	identity, ok := c.GetString("identity")
	if !ok {
		// c.NoContent(http.StatusBadRequest)
		return errors.New("identity param not found")
	}

	authorization := c.Header().Get("Authorization")

	bearer, token, err := auth.ParseAuthorization(authorization)
	if err != nil {
		return err
	}

	// if we are encoded, then it means token is base64 encode and we need to split it for the real token.
	if encodedAuthorization := c.Header().Get("Base64Authorization"); encodedAuthorization != "" {
		userIdentity, userToken, err := auth.ParseToken(token)
		if err != nil {
			return err
		}

		// if userIdentity does not match given identity then fail this.
		if userIdentity != identity {
			return errors.New("Identity does not match token identity")
		}

		if err := au.Service.Authenticate(identity, bearer, userToken); err != nil {
			return err
		}

		return nil
	}

	if err := au.Service.Authenticate(identity, bearer, token); err != nil {
		return err
	}

	return nil
}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
