package sql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/influx6/faux/db"
	"github.com/influx6/faux/db/sql/tables"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/sentries/stdout"
	"github.com/jmoiron/sqlx"
)

// contains templates of sql statement for use in operations.
const (
	countTemplate         = "SELECT count(*) FROM %s"
	selectAllTemplate     = "SELECT * FROM %s ORDER BY %s %s"
	selectLimitedTemplate = "SELECT * FROM %s ORDER BY %s %s LIMIT %d OFFSET %d"
	selectItemTemplate    = "SELECT * FROM %s WHERE %s=%s"
	insertTemplate        = "INSERT INTO %s %s VALUES %s"
	updateTemplate        = "UPDATE %s SET %s WHERE %s=%s"
	deleteTemplate        = "DELETE FROM %s WHERE %s=%s"
)

//===============================================================================================================

// Config is a configuration struct for the DB connection for DBMaker.
type Config struct {
	User         string `json:"user"`
	UserPassword string `json:"user_password"`
	DBPort       string `json:"db_port"`
	DBIP         string `json:"dp_ip"`
	DBName       string `json:"db_name"`
	DBDriver     string `json:"db_driver"`
}

// dBMaker defines a structure which returns a new db connection for
// use in creating new sql db instances for db ops.
type dBMaker struct {
	config Config
	log    metrics.Metrics
}

// New returns a new instance of a sqlx.DB connected to the db with the provided
// credentials pulled from the host environment.
func (dl dBMaker) New() (*sqlx.DB, error) {
	if dl.config.DBIP == "" {
		dl.config.DBIP = "0.0.0.0"
	}

	addr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dl.config.User, dl.config.UserPassword, dl.config.DBIP, dl.config.DBPort, dl.config.DBName)
	db, err := sqlx.Connect(dl.config.DBDriver, addr)
	if err != nil {
		dl.log.Emit(stdout.Error("Failed to connect to SQLServer: %+q", err).WithFields(metrics.Fields{
			"ip":     dl.config.DBIP,
			"port":   dl.config.DBPort,
			"dbName": dl.config.DBName,
			"driver": dl.config.DBDriver,
		}))

		return nil, err
	}

	return db, nil
}

//===============================================================================================================

// DB defines an interface which exposes a method to return a new
// underline sql.Database.
type DB interface {
	New() (*sqlx.DB, error)
}

// NewDB returns a new instance of a DB.
func NewDB(config Config, metrics metrics.Metrics) DB {
	return &dBMaker{
		config: config,
		log:    metrics,
	}
}

//===============================================================================================================

// SQL defines an struct which implements the db.Provider which allows us
// execute CRUD ops.
type SQL struct {
	d      DB
	inited bool
	l      metrics.Metrics
	tables []tables.TableMigration
}

// New returns a new instance of SQL.
func New(s metrics.Metrics, d DB, ts ...tables.TableMigration) *SQL {
	return &SQL{
		d:      d,
		l:      s,
		tables: ts,
	}
}

// migrate takes the individual query supplied and attempts to
// execute them returning any error found.
func (sq *SQL) migrate() error {
	if sq.d == nil {
		return nil
	}

	if sq.inited {
		return nil
	}

	dbi, err := sq.d.New()
	if err != nil {
		return err
	}

	defer dbi.Close()

	for _, table := range sq.tables {
		sq.l.Emit(stdout.Info("Executing Migration").WithFields(metrics.Fields{
			"query": table.String(),
			"table": table.TableName,
		}))

		if _, err := dbi.Exec(table.String()); err != nil {
			sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{"query": table.String(), "table": table.TableName}))
			return err
		}
	}

	sq.inited = true

	return nil
}

// Save takes the giving table name with the giving fields and attempts to save this giving
// data appropriately into the giving db.
func (sq *SQL) Save(identity db.TableIdentity, table db.TableFields) error {
	defer sq.l.Emit(stdout.Info("Save to DB").With("table", identity.Table()).Trace("db.Save").End())

	if err := sq.migrate(); err != nil {
		return err
	}

	db, err := sq.d.New()
	if err != nil {
		return err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	fields := table.Fields()
	fieldNames := fieldNames(fields)
	values := fieldValues(fieldNames, fields)

	fieldNames = append(fieldNames, "created_at")
	fieldNames = append(fieldNames, "updated_at")

	values = append(values, time.Now().UTC())
	values = append(values, time.Now().UTC())

	query := fmt.Sprintf(insertTemplate, identity.Table(), fieldNameMarkers(fieldNames), fieldMarkers(len(fieldNames)))
	sq.l.Emit(stdout.Info("DB:Query").With("query", query))

	if _, err := db.Exec(query, values...); err != nil {
		sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": identity.Table(),
		}))
		return err
	}

	return tx.Commit()
}

// Update takes the giving table name with the giving fields and attempts to update this giving
// data appropriately into the giving db.
// index - defines the string which should identify the key to be retrieved from the fields to target the
// data to be updated in the db.
func (sq *SQL) Update(identity db.TableIdentity, table db.TableFields, index string, indexValue interface{}) error {
	defer sq.l.Emit(stdout.Info("Update to DB").With("table", identity.Table()).Trace("db.Update").End())

	if err := sq.migrate(); err != nil {
		return err
	}

	db, err := sq.d.New()
	if err != nil {
		return err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	tableFields := table.Fields()
	tableFields["updated_at"] = time.Now().UTC()

	indexValueString, err := printLiteral(indexValue)
	if err != nil {
		sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{
			"err":   err,
			"table": identity.Table(),
		}))
		return err
	}

	sets, err := setValues(tableFields)
	if err != nil {
		sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{
			"err":   err,
			"table": identity.Table(),
		}))
		return err
	}

	query := fmt.Sprintf(updateTemplate, identity.Table(), sets, index, indexValueString)
	sq.l.Emit(stdout.Info("DB:Query").With("query", query))

	if _, err := db.Exec(query); err != nil {
		sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": identity.Table(),
		}))
		return err
	}

	return tx.Commit()
}

// GetAllPerPage retrieves the giving data from the specific db with the specific index and value.
func (sq *SQL) GetAllPerPage(table db.TableIdentity, order string, orderBy string, page int, responsePerPage int) ([]map[string]interface{}, int, error) {
	defer sq.l.Emit(stdout.Info("Retrieve all records from DB").With("table", table.Table()).WithFields(metrics.Fields{
		"page":            page,
		"order":           order,
		"orderBy":         orderBy,
		"responsePerPage": responsePerPage,
	}).Trace("db.GetAllPerPage").End())

	if err := sq.migrate(); err != nil {
		return nil, -1, err
	}

	db, err := sq.d.New()
	if err != nil {
		return nil, -1, err
	}

	defer db.Close()

	if page <= 0 && responsePerPage <= 0 {
		records, err := sq.GetAll(table, order, orderBy)
		return records, len(records), err
	}

	// Get total number of records.
	totalRecords, err := sq.Count(table)
	if err != nil {
		return nil, -1, err
	}

	var totalWanted, indexToStart int

	if page <= 1 && responsePerPage > 0 {
		totalWanted = responsePerPage
		indexToStart = 0
	} else {
		totalWanted = responsePerPage * page
		indexToStart = totalWanted / 2

		if page > 1 {
			indexToStart++
		}
	}

	sq.l.Emit(stdout.Info("DB:Query:GetAllPerPage").WithFields(metrics.Fields{
		"starting_index":       indexToStart,
		"total_records_wanted": totalWanted,
		"order":                order,
		"page":                 page,
		"responsePerPage":      responsePerPage,
	}))

	// If we are passed the total, just return nil records and total without error.
	if indexToStart > totalRecords {
		return nil, totalRecords, nil
	}

	query := fmt.Sprintf(selectLimitedTemplate, table.Table(), orderBy, order, totalWanted, indexToStart)
	sq.l.Emit(stdout.Info("DB:Query:GetAllPerPage").With("query", query))

	rows, err := db.Queryx(query)
	if err != nil {
		sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))
		return nil, -1, err
	}

	var fields []map[string]interface{}

	for rows.Next() {
		mo := make(map[string]interface{})

		if err := rows.MapScan(mo); err != nil {
			sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{
				"err":   err,
				"query": query,
				"table": table.Table(),
			}))

			return nil, -1, err
		}

		fields = append(fields, naturalizeMap(mo))
	}

	return fields, totalRecords, nil
}

// GetAllPerPageBy retrieves the giving data from the specific db with the specific index and value.
func (sq *SQL) GetAllPerPageBy(table db.TableIdentity, order string, orderBy string, page int, responsePerPage int, mx func(*sqlx.Rows) error) (int, error) {
	defer sq.l.Emit(stdout.Info("Retrieve all records from DB").With("table", table.Table()).WithFields(metrics.Fields{
		"order":           order,
		"page":            page,
		"responsePerPage": responsePerPage,
	}).Trace("db.GetAllPerPageBy").End())

	if err := sq.migrate(); err != nil {
		return -1, err
	}

	db, err := sq.d.New()
	if err != nil {
		return -1, err
	}

	defer db.Close()

	if page <= 0 && responsePerPage <= 0 {
		records, err := sq.GetAll(table, order, orderBy)
		return  len(records), err
	}

	// Get total number of records.
	totalRecords, err := sq.Count(table)
	if err != nil {
		return -1, err
	}

	switch strings.ToLower(order) {
	case "asc":
		order = "ASC"
	case "dsc", "desc":
		order = "DESC"
	default:
		order = "ASC"
	}

	var totalWanted, indexToStart int

	if page <= 1 && responsePerPage > 0 {
		totalWanted = responsePerPage
		indexToStart = 0
	} else {
		totalWanted = responsePerPage * page
		indexToStart = totalWanted / 2

		if page > 1 {
			indexToStart++
		}
	}

	sq.l.Emit(stdout.Info("DB:Query:GetAllPerPageBy").WithFields(metrics.Fields{
		"starting_index":       indexToStart,
		"total_records_wanted": totalWanted,
		"order":                order,
		"page":                 page,
		"responsePerPage":      responsePerPage,
	}))

	// If we are passed the total, just return nil records and total without error.
	if indexToStart > totalRecords {
		return totalRecords, nil
	}

	query := fmt.Sprintf(selectLimitedTemplate, table.Table(), orderBy, order, totalWanted, indexToStart)

	sq.l.Emit(stdout.Info("DB:Query:GetAllPerPageBy").With("query", query))

	rows, err := db.Queryx(query)
	if err != nil {
		sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))

		return -1, err
	}

	// var fields []map[string]interface{}

	// for rows.Next() {
	// 	mo := make(map[string]interface{})

	// 	if err := rows.MapScan(mo); err != nil {
	// 		sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{
	// 			"err":   err,
	// 			"query": query,
	// 			"table": table.Table(),
	// 		}))

	// 		return nil, -1, err
	// 	}

	// 	fields = append(fields, naturalizeMap(mo))
	// }

	if err := mx(rows); err != nil {
		return -1, err
	}

	return totalRecords, nil
}

// GetAll retrieves the giving data from the specific db with the specific index and value.
func (sq *SQL) GetAll(table db.TableIdentity, order string, orderBy string) ([]map[string]interface{}, error) {
	defer sq.l.Emit(stdout.Info("Retrieve all records from DB").With("table", table.Table()).Trace("db.GetAll").End())

	if err := sq.migrate(); err != nil {
		return nil, err
	}

	db, err := sq.d.New()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	switch strings.ToLower(order) {
	case "asc":
		order = "ASC"
	case "dsc", "desc":
		order = "DESC"
	default:
		order = "ASC"
	}

	var fields []map[string]interface{}

	query := fmt.Sprintf(selectAllTemplate, table.Table(), orderBy, order)
	sq.l.Emit(stdout.Info("DB:Query:GetAll").With("query", query))

	rows, err := db.Queryx(query)
	if err != nil {
		sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))
		return nil, err
	}

	for rows.Next() {
		mo := make(map[string]interface{})
		if err := rows.MapScan(mo); err != nil {
			sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{
				"err":   err,
				"query": query,
				"table": table.Table(),
			}))
			return nil, err
		}

		fields = append(fields, naturalizeMap(mo))
	}

	return fields, nil
}

// GetAllBy retrieves the giving data from the specific db with the specific index and value.
func (sq *SQL) GetAllBy(table db.TableIdentity, order string, orderBy string, mx func(*sqlx.Rows) error) error {
	defer sq.l.Emit(stdout.Info("Retrieve all records from DB").With("table", table.Table()).Trace("db.GetAllBy").End())

	if err := sq.migrate(); err != nil {
		return nil
	}

	db, err := sq.d.New()
	if err != nil {
		return nil
	}

	defer db.Close()

	switch strings.ToLower(order) {
	case "asc":
		order = "ASC"
	case "dsc", "desc":
		order = "DESC"
	default:
		order = "ASC"
	}

	// var fields []map[string]interface{}

	query := fmt.Sprintf(selectAllTemplate, table.Table(), orderBy, order)

	sq.l.Emit(stdout.Info("DB:Query:GetAll").With("query", query))

	rows, err := db.Queryx(query)
	if err != nil {
		sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))

		return err
	}

	if err := mx(rows); err != nil {
		return err
	}

	return nil
}

// Get retrieves the giving data from the specific db with the specific index and value.
func (sq *SQL) Get(table db.TableIdentity, consumer db.TableConsumer, index string, indexValue interface{}) error {
	defer sq.l.Emit(stdout.Info("Get record from DB").WithFields(metrics.Fields{
		"table":      table.Table(),
		"index":      index,
		"indexValue": indexValue,
	}).Trace("db.Get").End())

	if err := sq.migrate(); err != nil {
		return err
	}

	db, err := sq.d.New()
	if err != nil {
		return err
	}

	defer db.Close()

	indexValueString, err := printLiteral(indexValue)
	if err != nil {
		sq.l.Emit(stdout.Error("DB:Query: %+q", err).WithFields(metrics.Fields{
			"err":   err,
			"table": table.Table(),
		}))
		return err
	}

	query := fmt.Sprintf(selectItemTemplate, table.Table(), index, indexValueString)
	sq.l.Emit(stdout.Info("DB:Query").With("query", query))

	row := db.QueryRowx(query)
	if err := row.Err(); err != nil {
		sq.l.Emit(stdout.Error("DB:Query: %+q", err).WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))
		return err
	}

	mo := make(map[string]interface{})

	if err := row.MapScan(mo); err != nil {
		sq.l.Emit(stdout.Error("DB:Query: %+q", err).WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))

		return err
	}

	sq.l.Emit(stdout.Debug("Consumer:Get:Fields").WithFields(metrics.Fields{
		"table":    table.Table(),
		"response": mo,
	}))

	if err := consumer.WithFields(naturalizeMap(mo)); err != nil {
		sq.l.Emit(stdout.Error("Consumer:WithFields: %+q", err).WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))

		return err
	}

	return nil
}

// GetBy retrieves the giving data from the specific db with the specific index and value.
func (sq *SQL) GetBy(table db.TableIdentity, consumer func(*sqlx.Row) error, index string, indexValue interface{}) error {
	defer sq.l.Emit(stdout.Info("Get record from DB").WithFields(metrics.Fields{
		"table":      table.Table(),
		"index":      index,
		"indexValue": indexValue,
	}).Trace("db.GetBy").End())

	if err := sq.migrate(); err != nil {
		return err
	}

	db, err := sq.d.New()
	if err != nil {
		return err
	}

	defer db.Close()

	indexValueString, err := printLiteral(indexValue)
	if err != nil {
		sq.l.Emit(stdout.Error("DB:Query: %+q", err).WithFields(metrics.Fields{
			"err":   err,
			"table": table.Table(),
		}))
		return err
	}

	query := fmt.Sprintf(selectItemTemplate, table.Table(), index, indexValueString)
	sq.l.Emit(stdout.Info("DB:Query").With("query", query))

	row := db.QueryRowx(query)
	if err := row.Err(); err != nil {
		sq.l.Emit(stdout.Error("DB:Query: %+q", err).WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))
		return err
	}

	if err := consumer(row); err != nil {
		sq.l.Emit(stdout.Error("DB:Query: %+q", err).WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))

		return err
	}

	sq.l.Emit(stdout.Debug("Consumer:Get:Fields").WithFields(metrics.Fields{
		"table":    table.Table(),
		"index":      index,
		"indexValue": indexValue,
	}))

	return nil
}

// Count retrieves the total number of records from the specific table from the db.
func (sq *SQL) Count(table db.TableIdentity) (int, error) {
	defer sq.l.Emit(stdout.Info("Count record from DB").WithFields(metrics.Fields{
		"table": table.Table(),
	}).Trace("db.Get").End())

	if err := sq.migrate(); err != nil {
		return 0, err
	}

	db, err := sq.d.New()
	if err != nil {
		return 0, err
	}

	defer db.Close()

	var records int

	query := fmt.Sprintf(countTemplate, table.Table())
	sq.l.Emit(stdout.Info("DB:Query").With("query", query))

	if err := db.Get(&records, query); err != nil {
		sq.l.Emit(stdout.Error("DB:Query").WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))
		return 0, err
	}

	return records, nil
}

// Delete removes the giving data from the specific db with the specific index and value.
func (sq *SQL) Delete(table db.TableIdentity, index string, indexValue interface{}) error {
	defer sq.l.Emit(stdout.Info("Delete record from DB").WithFields(metrics.Fields{
		"table":      table.Table(),
		"index":      index,
		"indexValue": indexValue,
	}).Trace("db.GetAll").End())

	if err := sq.migrate(); err != nil {
		return err
	}

	db, err := sq.d.New()
	if err != nil {
		return err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	indexValueString, err := printLiteral(indexValue)
	if err != nil {
		sq.l.Emit(stdout.Error("DB:Query: %+q", err).WithFields(metrics.Fields{
			"err":   err,
			"table": table.Table(),
		}))
		return err
	}

	query := fmt.Sprintf(deleteTemplate, table.Table(), index, indexValueString)
	sq.l.Emit(stdout.Info("DB:Query").With("query", query))

	if _, err := db.Exec(query); err != nil {
		sq.l.Emit(stdout.Error(err).WithFields(metrics.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))
		return err
	}

	return tx.Commit()
}

// FieldMarkers returns a (?,...,>) string which represents
// all filedNames extrated from the provided TableField.
func fieldMarkers(total int) string {
	var markers []string

	for i := 0; i < total; i++ {
		markers = append(markers, "?")
	}

	return "(" + strings.Join(markers, ",") + ")"
}

// fieldNameMarkers returns a (fieldName,...,fieldName) string which represents
// all filedNames extrated from the provided TableField.
func fieldNameMarkers(fields []string) string {
	return "(" + strings.Join(fields, ", ") + ")"
}

// fieldValues returns a (fieldName,...,fieldName) string which represents
// all filedNames extrated from the provided TableField.
func fieldValues(names []string, fields map[string]interface{}) []interface{} {
	var vals []interface{}

	for _, name := range names {
		vals = append(vals, fields[name])
	}

	return vals
}

func setValues(fields map[string]interface{}) (string, error) {
	var vals []string

	for name, val := range fields {
		rv, err := printLiteral(val)
		if err != nil {
			return "", err
		}

		vals = append(vals, fmt.Sprintf("%s=%s", name, rv))
	}

	return strings.Join(vals, ","), nil
}

// naturalizeMap returns a new map where all values of []bytes are converted to strings
func naturalizeMap(fields map[string]interface{}) map[string]interface{} {
	nz := map[string]interface{}{}

	for key, val := range fields {
		switch rv := val.(type) {
		case []byte:
			nz[key] = string(rv)
			continue
		default:
			nz[key] = val
			continue
		}
	}

	return nz
}

// fieldNamesFromMap returns a (fieldName,...,fieldName) string which represents
// all filedNames extrated from the provided TableField.
func fieldNamesFromMap(fields map[string]interface{}) []string {
	var names []string

	for key := range fields {
		names = append(names, key)
	}

	return names
}

// fieldNames returns a (fieldName,...,fieldName) string which represents
// all filedNames extrated from the provided TableField.
func fieldNames(fields map[string]interface{}) []string {
	var names []string

	for key := range fields {
		names = append(names, key)
	}

	return names
}

// printLiteral attempts to provide a function to allow us easily convert
// simple values like int, float, uint to string for use in queries.
func printLiteral(item interface{}) (string, error) {
	switch rl := item.(type) {
	case int, int64, int32:
		return strconv.Itoa(rl.(int)), nil
	case float32, float64:
		return strconv.FormatFloat(rl.(float64), 'f', 2, 64), nil
	case string:
		return strconv.Quote(rl), nil
	case []byte:
		return strconv.Quote(string(rl)), nil
	case time.Time:
		return strconv.Quote(rl.String()), nil
	case byte:
		return strconv.QuoteRune(rune(rl)), nil
	default:
		return "", errors.New("Not basic type")
	}
}
