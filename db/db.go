package db

// TableName defines a struct which returns a given table name associated with the table.
type TableName struct {
	Name string
}

// Table returns the giving name associated with the struct.
func (t TableName) Table() string {
	return t.Name
}

// TableIdentity defines an interface which exposes a method returning table name
// associated with the giving implementing structure.
type TableIdentity interface {
	Table() string
}

// TableFields defines an interface which exposes method to return a map of all data
// associated with the defined structure.
type TableFields interface {
	Fields() (map[string]interface{}, error)
}

// TableConsumer defines an interface that exposes a Consume method.
type TableConsumer interface {
	Consume(map[string]interface{}) error
}

// Migration defines an interface which provides structures to setup a new db migration
// call.
type Migration interface {
	Migrate() error
}

// DB defines a type which allows CRUD operations provided by a underline
// db structure.
type DB interface {
	Save(t TableIdentity, f TableFields) error
	Count(t TableIdentity) (int, error)
	Update(t TableIdentity, f TableFields, index string) error
	Delete(t TableIdentity, index string, value interface{}) error
	Get(t TableIdentity, c TableConsumer, index string, value interface{}) error
	GetAll(t TableIdentity, order string, orderBy string) ([]map[string]interface{}, error)
	GetAllPerPage(t TableIdentity, order string, orderBy string, page int, responsePage int) ([]map[string]interface{}, int, error)
}
