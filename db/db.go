package db

// TableIdentity defines an interface which exposes a method returning table name
// associated with the giving implementing structure.
type TableIdentity interface {
	Table() string
}

// TableFields defines an interface which exposes method to return a map of all data
// associated with the defined structure.
type TableFields interface {
	Fields() map[string]interface{}
}

// TableConsumer defines an interface which accepts a map of data which will be loaded
// into the giving implementing structure.
type TableConsumer interface {
	WithFields(map[string]interface{}) error
}

// Consumer defines an interface that exposes a Consume method.
type Consumer interface {
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

//=============================================================================================================================================

// TableName defines a struct which returns a given table name associated with the table.
type TableName struct {
	Name string
}

// Table returns the giving name associated with the struct.
func (t TableName) Table() string {
	return t.Name
}

//=============================================================================================================================================

// Namer exposes a method which can be used to generate specific names based on a giving
// critieras.
type Namer interface {
	New(string) string
}

//=============================================================================================================================================

// TableNamer defines holds a underline naming mechanism to deliver new TableName instance.
type TableNamer struct {
	namer Namer
}

// NewTableNamer returns a new TableNamer instance.
func NewTableNamer(nm Namer) *TableNamer {
	return &TableNamer{
		namer: nm,
	}
}

// New returns a new TableName which is fed into the underline naming mechanism to
// generate a unique name for that table.
// eg namer = FeedNamer("sugo_company");  TableNamer(namer).New("users") => "sugo_company_users".
func (t *TableNamer) New(table string) TableName {
	return TableName{Name: t.namer.New(table)}
}

//=============================================================================================================================================
