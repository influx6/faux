// Package oauthmongo exposes a structure which exposes a service to handle storage
// of oauth related authentication and authorization.
// @mongo
package oauthmongo

import (
	"errors"
	"time"

	"github.com/influx6/faux/auth"
	"github.com/influx6/faux/auth/service"
	"github.com/influx6/faux/auth/service/providers/oauthmongo/mongo"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/crypt"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/sentries/stdout"
	uuid "github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// OAuthMongo defines struct which implements the OAuthService interface to
// provide OAuth authentication using boltdb as the underline session storage.
type OAuthMongo struct {
	col     string
	mgo     *mongo.DB
	client  *auth.Auth
	metrics metrics.Metrics
}

// New returns a new instance of a OAuthMongo.
func New(col string, metric metrics.Metrics, db mongo.Mongod, client *auth.Auth) (*OAuthMongo, error) {
	var au OAuthMongo
	au.col = col
	au.client = client
	au.metrics = metric
	au.mgo = mongo.New(metric, db)

	if err := au.mgo.WithIndex(context.New(), col, mgo.Index{
		Key:      []string{"identity"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	}); err != nil {
		return &au, err
	}

	return &au, nil
}

// Revoke attempts to revoke authorization as regarding the giving identitys and
// will remove any record associated with the identity.
func (au *OAuthMongo) Revoke(ctx context.Context, identity string) error {
	m := stdout.Info("OAuthMongo.Revoke").Trace()
	defer au.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := errors.New("Context has expired")
		au.metrics.Emit(stdout.Error("OAuthBolt.Revoke").
			With("error", err).With("identity", identity).With("collection", au.col))
		return err
	}

	if err := au.mgo.Exec(ctx, au.col, func(col *mgo.Collection) error {
		return col.Remove(bson.M{"identity": identity})
	}); err != nil {
		au.metrics.Emit(stdout.Error("OAuthMongo.Revoke").
			With("error", err).With("identity", identity))

		return err
	}

	au.metrics.Emit(stdout.Info("Completed : OAuthMongo.Revoke").With("identity", identity))

	return nil
}

// Approve receives the giving response and uses the underline oauth client to
// retrieve access token.
func (au *OAuthMongo) Approve(ctx context.Context, response service.IdentityResponse) error {
	m := stdout.Info("OAuthMongo.Approve").Trace()
	defer au.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := errors.New("Context has expired")
		au.metrics.Emit(stdout.Error("OAuthBolt.Approve").
			With("error", err).With("identity", response.Identity).With("collection", au.col))
		return err
	}

	if err := au.mgo.Exec(ctx, au.col, func(col *mgo.Collection) error {
		findQuery := bson.M{
			"identity": response.Identity,
		}

		var data service.Identity

		if err := col.Find(findQuery).One(&data); err != nil {
			return err
		}

		_, token, err := au.client.AuthorizeFromUser(response.Code)
		if err != nil {
			return err
		}

		if data.Data == nil {
			data.Data = make(map[string]interface{})
		}

		for key, val := range response.Data {
			data.Data[key] = val
		}

		data.Token = token
		data.Status = service.Resolved

		tokenUUID := uuid.NewV4().String()

		tokenID, err := crypt.BcryptGenerate([]byte(tokenUUID+":"+token.AccessToken), 20)
		if err != nil {
			return err
		}

		data.TokenID = tokenUUID
		data.PrivateID = string(tokenID)

		return col.Update(findQuery, bson.M{
			"$set": bson.M{
				"status":     data.Status,
				"data":       data.Data,
				"token":      data.Token,
				"token_id":   data.TokenID,
				"private_id": data.PrivateID,
			},
		})
	}); err != nil {
		au.metrics.Emit(stdout.Error("OAuthMongo.Approve").
			With("error", err).With("identity", response.Identity).With("collection", au.col))

		return err
	}

	au.metrics.Emit(stdout.Info("Completed : OAuthMongo.Approve").
		With("response.identity", response.Identity).With("response.data", response.Data))

	return nil
}

// Authenticate attempts to validate giving identity against provided token and auth type.
func (au *OAuthMongo) Authenticate(ctx context.Context, identity string, authtype string, token string) error {
	m := stdout.Info("OAuthMongo.Authenticate").Trace()
	defer au.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := errors.New("Context has expired")
		au.metrics.Emit(stdout.Error("OAuthBolt.Revoke").
			With("error", err).With("identity", identity).With("collection", au.col))

		return err
	}

	if err := au.mgo.Exec(ctx, au.col, func(col *mgo.Collection) error {
		findQuery := bson.M{
			"identity": identity,
		}

		var data service.Identity

		if err := col.Find(findQuery).One(&data); err != nil {
			return err
		}

		provider := token + ":" + data.Token.AccessToken
		if err := crypt.BcryptAuthenticate([]byte(data.PrivateID), []byte(provider)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		au.metrics.Emit(stdout.Error("OAuthMongo.Authenticate").
			With("error", err).With("identity", identity).With("collection", au.col))
		return err
	}

	au.metrics.Emit(stdout.Info("Completed : OAuthMongo.Authenticate").
		With("identity", identity))

	return nil
}

// New returns a new URl for the giving identity and secret which is suited
// for requesting access.
func (au *OAuthMongo) New(ctx context.Context, identity string, secret string) (string, error) {
	m := stdout.Info("OAuthMongo.New").With("identity", identity).Trace()
	defer au.metrics.Emit(m.End())

	var identityRequestURL string

	if ctx.IsExpired() {
		err := errors.New("Context has expired")
		au.metrics.Emit(stdout.Error("OAuthBolt.Revoke").
			With("error", err).With("identity", identity).With("collection", au.col))
		return identityRequestURL, err
	}

	if err := au.mgo.Exec(ctx, au.col, func(col *mgo.Collection) error {
		findQuery := bson.M{
			"identity": identity,
		}

		var data service.Identity

		if err := col.Find(findQuery).One(&data); err != nil {
			data.Identity = identity
			data.Status = service.Pending

			return col.Insert(data)
		}

		identityRequestURL = au.client.LoginURL(secret)

		// If identity status is still pending then identity has not received
		// it's completion yet. Return URL to seek completion.
		if data.Status == service.Pending {
			return nil
		}

		// If we are dealing with a zero time expiration, then properly corrupted data.
		// Re-initialize authorization process.
		if data.Token.Expires.IsZero() {
			if err := col.Remove(findQuery); err != nil {
				return err
			}

			return nil
		}

		// If we are not empty and we are still in good time, then we dont need to attempt to return url.
		if !data.Token.Expires.IsZero() && time.Now().Before(data.Token.Expires) {
			return service.ErrIdentityStillValid
		}

		if !data.Token.Expires.IsZero() && time.Now().After(data.Token.Expires) {
			col.Remove(findQuery)

			return service.ErrIdentityHasExpired
		}

		return nil
	}); err != nil {
		au.metrics.Emit(stdout.Error("OAuthMongo.New").
			With("error", err).With("identity", identity).With("collection", au.col))
		return identityRequestURL, err
	}

	au.metrics.Emit(stdout.Info("Completed : OAuthMongo.New").With("identity", identity))

	return identityRequestURL, nil
}

// Get attempts to retrieve a identity record associated with the identity.
func (au *OAuthMongo) Get(ctx context.Context, identity string) (service.Identity, error) {
	m := stdout.Info("OAuthMongo.Get").Trace()
	defer au.metrics.Emit(m.End())

	var data service.Identity

	if ctx.IsExpired() {
		err := errors.New("Context has expired")
		au.metrics.Emit(stdout.Error("OAuthBolt.Revoke").
			With("error", err).With("identity", identity).With("collection", au.col))
		return data, err
	}

	if err := au.mgo.Exec(ctx, au.col, func(col *mgo.Collection) error {
		findQuery := bson.M{
			"identity": identity,
		}

		if err := col.Find(findQuery).One(&data); err != nil {
			if err == mgo.ErrNotFound {
				return service.ErrIdentityNotFound
			}

			return err
		}

		return nil
	}); err != nil {
		au.metrics.Emit(stdout.Error("OAuthMongo.Get").
			With("error", err).With("identity", identity).With("collection", au.col))
		return data, err
	}

	au.metrics.Emit(stdout.Info("Completed : OAuthMongo.Get").With("identity", identity))
	return data, nil
}

// Identities returns all available valid identites within the store.
func (au *OAuthMongo) Identities(ctx context.Context) ([]service.Identity, error) {
	if ctx.IsExpired() {
		err := errors.New("Context has expired")
		au.metrics.Emit(stdout.Error("OAuthBolt.Revoke").
			With("error", err).With("collection", au.col))
		return nil, err
	}

	var items []service.Identity

	if err := au.mgo.Exec(ctx, au.col, func(col *mgo.Collection) error {
		return col.Find(nil).All(&items)
	}); err != nil {
		au.metrics.Emit(stdout.Error("OAuthMongo.Identities").
			With("error", err).With("collection", au.col))
		return nil, err
	}

	au.metrics.Emit(stdout.Info("Completed : OAuthMongo.Identities"))
	return items, nil
}
