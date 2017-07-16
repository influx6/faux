package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"

	"encoding/json"

	"encoding/base64"

	"strings"

	"github.com/dimfeld/httptreemux"
	"github.com/influx6/faux/auth"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/metrics"
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

// Register defines a function to create a new oauth request for the underline
// OAuthService.
func (au AuthAPI) Register(c *httputil.Context) error {
	url := c.Request().URL
	identity := c.QueryParam("identity")

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

	return nil
}

// Retreive defines a function to return a existing oauth access record through the underline
// OAuthService.
func (au AuthAPI) Retrieve(c *httputil.Context) error {

	return nil
}

// Authenticate defines a function to validate a received token against a
// existing oauth access record through the underline OAuthService.
func (au AuthAPI) Authenticate(c *httputil.Context) error {

	return nil
}

// OAuthRelay defines a http service structure which registers giving OAuthService provides
// for giving OAuth service providers.
type OAuthRelay struct {
	*httptreemux.TreeMux
	metrics   metrics.Metrics
	providers map[string]OAuthService
}

// New returns a new instance of a OAuthRelay which registers giving OAuthService providers
// intro appropriate route groups for request processing.
// All services must be accessed with request having prefix routes "/oauth/", such
// that a giving sevice will be accessed through "/oauth/SERVICENAME".
// All services provide the following endpoints:
//
//	Note that all SERVICENAME will be in lowercase.
//
//	- CallbackURL (GET /oauth/response)
//		This endpoint will provide the redirectURL as it's response will be processed and stored here as
//		needed with regards to status error if not stil.
//
//	- Revoke (DELETE /oauth/SERVICENAME/:identity)
//		This endpoint will issue a correlating oauth token record revokation where associated
//		identity will be removed and clear if existing.
//
// 	- Register (GET /oauth/SERVICENAME/:identity)
//		This endpoint will issue a new authentication url has a JSON payload which can be delievred to the
//		user for authentication and authorization, every service is expected to include the needed
// 		scopes in the url that is delivered.
//
//	- New (GET /oauth/SERVICENAME/:identity/token)
//		This endpoint will issue a correlating oauth token record associated with the given identity else
//		return appropriate status error.
//
//	- Get (GET /oauth/SERVICENAME/:identity/auth)
//		This endpoint will issue a request to validate a giving request authorization header matches against
//		the access token record.
//
func New(metrics metrics.Metrics) *OAuthRelay {
	rl := &OAuthRelay{
		metrics:   metrics,
		TreeMux:   httptreemux.New(),
		providers: make(map[string]OAuthService),
	}

	rl.GET("/oauth/identity/response", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		var stateError, stateSecret, stateCode string

		if delError := r.FormValue("error"); delError != "" {
			stateError = delError
		} else {
			stateError = params["error"]
		}

		if state := r.FormValue("state"); state != "" {
			stateSecret = state
		} else {
			stateSecret = params["state"]
		}

		if code := r.FormValue("code"); code != "" {
			stateCode = code
		} else {
			stateCode = params["code"]
		}

		if stateSecret == "" {
			// httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to retrieve request secret", errors.New("State secret not found else was empty"))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// We issue a state secret where its a combination in this format: SERVICENAME:REQUEST-URI:IDENTITY.
		// If the secret does not match this format then it didnt come from us.
		decodedSecret, err := base64.StdEncoding.DecodeString(stateSecret)
		if err != nil {
			// httputil.WriteErrorMessage(w, http.StatusUnauthorized, "Failed to decode request secret", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		sections := strings.Split(string(decodedSecret), ":")
		if len(sections) != 3 {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			// httputil.WriteErrorMessage(w, http.StatusUnauthorized, "Failed to decode request secret", err)
			return
		}

		identity := sections[2]
		requestURI := sections[1]
		serviceName := sections[0]

		provider, ok := rl.providers[serviceName]
		if !ok {
			http.Redirect(w, r, requestURI, http.StatusTemporaryRedirect)
			// httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to find provider for service", fmt.Errorf("ServiceName %q provider not allowed", serviceName))
			return
		}

		var response IdentityResponse
		response.Code = stateCode
		response.Identity = identity

		if err := provider.Process(identity, response); err != nil {
			httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to process response in authorization process", err)
			return
		}

		_ = stateCode
		_ = stateError
		_ = stateSecret
	})

	return rl
}

// Register adds the giving *auth.Auth under the underline service namespace.
func (rl *OAuthRelay) Register(service string, provider OAuthService) {
	service = strings.ToLower(service)

	rl.providers[service] = provider

	authGroup := rl.TreeMux.NewGroup(fmt.Sprintf("/oauth/%s", service))

	authGroup.GET("/:identity", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		identity, ok := params["identity"]
		if !ok {
			httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to retrieve identity", errors.New("Expected identity params"))
			return
		}

		inibase := fmt.Sprintf("%s:%s:%s:%s", service, r.URL.RequestURI(), randString(15), identity)
		secret := base64.StdEncoding.EncodeToString([]byte(inibase))

		redirectURL, err := provider.New(identity, secret)
		if err != nil {
			httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to retrieve redirect url", err)
			return
		}

		if err := json.NewEncoder(w).Encode(IdentityPath{
			Identity: identity,
			Login:    redirectURL,
		}); err != nil {
			httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to respond with data", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	authGroup.DELETE("/:identity", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		identity, ok := params["identity"]
		if !ok {
			httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to retrieve identity", errors.New("Expected identity params"))
			return
		}

		if err := provider.Revoke(identity); err != nil {
			httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to respond with data", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	authGroup.GET("/:identity/token", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		identity, ok := params["identity"]
		if !ok {
			httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to retrieve identity", errors.New("Expected identity params"))
			return
		}

		identityInfo, err := provider.Get(identity)
		if err != nil {
			httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to retrieve redirect url", err)
			return
		}

		if err := json.NewEncoder(w).Encode(identityInfo); err != nil {
			httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to respond with data", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	authGroup.GET("/:identity/auth", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		identity, ok := params["identity"]
		if !ok {
			httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to retrieve identity", errors.New("Expected identity params"))
			return
		}

		authorization := r.Header.Get("Authorization")
		// if !ok {
		// 	httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to retrieve authorization header", errors.New("Expected Authorization headers"))
		// 	return
		// }

		bearer, token, err := auth.ParseAuthorization(authorization)
		if err != nil {
			httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to parse authorization header", errors.New("Expected Authorization value to match format: 'AUTH_TYPE AUTH_TOKEN'"))
			return
		}

		// if we are encoded, then it means token is base64 encode and we need to split it for the real token.
		if encodedAuthorization := r.Header.Get("Base64Authorization"); encodedAuthorization != "" {
			userIdentity, userToken, err := auth.ParseToken(token)
			if err != nil {
				httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to parse authorization header", errors.New("Expected token value"))
				return
			}

			// if userIdentity does not match given identity then fail this.
			if userIdentity != identity {
				httputil.WriteErrorMessage(w, http.StatusBadRequest, "Decoded token user identity does not match request identity", errors.New("Identity does not match token identity"))
				return
			}

			if err := provider.Authenticate(identity, bearer, userToken); err != nil {
				httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to authorizate Authorization header", errors.New("Token does not validate with in-house access token"))
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		if err := provider.Authenticate(identity, bearer, token); err != nil {
			httputil.WriteErrorMessage(w, http.StatusBadRequest, "Failed to authorizate Authorization header", errors.New("Token does not validate with in-house access token"))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
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
