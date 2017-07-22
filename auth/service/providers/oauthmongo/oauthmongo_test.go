package oauthmongo_test

import (
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/influx6/faux/auth"
	"github.com/influx6/faux/auth/google"
	"github.com/influx6/faux/auth/service"
	"github.com/influx6/faux/auth/service/providers/oauthmongo"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/db/mongo"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/sentries/stdout"
	"github.com/influx6/faux/tests"
)

var (
	id           = "323"
	clientId     = "43434as43423d232fr232"
	clientSecret = "af3434Ju83434HK23232"
	events       = metrics.New(stdout.Stderr{})

	config = mongo.Config{
		Mode:     mgo.Monotonic,
		DB:       os.Getenv("dap_MONGO_DB"),
		Host:     os.Getenv("dap_MONGO_HOST"),
		User:     os.Getenv("dap_MONGO_USER"),
		AuthDB:   os.Getenv("dap_MONGO_AUTHDB"),
		Password: os.Getenv("dap_MONGO_PASSWORD"),
	}

	testCol = "ignitor_test_collection"

	expectedURL = `https://accounts.google.com/o/oauth2/auth?client_id=43434as43423d232fr232&redirect_uri=http%3A%2F%2Flocalhost%3A80%2F&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile&state=RestMan%3A323`
)

func TestOAuthMongo(t *testing.T) {
	defer os.Remove("oauth-bolted.db")
	ctx := context.New()

	client := google.New(auth.Credential{
		ClientID:     clientId,
		ClientSecret: clientSecret,
	}, "http://localhost:80/")

	mgod, err := oauthmongo.New(testCol, events, mongo.New(config), client)
	if err != nil {
		tests.Failed("Should have successfully created mongo service")
	}
	tests.Passed("Should have successfully created mongo service")

	newURL, err := mgod.New(ctx, id, "RestMan:"+id)
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

	identity, err := mgod.Get(ctx, id)
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

	identities, err := mgod.Identities(ctx)
	if err != nil {
		tests.Failed("Should have successfully retrieved stored identities: %+q", err)
	}
	tests.Passed("Should have successfully received stored identities")

	if len(identities) < 1 {
		tests.Failed("Should have retrieved alteast one identity record: 0 found")
	}
	tests.Passed("Should have retrieved alteast one identity record.")
}
