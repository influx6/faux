package mongo

import (
	"encoding/json"
	"errors"
	"runtime"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

// Config embodies the data used to connect to user's mongo connection.
type Config struct {
	AuthDB     string `toml:"authdb"`
	DB         string `toml:"db"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	Host       string `toml:"host"`
	Collection string `toml:"collection"`
}

// CloneWithCollection returns a new Config cloned from this instance
// with the Collection changed to the provided name.
func (m Config) CloneWithCollection(col string) Config {
	copy := m
	copy.Collection = col
	return copy
}

// Validate returns an error if the config is invalid.
func (mgc Config) Validate() error {
	if mgc.User == "" {
		return errors.New("Config.User is required")
	}
	if mgc.Password == "" {
		return errors.New("Config.Password is required")
	}
	if mgc.AuthDB == "" {
		return errors.New("Config.AuthDB is required")
	}
	if mgc.Host == "" {
		return errors.New("Config.Host is required")
	}
	if mgc.DB == "" {
		return errors.New("Config.DB is required")
	}
	return nil
}

// MongoDB defines a mongo connection manager that builds
// allows usage of a giving configuration to generate new mongo
// sessions and database instances.
type MongoDB struct {
	Config
	ml     sync.Mutex
	master *mgo.Session
}

// NewMongoDB returns a new instance of a MongoDB.
func NewMongoDB(conf Config) *MongoDB {
	mg := &MongoDB{
		Config: conf,
	}

	// Add finalizer to ensure closure of master session.
	runtime.SetFinalizer(mg, func() {
		mg.ml.Lock()
		defer mg.ml.Unlock()
		if mg.master != nil {
			mg.master.Close()
			mg.master = nil
		}
	})

	return mg
}

// New returns a new session and database from the giving configuration.
//
// Argument:
//  isread: bool
//
// 1. If `isread` is false, then the mgo.Session is cloned so that we re-use the existing
// sessiby not closing, so others get use ofn connection, in such case, it lets you optimize writes, so try not
// the session instance connection for other writes.
//
// 2. If `isread` is true, then session is copied which creates a new unique session which you
// should close after use, this lets you handle large reads that may contain complicated queries.
//
func (m *MongoDB) New(isread bool) (*mgo.Collection, *mgo.Database, *mgo.Session, error) {
	m.ml.Lock()
	defer m.ml.Unlock()

	// if m.master is alive then continue else, reset as empty.
	if err := m.master.Ping(); err != nil {
		m.master = nil
	}

	if m.master != nil && isread {
		copy := m.master.Copy()
		db := copy.DB(m.Config.DB)
		return db.C(m.Config.Collection), db, copy, nil
	}

	if m.master != nil && !isread {
		clone := m.master.Clone()
		db := clone.DB(m.Config.DB)
		return db.C(m.Config.Collection), db, clone, nil
	}

	ses, err := getSession(m.Config)
	if err != nil {
		return nil, nil, nil, err
	}

	m.master = ses

	if isread {
		copy := m.master.Copy()
		db := copy.DB(m.Config.DB)
		return db.C(m.Config.Collection), db, copy, nil
	}

	clone := m.master.Clone()
	db := clone.DB(m.Config.DB)
	return db.C(m.Config.Collection), db, clone, nil
}

// getSession attempts to retrieve the giving session for the given config.
func getSession(config Config) (*mgo.Session, error) {
	info := mgo.DialInfo{
		Addrs:    []string{config.Host},
		Timeout:  60 * time.Second,
		Database: config.AuthDB,
		Username: config.User,
		Password: config.Password,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	ses, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil, err
	}

	ses.SetMode(mgo.Monotonic, true)

	return ses, nil
}

//==========================================================================================

// JSONIndent returns the stringified version of the giving data and indents
// its result. Uses json.Marshal underneath.
func JSONIndent(ms interface{}) (string, error) {
	data, err := json.MarshalIndent(ms, "", "\t")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// JSON returns a stringified version of the provided argument
// using json.Marshal.
func JSON(ms interface{}) (string, error) {
	data, err := json.Marshal(ms)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
