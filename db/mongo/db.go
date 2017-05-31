package mongo

import (
	"encoding/json"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

// Config provides configuration for connecting to a db.
type Config struct {
	Host     string
	AuthDB   string
	DB       string
	User     string
	Password string
}

// MongoServer defines a mongo connection manager that builds 
// allows usage of a giving configuration to generate new mongo 
// sessions and database instances.
type MongoServer struct {
	Config
	sl sync.Mutex
	sessions map[string*mgo.Session
}

// New returns a new instance of a MongoServer.
func New(config Config) *MongoServer {
	var mn MongoServer
	mn.Config = config
	mn.sessions = make(map[string]*mgo.Session)

	return &mn
}
// New returns a new session and database from the giving configuration.
func (m *Mongnod) NewSession() (*mgo.Database, *mgo.Session, error) {
	key := m.Config.Host + ":" + m.Config.DB

	m.sl.Lock()	
	ms, ok := m.sessions[key]
	m.sl.Unlock()	

	if ok {
		ses := ms.Copy()
		return ses.DB(m.Config.DB), ses, nil
	}

	// If not found, then attemp to connect and add to session master list.
	// We need this object to establish a session to our MongoDB.
	info := mgo.DialInfo{
		Addrs:    []string{m.Config.Host},
		Timeout:  60 * time.Second,
		Database: m.Config.AuthDB,
		Username: m.Config.User,
		Password: m.Config.Password,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	ses, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil, nil, err
	}

	ses.SetMode(mgo.Monotonic, true)

	// Add to master list.
	m.sl.Lock()	
	m.sessions[key] = ses.Copy()
	m.sl.Unlock()	

	return ses.DB(m.Config.DB), ses, nil
}

//==========================================================================================

// QueryIndent returns the stringified version of the giving data and indents
// its result. Uses json.Marshal underneath.
func QueryIndent(ms interface{}) string {
	data, err := json.MarshalIndent(ms, "", "\n")
	if err != nil {
		return ""
	}

	return string(data)
}

// Query returns a stringified version of the provided argument
// using json.Marshal.
func Query(ms interface{}) string {
	data, err := json.Marshal(ms)
	if err != nil {
		return ""
	}

	return string(data)
}
