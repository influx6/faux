package mongo

import (
	"encoding/json"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

//==============================================================================

var mstore struct {
	ml   sync.RWMutex
	list map[string]*mgo.Session
}

// masterListLock provides a mutex for controlling access to the masterList.
var masterListLock sync.RWMutex

// masterList contains a set of session lists of connections that have been
// created
var masterList = make(map[string]*mgo.Session)

//==============================================================================

// Config provides configuration for connecting to a db.
type Config struct {
	Host     string
	AuthDB   string
	DB       string
	User     string
	Password string
}

//==============================================================================

// EventLog defines event logger that allows us to record events for a specific
// action that occured.
type EventLog interface {
	Dev(context interface{}, name string, message string, data ...interface{})
	User(context interface{}, name string, message string, data ...interface{})
	Error(context interface{}, name string, err error, message string, data ...interface{})
}

//==============================================================================

// Mongnod defines a mongo connection manager that builds off a mongo instance.
type Mongnod struct {
	C   Config
	Log EventLog
}

// New connects and initializes the master session for the mongo list.
func (m *Mongnod) New(context interface{}) (*mgo.Database, *mgo.Session, error) {
	key := m.C.Host + ":" + m.C.DB

	masterListLock.Lock()
	ms, ok := masterList[key]
	masterListLock.Unlock()

	if ok {
		ses := ms.Copy()
		return ses.DB(m.C.DB), ses, nil
	}

	// If not found, then attemp to connect and add to session master list.
	// We need this object to establish a session to our MongoDB.
	info := mgo.DialInfo{
		Addrs:    []string{m.C.Host},
		Timeout:  60 * time.Second,
		Database: m.C.AuthDB,
		Username: m.C.User,
		Password: m.C.Password,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	ses, err := mgo.DialWithInfo(&info)
	if err != nil {
		m.Log.Error(context, "New", err, "Completed")
		return nil, nil, err
	}

	ses.SetMode(mgo.Monotonic, true)

	// Add to master list.
	masterListLock.Lock()
	masterList[key] = ses.Copy()
	masterListLock.Unlock()

	return ses.DB(m.C.DB), ses, nil
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
