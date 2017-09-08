package oauthbolt_test

import (
	"os"
	"testing"

	"github.com/influx6/faux/auth"
	"github.com/influx6/faux/auth/google"
	"github.com/influx6/faux/auth/service"
	"github.com/influx6/faux/auth/service/providers/oauthbolt"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/sentries/stdout"
	"github.com/influx6/faux/tests"
)

var (
	id           = "323"
	clientId     = "43434as43423d232fr232"
	clientSecret = "af3434Ju83434HK23232"
	events       = metrics.New(stdout.Stderr{})

	expectedURL = `https://accounts.google.com/o/oauth2/auth?client_id=43434as43423d232fr232&redirect_uri=http%3A%2F%2Flocalhost%3A80%2F&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile&state=RestMan%3A323`
)

func TestOAuthBolt(t *testing.T) {
	defer os.Remove("oauth-bolted.db")
	ctx := context.New()

	client := google.New(auth.Credential{
		ClientID:     clientId,
		ClientSecret: clientSecret,
	}, "http://localhost:80/")

	bolt, err := oauthbolt.New(events, client)
	if err != nil {
		tests.Failed("Should have successfully created bolt service")
	}
	tests.Passed("Should have successfully created bolt service")

	newURL, err := bolt.New(ctx, id, "RestMan:"+id)
	if err != nil {
		tests.Failed("Should have successfully received login url: %+q", err)
	}
	tests.Passed("Should have successfully received login url.")

	if newURL != expectedURL {
		tests.Info("Received URL: %q", newURL)
		tests.Info("Expected URL: %q", expectedURL)
		tests.Failed("Should have correctly matched expected URL")
	}
	tests.Passed("Should have correctly matched expected URL")

	identity, err := bolt.Get(ctx, id)
	if err != nil {
		tests.Failed("Should have successfully retrieved identity: %+q", err)
	}
	tests.Passed("Should have successfully received identity")

	if identity.Identity != id {
		tests.Info("Received ID: %q", identity.Identity)
		tests.Info("Expected URL: %q", id)
		tests.Failed("Should have correctly match identity with id")
	}
	tests.Passed("Should have correctly match identity with id")

	if identity.Status > service.Pending {
		tests.Info("Received Status: %q", identity.Status)
		tests.Info("Expected Status: %q", service.Pending)
		tests.Failed("Should have correctly match identity as pending")
	}
	tests.Passed("Should have correctly match identity as pending")

	identities, err := bolt.Identities(ctx)
	if err != nil {
		tests.Failed("Should have successfully retrieved stored identities: %+q", err)
	}
	tests.Passed("Should have successfully received stored identities")

	if len(identities) < 1 {
		tests.Failed("Should have retrieved alteast one identity record: 0 found")
	}
	tests.Passed("Should have retrieved alteast one identity record.")
}
