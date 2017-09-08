package sql_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/influx6/backoffice/models/profile"
	"github.com/influx6/faux/db"
	"github.com/influx6/faux/db/sql"
	"github.com/influx6/faux/db/sql/tables"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/sentries/stdout"
	"github.com/influx6/faux/naming"
	"github.com/influx6/faux/tests"
	uuid "github.com/satori/go.uuid"
)

// contains different environment flags for use to setting up
// a db connection.
var (
	basicNamer = naming.NewNamer("%s_%s", naming.PrefixNamer{Prefix: "test"})

	log = metrics.New(stdout.Stdout{})

	DBPortEnv     = "MYSQL_PORT"
	DBIPEnv       = "MYSQL_IP"
	DBUserEnv     = "MYSQL_USER"
	DBDatabaseEnv = "MYSQL_DATABASE"
	DBUserPassEnv = "MYSQL_PASSWORD"

	user     = strings.TrimSpace(os.Getenv(DBUserEnv))
	userPass = strings.TrimSpace(os.Getenv(DBUserPassEnv))
	port     = strings.TrimSpace(os.Getenv(DBPortEnv))
	ip       = strings.TrimSpace(os.Getenv(DBIPEnv))
	dbName   = strings.TrimSpace(os.Getenv(DBDatabaseEnv))

	mydb = sql.NewDB(sql.Config{
		DBDriver:     "mysql",
		DBIP:         ip,
		DBPort:       port,
		DBName:       dbName,
		User:         user,
		UserPassword: userPass,
	}, log)
)

func TestSQLAPI(t *testing.T) {
	userTable := db.TableName{Name: basicNamer.New("users")}

	userTableMigration := tables.TableMigration{
		TableName:   basicNamer.New("users"),
		Timestamped: true,
		Indexes:     []tables.IndexMigration{},
		Fields: []tables.FieldMigration{
			{
				FieldName: "email",
				FieldType: "VARCHAR(255)",
				NotNull:   true,
			},
			{
				FieldName:  "public_id",
				FieldType:  "VARCHAR(255)",
				PrimaryKey: true,
				NotNull:    true,
			},
			{
				FieldName: "private_id",
				FieldType: "VARCHAR(255)",
				NotNull:   true,
			},
			{
				FieldName: "hash",
				FieldType: "VARCHAR(255)",
				NotNull:   true,
			},
		},
	}

	nw, err := New(NewUser{
		Email:    "bob@guma.com",
		Password: "glow",
	})

	if err != nil {
		tests.Failed("Should have successfully created new user: %+q.", err)
	}
	tests.Passed("Should have successfully created new ")

	db := sql.New(log, mydb, userTableMigration)

	t.Logf("Given the need to validate sql api operations")
	{

		t.Log("\tWhen saving user record")
		{
			if err := db.Save(userTable, nw); err != nil {
				tests.Failed("Should have successfully saved record to db table %q: %+q.", userTable.Table(), err)
			}
			tests.Passed("Should have successfully saved record to db table %q.", userTable.Table())
		}

		t.Log("\tWhen counting user records")
		{
			total, err := db.Count(userTable)
			if err != nil {
				tests.Failed("Should have successfully saved record to db table %q: %+q.", userTable.Table(), err)
			}
			tests.Passed("Should have successfully saved record to db table %q.", userTable.Table())

			if total <= 0 {
				tests.Failed("Should have successfully recieved a count greater than 0.")

			}
			tests.Passed("Should have successfully recieved a count greater than 0.")
		}

		t.Log("\tWhen retrieving all user record")
		{
			records, err := db.GetAll(userTable, "asc", "public_id")
			if err != nil {
				tests.Failed("Should have successfully retrieved all records from db table %q: %+q.", userTable.Table(), err)
			}
			tests.Passed("Should have successfully retrieved all records from db table %q.", userTable.Table())

			if len(records) == 0 {
				tests.Failed("Should have successfully retrieved atleast one record from db table %q.", userTable.Table())
			}
			tests.Passed("Should have successfully retrieved atleast one record from db table %q.", userTable.Table())
		}

		t.Log("\tWhen retrieving all user record based on page")
		{
			_, total, err := db.GetAllPerPage(userTable, "asc", "public_id", 2, 2)
			if err != nil {
				tests.Failed("Should have successfully retrieved all records from db table %q: %+q.", userTable.Table(), err)
			}
			tests.Passed("Should have successfully retrieved all records from db table %q.", userTable.Table())

			if total == -1 {
				tests.Failed("Should have successfully retrieved records based on pages from db table %q.", userTable.Table())
			}
			tests.Passed("Should have successfully retrieved records based on pages from db table %q.", userTable.Table())
		}

		t.Log("\tWhen retrieving user record")
		{
			var nu User
			if err := db.Get(userTable, &nu, "public_id", nw.PublicID); err != nil {
				tests.Failed("Should have successfully retrieved record from db table %q: %+q.", userTable.Table(), err)
			}
			tests.Passed("Should have successfully retrieved record from db table %q.", userTable.Table())

			if nu.PublicID != nw.PublicID {
				tests.Info("Expected: %+q", nw.Fields())
				tests.Info("Recieved: %+q", nu.Fields())
				tests.Failed("Should have successfully matched original user with user retrieved from db.")
			}
			tests.Passed("Should have successfully matched original user with user retrieved from db.")
		}

		t.Log("\tWhen updating user record")
		{
			if err := db.Update(userTable, nw, "public_id", nw.PublicID); err != nil {
				tests.Failed("Should have successfully updated record to db table %q: %+q.", userTable.Table(), err)
			}
			tests.Passed("Should have successfully updated record to db table %q.", userTable.Table())
		}

		t.Logf("\tWhen deleting user record")
		{
			if err := db.Delete(userTable, "public_id", nw.PublicID); err != nil {
				tests.Failed("Should have successfully deleted record to db table %q: %+q.", userTable.Table(), err)
			}
			tests.Passed("Should have successfully deleted record to db table %q.", userTable.Table())
		}
	}
}

//=============================================================================================================

const (
	hashComplexity = 10
	timeFormat     = "Mon Jan 2 15:04:05 -0700 MST 2006"
)

// NewUser defines the set of data received to create a new user.
type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User is a type defining the given user related fields for a given.
type User struct {
	Email     string           `json:"email"`
	PublicID  string           `json:"public_id"`
	PrivateID string           `json:"private_id,omitempty"`
	Hash      string           `json:"hash,omitempty"`
	Profile   *profile.Profile `json:"profile,omitempty"`
}

// New returns a new User instance based on the provided data.
func New(nw NewUser) (*User, error) {
	var u User
	u.Email = nw.Email
	u.PublicID = uuid.NewV4().String()
	u.PrivateID = uuid.NewV4().String()

	u.ChangePassword(nw.Password)

	return &u, nil
}

// Authenticate attempts to authenticate the giving password to the provided
func (u User) Authenticate(password string) error {
	pass := []byte(u.PrivateID + ":" + password)
	return bcrypt.CompareHashAndPassword([]byte(u.Hash), pass)
}

// Table returns the given table which the given struct corresponds to.
func (u User) Table() string {
	return basicNamer.New("users")
}

// SafeFields returns a map representing the data of the user with important
// security fields removed.
func (u User) SafeFields() map[string]interface{} {
	fields := u.Fields()

	delete(fields, "hash")
	delete(fields, "private_id")

	return fields
}

// Fields returns a map representing the data of the
func (u User) Fields() map[string]interface{} {
	fields := map[string]interface{}{
		"hash":       u.Hash,
		"email":      u.Email,
		"private_id": u.PrivateID,
		"public_id":  u.PublicID,
	}

	if u.Profile != nil {
		fields["profile"] = u.Profile.Fields()
	}

	return fields
}

// ChangePassword uses the provided password to set the users password hash.
func (u *User) ChangePassword(password string) error {
	pass := []byte(u.PrivateID + ":" + password)
	hash, err := bcrypt.GenerateFromPassword(pass, hashComplexity)
	if err != nil {
		return err
	}

	u.Hash = string(hash)
	return nil
}

// WithFields attempts to syncing the giving data within the provided
// map into it's own fields.
func (u *User) WithFields(fields map[string]interface{}) error {
	if email, ok := fields["email"].(string); ok {
		u.Email = email
	} else {
		return errors.New("Expected 'email' key")
	}

	if public, ok := fields["public_id"].(string); ok {
		u.PublicID = public
	} else {
		return errors.New("Expected 'public_id' key")
	}

	if private, ok := fields["private_id"].(string); ok {
		u.PrivateID = private
	} else {
		return errors.New("Expected 'private_id' key")
	}

	if hash, ok := fields["hash"].(string); ok {
		u.Hash = hash
	} else {
		return errors.New("Expected 'hash' key")
	}

	return nil
}
