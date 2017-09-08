package service_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/influx6/faux/auth"
	"github.com/influx6/faux/auth/google"
	"github.com/influx6/faux/auth/service"
	"github.com/influx6/faux/auth/service/providers/oauthbolt"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/sentries/stdout"
	"github.com/influx6/faux/tests"
)

var (
	id           = "323"
	clientId     = "43434as43423d232fr232"
	clientSecret = "af3434Ju83434HK23232"
	events       = metrics.New(stdout.Stderr{})
)

func TestAuthAPI(t *testing.T) {
	defer os.Remove("oauth-bolted.db")

	client := google.New(auth.Credential{
		ClientID:     clientId,
		ClientSecret: clientSecret,
	}, "http://localhost:80/")

	googleService, err := oauthbolt.New(events, client)
	if err != nil {
		tests.Failed("Should have successfully created bolt service")
	}
	tests.Passed("Should have successfully created bolt service")

	api := service.New(googleService, "google")

	testRegister(t, api)
	testRetrieval(t, api)
	testRetrievalAll(t, api)
	testRevoking(t, api)
	testRetrievalAllWithNoRecord(t, api)
}

func testRetrievalAllWithNoRecord(t *testing.T, api service.AuthAPI) {
	newLoginCtx, newLoginRw := newContext(newRequest(map[string]string{}, nil))

	if err := api.RetrieveAll(newLoginCtx); err != nil {
		tests.Failed("Should have processed request successfully: %+q", err)
	}
	tests.Passed("Should have processed request successfully")

	var res []service.Identity

	if err := json.NewDecoder(newLoginRw.Body).Decode(&res); err != nil {
		tests.Failed("Should have successfully decoded request response: %+q.", err)
	}
	tests.Passed("Should have successfully decoded request response.")

	if len(res) == 1 {
		tests.Failed("Should have successfully received empty records.")
	}
	tests.Passed("Should have successfully received empty records.")
}

func testRetrievalAll(t *testing.T, api service.AuthAPI) {
	newLoginCtx, newLoginRw := newContext(newRequest(map[string]string{}, nil))

	if err := api.RetrieveAll(newLoginCtx); err != nil {
		tests.Failed("Should have processed request successfully: %+q", err)
	}
	tests.Passed("Should have processed request successfully")

	var res []service.Identity

	if err := json.NewDecoder(newLoginRw.Body).Decode(&res); err != nil {
		tests.Failed("Should have successfully decoded request response: %+q.", err)
	}
	tests.Passed("Should have successfully decoded request response.")

	if len(res) != 1 {
		tests.Failed("Should have successfully received 1 record.")
	}
	tests.Passed("Should have successfully received 1 record.")
}

func testRevoking(t *testing.T, api service.AuthAPI) {
	newLoginCtx, newLoginRw := newContext(newRequest(map[string]string{
		"identity": id,
	}, nil))

	if err := api.Revoke(newLoginCtx); err != nil {
		tests.Failed("Should have processed request successfully: %+q", err)
	}
	tests.Passed("Should have processed request successfully")

	if newLoginRw.Code != http.StatusNoContent {
		tests.Failed("Should have successfully received StatusNoContent code.")
	}
	tests.Passed("Should have successfully received StatusNoContent code.")
}

func testRetrieval(t *testing.T, api service.AuthAPI) {
	newLoginCtx, newLoginRw := newContext(newRequest(map[string]string{
		"identity": id,
	}, nil))

	if err := api.Retrieve(newLoginCtx); err != nil {
		tests.Failed("Should have processed request successfully: %+q", err)
	}
	tests.Passed("Should have processed request successfully")

	var res service.Identity

	if err := json.NewDecoder(newLoginRw.Body).Decode(&res); err != nil {
		tests.Failed("Should have successfully decoded request response: %+q.", err)
	}
	tests.Passed("Should have successfully decoded request response.")

	if res.Identity != id {
		tests.Failed("Should have successfully matched request identity with response.identity.")
	}
	tests.Passed("Should have successfully matched request identity with response.identity.")
}

func testRegister(t *testing.T, api service.AuthAPI) {
	newLoginCtx, newLoginRw := newContext(newRequest(map[string]string{
		"identity": id,
	}, nil))

	if err := api.Register(newLoginCtx); err != nil {
		tests.Failed("Should have processed request successfully: %+q", err)
	}
	tests.Passed("Should have processed request successfully")

	if newLoginRw.Code != http.StatusOK {
		tests.Failed("Should have successfully matched expected OK status: %+q.", newLoginRw.Code)
	}
	tests.Passed("Should have successfully matched expected OK status.")

	var res service.IdentityPath

	if err := json.NewDecoder(newLoginRw.Body).Decode(&res); err != nil {
		tests.Failed("Should have successfully decoded request response: %+q.", err)
	}
	tests.Passed("Should have successfully decoded request response.")

	if res.Identity != id {
		tests.Failed("Should have successfully matched request identity with response.identity.")
	}
	tests.Passed("Should have successfully matched request identity with response.identity.")

	loginURl, err := url.Parse(res.Login)
	if err != nil {
		tests.Failed("Should have successfully parsed request url : %+q.", err)
	}
	tests.Passed("Should have successfully parsed request url.")

	stateValue := loginURl.Query().Get("state")

	decodeStateValue, err := base64.StdEncoding.DecodeString(stateValue)
	if err != nil {
		tests.Failed("Should have successfully parsed request url.secret : %+q.", err)
	}
	tests.Passed("Should have successfully decoded request url.secret.")

	encodedURL := base64.StdEncoding.EncodeToString([]byte(newLoginCtx.Request().URL.String()))

	if !strings.Contains(string(decodeStateValue), string(encodedURL)) {
		tests.Failed("Should have successfully found request url in state secret : %+q.", err)
	}
	tests.Passed("Should have successfully found request url in state secret.")

}

func newRequest(params map[string]string, body io.Reader) *http.Request {
	var uparams []string

	for key, name := range params {
		uparams = append(uparams, fmt.Sprintf("%s=%s", key, name))
	}

	newReq, err := http.NewRequest("GET", "http://localhost:80/?"+strings.Join(uparams, "&"), body)
	if err != nil {
		tests.Failed("Should have successfully created request: %q.", err)
	}
	tests.Passed("Should have successfully created request.")

	return newReq
}

func newContext(r *http.Request) (*httputil.Context, *httptest.ResponseRecorder) {
	res := httptest.NewRecorder()
	ctx := httputil.NewContext(httputil.SetRequest(r),
		httputil.SetResponseWriter(res), httputil.SetMetrics(events))

	return ctx, res
}
