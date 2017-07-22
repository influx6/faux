package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"encoding/base64"
	"encoding/json"

	"github.com/influx6/faux/auth"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/httputil"
)

//====================================================================================================

// contains series of errors used for Authentication.
var (
	ErrIdentityNotFound   = errors.New("Identity not found")
	ErrIdentityStillValid = errors.New("Identity Credentials Still Valid")
	ErrIdentityHasExpired = errors.New("Identity Credentials Has Expired")
)

//====================================================================================================

// IdentityStatus defines the status int type which specifies the current state of a giving identity.
type IdentityStatus int

// Defines series of IdentityStatus types.
const (
	Pending IdentityStatus = iota + 1
	Resolved
	Expired
)

// IdentityPath represent the response delivered when a giving oauth relay is called for
// a new user login.
type IdentityPath struct {
	Identity string `json:"identity" bson:"identity"`
	Login    string `json:"login" bson:"login"`
}

// IdentityResponse defines a type that contains the initial response received to process a
// authentication/authorization login.
type IdentityResponse struct {
	Code     string                 `json:"code" bson:"code"`
	Identity string                 `json:"identity" bson:"identity"`
	Data     map[string]interface{} `json:"data" bson:"data"`
}

// Identity defines the response delivered to all request for the retrieval of a giving identity token
// details.
type Identity struct {
	Identity  string                 `json:"identity" bson:"identity"`
	TokenID   string                 `json:"token_id" bson:"token_id"`
	PrivateID string                 `json:"private_id" bson:"private_id"`
	Token     auth.Token             `json:"token" bson:"token"`
	Status    IdentityStatus         `json:"status" bson:"status"`
	Data      map[string]interface{} `json:"data" bson:"data"`
}

// OAuthService defines a API which exposes the consistent operations needed to both
// manage and deploy a oauth service, which will manage both OAuth authentication and
// retireve authorization from such service.
type OAuthService interface {
	Revoke(ctx context.Context, identity string) error
	Identities(ctx context.Context) ([]Identity, error)
	Get(ctx context.Context, identity string) (Identity, error)
	Approve(ctx context.Context, response IdentityResponse) error
	New(ctx context.Context, identity string, secret string) (string, error)
	Authenticate(ctx context.Context, identity string, bearerType string, token string) error
}

// AuthAPI defines a core which exposes
type AuthAPI struct {
	Service     OAuthService
	ServiceName string
}

// New returns a new instance of the AuthAPI.
func New(service OAuthService, serviceName string) AuthAPI {
	return AuthAPI{
		Service:     service,
		ServiceName: serviceName,
	}
}

// Approve defines a function for approving a access token received
// from a request.
func (au AuthAPI) Approve(c *httputil.Context) error {
	if stateError, ok := c.GetString("error"); ok {
		return fmt.Errorf("Error occured from OAUTH service: %q", stateError)
	}

	defer c.Request().Body.Close()

	var data map[string]interface{}

	if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return err
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
	serviceName := sections[0]
	encodedRequestURI := sections[1]

	requestURI, err := base64.StdEncoding.DecodeString(encodedRequestURI)
	if err != nil {
		return err
	}

	if serviceName != au.ServiceName {
		return fmt.Errorf("ServiceName mismatch with %q for %q", serviceName, au.ServiceName)
	}

	var response IdentityResponse
	response.Data = data
	response.Code = code
	response.Identity = identity

	if err2 := au.Service.Approve(c, response); err2 != nil {
		return err2
	}

	// if we received approval then retrieve identity and add as cookie.
	identityData, err := au.Service.Get(c, identity)
	if err != nil {
		return err
	}

	// If expired then reovke.
	if identityData.Status > Resolved {
		if err := au.Service.Revoke(c, identity); err != nil {
			return err
		}
	}

	var authCookie http.Cookie
	// authCookie.HttpOnly = true
	authCookie.Name = "oauth-data"
	authCookie.Value = identityData.TokenID
	authCookie.Expires = time.Now().Add(httputil.TwentyFourHoursDuration * 3) // 3 days

	c.SetCookie(&authCookie)

	return c.Redirect(http.StatusOK, string(requestURI))
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

	fromURL := base64.StdEncoding.EncodeToString([]byte(url.String()))
	inibase := fmt.Sprintf("%s:%s:%s:%s", au.ServiceName, fromURL, randString(15), identity)
	secret := base64.StdEncoding.EncodeToString([]byte(inibase))

	redirectURL, err := au.Service.New(c, identity, secret)
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

	if err := au.Service.Revoke(c, identity); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// RetrieveAll defines a function to return a existing oauth access record through the underline
// OAuthService.
func (au AuthAPI) RetrieveAll(c *httputil.Context) error {
	identities, err := au.Service.Identities(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, identities)
}

// Retrieve defines a function to return a existing oauth access record through the underline
// OAuthService.
func (au AuthAPI) Retrieve(c *httputil.Context) error {
	identity, ok := c.GetString("identity")
	if !ok {
		// c.NoContent(http.StatusBadRequest)
		return errors.New("identity param not found")
	}

	identityInfo, err := au.Service.Get(c, identity)
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

		if err := au.Service.Authenticate(c, identity, bearer, userToken); err != nil {
			return err
		}

		return nil
	}

	if err := au.Service.Authenticate(c, identity, bearer, token); err != nil {
		return err
	}

	return nil
}

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func randString(n int) string {
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
