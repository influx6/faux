package oauthbolt

import (
	"bytes"
	"encoding/json"
	"time"

	"context"

	"github.com/boltdb/bolt"
	"github.com/influx6/faux/auth"
	"github.com/influx6/faux/auth/service"
	"github.com/influx6/faux/crypt"
	uuid "github.com/satori/go.uuid"
)

var (
	recordName = []byte("oauth-records")
)

// OAuthBolt defines struct which implements the OAuthService interface to
// provide OAuth authentication using boltdb as the underline session storage.
type OAuthBolt struct {
	bolt   *bolt.DB
	client *auth.Auth
}

// New returns a new instance of a OAuthBolt.
func New(client *auth.Auth) (*OAuthBolt, error) {
	var au OAuthBolt
	au.client = client

	db, err := bolt.Open("oauth-bolted.db", 0600, &bolt.Options{
		Timeout: 30 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	if uerr := db.Update(func(tx *bolt.Tx) error {
		_, cerr := tx.CreateBucketIfNotExists(recordName)
		return cerr
	}); err != nil {
		return nil, uerr
	}

	au.bolt = db

	return &au, nil
}

// Revoke attempts to revoke authorization as regarding the giving identitys and
// will remove any record associated with the identity.
func (au *OAuthBolt) Revoke(ctx context.Context, identity string) error {
	if err := au.bolt.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(recordName)

		return bucket.Delete([]byte(identity))
	}); err != nil {
		return err
	}

	return nil
}

// Approve receives the giving response and uses the underline oauth client to
// retrieve access token.
func (au *OAuthBolt) Approve(ctx context.Context, response service.IdentityResponse) error {
	if err := au.bolt.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(recordName)

		identityData := bucket.Get([]byte(response.Identity))

		_, token, err := au.client.AuthorizeFromUser(response.Code)
		if err != nil {
			return err
		}

		var data service.Identity

		if exerr := json.Unmarshal(identityData, &data); exerr != nil {
			return exerr
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

		updatedIdentity, err := json.Marshal(data)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(response.Identity), updatedIdentity)
	}); err != nil {
		return err
	}

	return nil
}

// Authenticate attempts to validate giving identity against provided token and auth type.
func (au *OAuthBolt) Authenticate(ctx context.Context, identity string, authtype string, token string) error {
	if err := au.bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(recordName)
		identityData := bucket.Get([]byte(identity))

		var data service.Identity

		if err := json.Unmarshal(identityData, &data); err != nil {
			return err
		}

		provider := token + ":" + data.Token.AccessToken
		if err := crypt.BcryptAuthenticate([]byte(data.PrivateID), []byte(provider)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// New returns a new URl for the giving identity and secret which is suited
// for requesting access.
func (au *OAuthBolt) New(ctx context.Context, identity string, secret string) (string, error) {
	var identityRequestURL string

	identityName := []byte(identity)

	if err := au.bolt.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(recordName)

		identityData := bucket.Get(identityName)

		identityRequestURL = au.client.LoginURL(secret)

		var data service.Identity

		// Unknown identity hence retrieve URL and deliver response.
		if len(identityData) == 0 {
			data.Identity = identity
			data.Status = service.Pending

			initialData, err := json.Marshal(data)
			if err != nil {
				return err
			}

			return bucket.Put(identityName, initialData)
		}

		if err := json.Unmarshal(identityData, &data); err != nil {
			return err
		}

		// If identity status is still pending then identity has not received
		// it's completion yet. Return URL to seek completion.
		if data.Status == service.Pending {
			return nil
		}

		// If we are dealing with a zero time expiration, then properly corrupted data.
		// Re-initialize authorization process.
		if data.Token.Expires.IsZero() {
			if err := bucket.Delete(identityName); err != nil {
				return err
			}

			return nil
		}

		// If we are not empty and we are still in good time, then we dont need to attempt to return url.
		if !data.Token.Expires.IsZero() && time.Now().Before(data.Token.Expires) {
			return service.ErrIdentityStillValid
		}

		if !data.Token.Expires.IsZero() && time.Now().After(data.Token.Expires) {
			bucket.Delete([]byte(identity))

			return service.ErrIdentityHasExpired
		}

		return nil
	}); err != nil {
		return identityRequestURL, err
	}

	return identityRequestURL, nil
}

// Get attempts to retrieve a identity record associated with the identity.
func (au *OAuthBolt) Get(ctx context.Context, identity string) (service.Identity, error) {
	var data service.Identity

	if err := au.bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(recordName)
		identityData := bucket.Get([]byte(identity))

		if len(identityData) == 0 {
			return service.ErrIdentityNotFound
		}

		if err := json.Unmarshal(identityData, &data); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return data, err
	}

	return data, nil
}

// Identities returns all available valid identites within the store.
func (au *OAuthBolt) Identities(ctx context.Context) ([]service.Identity, error) {
	var data []service.Identity

	if err := au.bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(recordName)

		var allRecords [][]byte

		if err := bucket.ForEach(func(_ []byte, v []byte) error {
			allRecords = append(allRecords, v)
			return nil
		}); err != nil {
			return err
		}

		var buf bytes.Buffer
		buf.WriteString("[")
		buf.Write(bytes.Join(allRecords, []byte(",")))
		buf.WriteString("]")

		if err := json.NewDecoder(&buf).Decode(&data); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return data, nil
}
