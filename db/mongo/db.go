package mongo

import (
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2"
)

// Mongod defines a interface which exposes a method for retrieving a
// mongo.Database and mongo.Session.
type Mongod interface {
	New() (*mgo.Database, *mgo.Session, error)
}

// Config provides configuration for connecting to a db.
type Config struct {
	Host     string
	AuthDB   string
	DB       string
	User     string
	Password string
	Mode     mgo.Mode
}

// GetSession attempts to retrieve the giving session for the given config.
func GetSession(config Config) (*mgo.Session, error) {
	// key := config.Host + ":" + config.DB

	// If not found, then attemp to connect and add to session master list.
	// We need this object to establish a session to our MongoDB.
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

	if config.Mode < 0 {
		config.Mode = mgo.Monotonic
	}

	ses.SetMode(config.Mode, true)

	return ses, nil
}

// New returns a new instance of a MongoServer.
func New(config Config) Mongod {
	var mn mongoServer
	mn.Config = config

	return &mn
}

// mongoServer defines a mongo connection manager that builds
// allows usage of a giving configuration to generate new mongo
// sessions and database instances.
type mongoServer struct {
	Config
}

// New returns a new session and database from the giving configuration.
func (m *mongoServer) New() (*mgo.Database, *mgo.Session, error) {
	ses, err := GetSession(m.Config)
	if err != nil {
		return nil, nil, err
	}

	return ses.DB(m.Config.DB), ses, nil
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
