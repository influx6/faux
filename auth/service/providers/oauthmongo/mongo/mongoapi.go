// Package mongoapi provides a auto-generated package which contains a mongo base pkg for db operations.
//
//
package mongo

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"

	"github.com/influx6/faux/context"

	"github.com/influx6/faux/metrics"

	"github.com/influx6/faux/metrics/sentries/stdout"
)

// Mongod defines a interface which exposes a method for retrieving a
// mongo.Database and mongo.Session.
type Mongod interface {
	New() (*mgo.Database, *mgo.Session, error)
}

// DB defines a structure which provide DB CRUD operations
// using mongo as the underline db.
type DB struct {
	db      Mongod
	metrics metrics.Metrics
}

// New returns a new instance of DB.
func New(m metrics.Metrics, mo Mongod) *DB {
	return &DB{
		db:      mo,
		metrics: m,
	}
}

// WithIndex applies the provided index slice to the provided collection configuration.
func (mdb *DB) WithIndex(ctx context.Context, col string, indexes ...mgo.Index) error {
	m := stdout.Info("DB.WithIndex").Trace("DB.WithIndex")
	defer mdb.metrics.Emit(m.End())

	if len(indexes) == 0 {
		return nil
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to create session for indexes").
			With("collection", col).
			With("error", err.Error()))
		return err
	}

	defer session.Close()

	collection := database.C(col)

	for _, index := range indexes {

		if err := collection.EnsureIndex(index); err != nil {
			mdb.metrics.Emit(stdout.Error("Failed to ensure session index").
				With("collection", col).
				With("index", index).
				With("error", err.Error()))

			return err
		}

		mdb.metrics.Emit(stdout.Info("Succeeded in ensuring collection index").
			With("collection", col).
			With("index", index))
	}

	mdb.metrics.Emit(stdout.Notice("Finished adding index").
		With("collection", col))

	return nil
}

// Exec provides a function which allows the execution of a custom function against the collection.
func (mdb *DB) Exec(ctx context.Context, col string, fx func(col *mgo.Collection) error) error {
	m := stdout.Info("DB.Exec").Trace("DB.Exec")
	defer mdb.metrics.Emit(m.End())

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to execute operation").
			With("collection", col).
			With("error", err.Error()))
		return err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to create session").
			With("collection", col).
			With("error", err.Error()))
		return err
	}

	defer session.Close()

	if ctx.IsExpired() {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(stdout.Error("Failed to finish, context has expired").
			With("collection", col).
			With("error", err.Error()))
		return err
	}

	if err := fx(database.C(col)); err != nil {
		mdb.metrics.Emit(stdout.Error("Failed to execute operation").
			With("collection", col).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(stdout.Notice("Operation executed").
		With("collection", col))

	return nil
}
