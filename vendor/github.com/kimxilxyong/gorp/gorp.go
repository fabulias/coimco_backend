// Copyright 2012 James Cooper. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package gorp provides a simple way to marshal Go structs to and from
// SQL databases.  It uses the database/sql package, and should work with any
// compliant database/sql driver.
//
// Gorp with Indexes, a fork of gorp by Kim Il
// Source code and project home:
// https://github.com/kimxilxyong/gorp
//
// Original:
// Source code and project home:
// https://github.com/go-gorp/gorp
//
// History:
// 2015.05.16 Forked and initial index support
// 2015.05.25 Added named parameter support for the new field tags
// 2015.05.29 Added support for PostgreSQL

package gorp

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Oracle String (empty string is null)
type OracleString struct {
	sql.NullString
}

// Scan implements the Scanner interface.
func (os *OracleString) Scan(value interface{}) error {
	if value == nil {
		os.String, os.Valid = "", false
		return nil
	}
	os.Valid = true
	return os.NullString.Scan(value)
}

// Value implements the driver Valuer interface.
func (os OracleString) Value() (driver.Value, error) {
	if !os.Valid || os.String == "" {
		return nil, nil
	}
	return os.String, nil
}

// SqlTyper is a type that returns its database type.  Most of the
// time, the type can just use "database/sql/driver".Valuer; but when
// it returns nil for its empty value, it needs to implement SqlTyper
// to have its column type detected properly during table creation.
type SqlTyper interface {
	SqlType() driver.Valuer
}

// A nullable Time value
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	switch t := value.(type) {
	case time.Time:
		nt.Time, nt.Valid = t, true
	case []byte:
		nt.Valid = false
		for _, dtfmt := range []string{
			"2006-01-02 15:04:05.999999999",
			"2006-01-02T15:04:05.999999999",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
			"2006-01-02 15:04",
			"2006-01-02T15:04",
			"2006-01-02",
			"2006-01-02 15:04:05-07:00",
		} {
			var err error
			if nt.Time, err = time.Parse(dtfmt, string(t)); err == nil {
				nt.Valid = true
				break
			}
		}
	}
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

var zeroVal reflect.Value
var versFieldConst = "[gorp_ver_field]"

// OptimisticLockError is returned by Update() or Delete() if the
// struct being modified has a Version field and the value is not equal to
// the current value in the database
type OptimisticLockError struct {
	// Table name where the lock error occurred
	TableName string

	// Primary key values of the row being updated/deleted
	Keys []interface{}

	// true if a row was found with those keys, indicating the
	// LocalVersion is stale.  false if no value was found with those
	// keys, suggesting the row has been deleted since loaded, or
	// was never inserted to begin with
	RowExists bool

	// Version value on the struct passed to Update/Delete. This value is
	// out of sync with the database.
	LocalVersion int64
}

// Error returns a description of the cause of the lock error
func (e OptimisticLockError) Error() string {
	if e.RowExists {
		return fmt.Sprintf("gorp: OptimisticLockError table=%s keys=%v out of date version=%d", e.TableName, e.Keys, e.LocalVersion)
	}

	return fmt.Sprintf("gorp: OptimisticLockError no row found for table=%s keys=%v", e.TableName, e.Keys)
}

// The TypeConverter interface provides a way to map a value of one
// type to another type when persisting to, or loading from, a database.
//
// Example use cases: Implement type converter to convert bool types to "y"/"n" strings,
// or serialize a struct member as a JSON blob.
type TypeConverter interface {
	// ToDb converts val to another type. Called before INSERT/UPDATE operations
	ToDb(val interface{}) (interface{}, error)

	// FromDb returns a CustomScanner appropriate for this type. This will be used
	// to hold values returned from SELECT queries.
	//
	// In particular the CustomScanner returned should implement a Binder
	// function appropriate for the Go type you wish to convert the db value to
	//
	// If bool==false, then no custom scanner will be used for this field.
	FromDb(target interface{}) (CustomScanner, bool)
}

// CustomScanner binds a database column value to a Go type
type CustomScanner struct {
	// After a row is scanned, Holder will contain the value from the database column.
	// Initialize the CustomScanner with the concrete Go type you wish the database
	// driver to scan the raw column into.
	Holder interface{}
	// Target typically holds a pointer to the target struct field to bind the Holder
	// value to.
	Target interface{}
	// Binder is a custom function that converts the holder value to the target type
	// and sets target accordingly.  This function should return error if a problem
	// occurs converting the holder to the target.
	Binder func(holder interface{}, target interface{}) error
}

// Bind is called automatically by gorp after Scan()
func (me CustomScanner) Bind() error {
	return me.Binder(me.Holder, me.Target)
}

// DbMap is the root gorp mapping object. Create one of these for each
// database schema you wish to map.  Each DbMap contains a list of
// mapped tables.
//
// Example:
//
//     dialect := gorp.MySQLDialect{"InnoDB", "UTF8"}
//     dbmap := &gorp.DbMap{Db: db, Dialect: dialect}
//
type DbMap struct {
	// Db handle to use with this map
	Db *sql.DB

	// Dialect implementation to use with this map
	Dialect Dialect

	TypeConverter TypeConverter

	tables    []*TableMap
	logger    GorpLogger
	logPrefix string

	DebugLevel        int
	LastOpInfo        CRUDInfo // info about the last operation on this database
	CheckAffectedRows bool     // if true an error is raised if affected rows was 0
}

// TableMap represents a mapping between a Go struct and a database table
// Use dbmap.AddTable() or dbmap.AddTableWithName() to create these
type TableMap struct {
	// Name of database table.
	TableName      string
	SchemaName     string
	gotype         reflect.Type
	Columns        []*ColumnMap
	Indexes        []*IndexMap    // list of indexes for this table
	Relations      []*RelationMap // list of detail/child tables for this table
	keys           []*ColumnMap
	uniqueTogether [][]string
	version        *ColumnMap
	insertPlan     bindPlan
	updatePlan     bindPlan
	deletePlan     bindPlan
	getPlan        bindPlan
	dbmap          *DbMap
}

func (t TableMap) String() string {
	var s string
	s = "TableName: " + t.TableName + "\n"
	s = s + "Type: " + t.gotype.Name() + "\n"
	for _, c := range t.Columns {
		s = s + "Column: " + c.ColumnName + "\n"
	}
	return s
}

// RelationMap represents a mapping between a master table and a detail table
// Use tablemap.AddRelation() or field tag `db:"relation:<foreignkey field in detail table>"`to create these
// Example:	Comments  []*Comment `db:"relation:PostId"`
type RelationMap struct {
	DetailTable         *TableMap
	ForeignKeyFieldName string
	DetailTableType     interface{}
	MasterFieldName     string
	Limit               uint64 // Limits the number of rows a child query returns
	Offset              uint64 // Starting at row offset when queriyng a child table
}

func (r RelationMap) String() string {
	var s string
	s = "TableMap: " + r.DetailTable.String() + "\n"
	s = s + "ForeignKeyFieldName: " + r.ForeignKeyFieldName + "\n"
	s = s + "DetailTableType: " + reflect.TypeOf(r.DetailTableType).Name() + "\n"
	s = s + "MasterFieldName: " + r.MasterFieldName
	return s
}

// CRUDInfo contains info about a CRUD operation
type CRUDInfo struct {
	Type                CRUDType
	BindPlanUsed        *bindPlan
	RowCount            int64
	ChildUpdateRowCount int64
	ChildInsertRowCount int64
}

// Resets all field to its zero value
func (i *CRUDInfo) Reset() {
	i.Type = Unknown
	i.BindPlanUsed = nil
	i.RowCount = 0
	i.ChildUpdateRowCount = 0
	i.ChildInsertRowCount = 0
}

type CRUDType int

const (
	Unknown CRUDType = iota
	Insert
	Select
	Update
	Delete
)

func (t CRUDType) String() string {
	var s string
	if t&Insert == Insert {
		s = "Insert"
	} else if t&Select == Select {
		s = "Select"
	} else if t&Update == Update {
		s = "Update"
	} else if t&Delete == Delete {
		s = "Delete"
	} else {
		s = "Unknown"
	}
	return s
}

// ResetSql removes cached insert/update/select/delete SQL strings
// associated with this TableMap.  Call this if you've modified
// any column names or the table name itself.
func (t *TableMap) ResetSql() {
	t.insertPlan = bindPlan{}
	t.updatePlan = bindPlan{}
	t.deletePlan = bindPlan{}
	t.getPlan = bindPlan{}
}

// SetKeys lets you specify the fields on a struct that map to primary
// key columns on the table.  If isAutoIncr is set, result.LastInsertId()
// will be used after INSERT to bind the generated id to the Go struct.
//
// Automatically calls ResetSql() to ensure SQL statements are regenerated.
//
// Panics if isAutoIncr is true, and fieldNames length != 1
//
func (t *TableMap) SetKeys(isAutoIncr bool, fieldNames ...string) *TableMap {
	if isAutoIncr && len(fieldNames) != 1 {
		panic(fmt.Sprintf(
			"gorp: SetKeys: fieldNames length must be 1 if key is auto-increment. (Saw %v fieldNames)",
			len(fieldNames)))
	}
	t.keys = make([]*ColumnMap, 0)
	for _, name := range fieldNames {
		colmap := t.ColMap(name)
		colmap.isPK = true
		colmap.isAutoIncr = isAutoIncr
		t.keys = append(t.keys, colmap)
	}
	t.ResetSql()

	return t
}

// SetUniqueTogether lets you specify uniqueness constraints across multiple
// columns on the table. Each call adds an additional constraint for the
// specified columns.
//
// Automatically calls ResetSql() to ensure SQL statements are regenerated.
//
// Panics if fieldNames length < 2.
//
func (t *TableMap) SetUniqueTogether(fieldNames ...string) *TableMap {
	if len(fieldNames) < 2 {
		panic(fmt.Sprintf(
			"gorp: SetUniqueTogether: must provide at least two fieldNames to set uniqueness constraint."))
	}

	columns := make([]string, 0)
	for _, name := range fieldNames {
		columns = append(columns, name)
	}
	t.uniqueTogether = append(t.uniqueTogether, columns)
	t.ResetSql()

	return t
}

// ColMap returns the ColumnMap pointer matching the given struct field
// name.  It panics if the struct does not contain a field matching this
// name.
func (t *TableMap) ColMap(field string) *ColumnMap {
	col := colMapOrNil(t, field)
	if col == nil {
		e := fmt.Sprintf("No ColumnMap in table %s type %s with field %s",
			t.TableName, t.gotype.Name(), field)

		panic(e)
	}
	return col
}

func colMapOrNil(t *TableMap, field string) *ColumnMap {
	for _, col := range t.Columns {
		if strings.ToLower(col.fieldName) == strings.ToLower(field) || strings.ToLower(col.ColumnName) == strings.ToLower(field) {
			// Ignore ignored Columns, "-" is the indentifier for columns which should be ignored
			//if col.ColumnName != "-" { // old code, depending on a hardcoded string is avoided now
			if col.Transient == false {
				return col
			}
		}
	}
	return nil
}

// SetVersionCol sets the column to use as the Version field.  By default
// the "Version" field is used.  Returns the column found, or panics
// if the struct does not contain a field matching this name.
//
// Automatically calls ResetSql() to ensure SQL statements are regenerated.
func (t *TableMap) SetVersionCol(field string) *ColumnMap {
	c := t.ColMap(field)
	t.version = c
	t.ResetSql()
	return c
}

// SqlForCreateTable gets a sequence of SQL commands that will create
// the specified table and any associated schema
func (t *TableMap) SqlForCreate(ifNotExists bool) string {
	s := bytes.Buffer{}
	dialect := t.dbmap.Dialect

	if strings.TrimSpace(t.SchemaName) != "" {
		schemaCreate := "create schema"
		if ifNotExists {
			s.WriteString(dialect.IfSchemaNotExists(schemaCreate, t.SchemaName))
		} else {
			s.WriteString(schemaCreate)
		}
		s.WriteString(fmt.Sprintf(" %s;", t.SchemaName))
	}

	tableCreate := "create table"
	if ifNotExists {
		s.WriteString(dialect.IfTableNotExists(tableCreate, t.SchemaName, t.TableName))
	} else {
		s.WriteString(tableCreate)
	}
	s.WriteString(fmt.Sprintf(" %s (", dialect.QuotedTableForQuery(t.SchemaName, t.TableName)))

	x := 0
	for _, col := range t.Columns {
		if !col.Transient {
			if x > 0 {
				s.WriteString(", ")
			}
			stype := dialect.ToSqlType(col.gotype, col.MaxSize, col.isAutoIncr)
			s.WriteString(fmt.Sprintf("%s %s", dialect.QuoteField(col.ColumnName), stype))

			if col.isPK || col.isNotNull {
				s.WriteString(" not null")
			}
			if col.isPK && len(t.keys) == 1 {
				s.WriteString(" primary key")
			}
			if col.Unique {
				s.WriteString(" unique")
			}
			if col.isAutoIncr {
				s.WriteString(fmt.Sprintf(" %s", dialect.AutoIncrStr()))
			}

			x++
		}
	}
	if len(t.keys) > 1 {
		s.WriteString(", primary key (")
		for x := range t.keys {
			if x > 0 {
				s.WriteString(", ")
			}
			s.WriteString(dialect.QuoteField(t.keys[x].ColumnName))
		}
		s.WriteString(")")
	}
	if len(t.uniqueTogether) > 0 {
		for _, columns := range t.uniqueTogether {
			s.WriteString(", unique (")
			for i, column := range columns {
				if i > 0 {
					s.WriteString(", ")
				}
				s.WriteString(dialect.QuoteField(column))
			}
			s.WriteString(")")
		}
	}
	s.WriteString(") ")
	s.WriteString(dialect.CreateTableSuffix())
	s.WriteString(dialect.QuerySuffix())
	return s.String()
}

type bindPlan struct {
	query             string
	argFields         []string
	keyFields         []string
	versField         string
	autoIncrIdx       int
	autoIncrFieldName string
}

func (plan bindPlan) createBindInstance(elem reflect.Value, conv TypeConverter) (bindInstance, error) {
	bi := bindInstance{query: plan.query, autoIncrIdx: plan.autoIncrIdx, autoIncrFieldName: plan.autoIncrFieldName, versField: plan.versField}
	if plan.versField != "" {
		bi.existingVersion = elem.FieldByName(plan.versField).Int()
	}

	var err error

	for i := 0; i < len(plan.argFields); i++ {
		k := plan.argFields[i]
		if k == versFieldConst {
			newVer := bi.existingVersion + 1
			bi.args = append(bi.args, newVer)
			if bi.existingVersion == 0 {
				elem.FieldByName(plan.versField).SetInt(int64(newVer))
			}
		} else {
			val := elem.FieldByName(k).Interface()
			if conv != nil {
				val, err = conv.ToDb(val)
				if err != nil {
					return bindInstance{}, err
				}
			}
			bi.args = append(bi.args, val)
		}
	}

	for i := 0; i < len(plan.keyFields); i++ {
		k := plan.keyFields[i]
		val := elem.FieldByName(k).Interface()
		if conv != nil {
			val, err = conv.ToDb(val)
			if err != nil {
				return bindInstance{}, err
			}
		}
		bi.keys = append(bi.keys, val)
	}

	return bi, nil
}

type bindInstance struct {
	query             string
	args              []interface{}
	keys              []interface{}
	existingVersion   int64
	versField         string
	autoIncrIdx       int
	autoIncrFieldName string
}

func (t *TableMap) bindInsert(elem reflect.Value) (bindInstance, error) {
	plan := t.insertPlan
	if plan.query == "" {
		plan.autoIncrIdx = -1

		s := bytes.Buffer{}
		s2 := bytes.Buffer{}
		s.WriteString(fmt.Sprintf("insert into %s (", t.dbmap.Dialect.QuotedTableForQuery(t.SchemaName, t.TableName)))

		x := 0
		first := true
		for y := range t.Columns {
			col := t.Columns[y]
			if !(col.isAutoIncr && t.dbmap.Dialect.AutoIncrBindValue() == "") {
				if !col.Transient {
					if !first {
						s.WriteString(",")
						s2.WriteString(",")
					}
					s.WriteString(t.dbmap.Dialect.QuoteField(col.ColumnName))

					if col.isAutoIncr {
						s2.WriteString(t.dbmap.Dialect.AutoIncrBindValue())
						plan.autoIncrIdx = y
						plan.autoIncrFieldName = col.fieldName
					} else {
						if col.DefaultValue == "" {
							s2.WriteString(t.dbmap.Dialect.BindVar(x))
							if col == t.version {
								plan.versField = col.fieldName
								plan.argFields = append(plan.argFields, versFieldConst)
							} else {
								plan.argFields = append(plan.argFields, col.fieldName)
							}
							x++
						} else {

							// Check if this column is a NOT NULL
							if err := checkForNotNull(elem, col, t); err != nil {
								return bindInstance{}, err
							}
							plan.argFields = append(plan.argFields, col.fieldName)
							s2.WriteString(col.DefaultValue)
						}
					}
					first = false
				}
			} else {
				plan.autoIncrIdx = y
				plan.autoIncrFieldName = col.fieldName
			}
		}
		s.WriteString(") values (")
		s.WriteString(s2.String())
		s.WriteString(")")
		if plan.autoIncrIdx > -1 {
			s.WriteString(t.dbmap.Dialect.AutoIncrInsertSuffix(t.Columns[plan.autoIncrIdx]))
		}
		s.WriteString(t.dbmap.Dialect.QuerySuffix())

		plan.query = s.String()
		t.insertPlan = plan
	}

	return plan.createBindInstance(elem, t.dbmap.TypeConverter)
}

func (t *TableMap) bindUpdate(elem reflect.Value) (bindInstance, error) {
	plan := t.updatePlan
	if plan.query == "" {

		s := bytes.Buffer{}
		s.WriteString(fmt.Sprintf("update %s set ", t.dbmap.Dialect.QuotedTableForQuery(t.SchemaName, t.TableName)))
		x := 0

		for y := range t.Columns {
			col := t.Columns[y]
			if !col.isAutoIncr && !col.Transient {
				if x > 0 {
					s.WriteString(", ")
				}
				s.WriteString(t.dbmap.Dialect.QuoteField(col.ColumnName))
				s.WriteString("=")
				s.WriteString(t.dbmap.Dialect.BindVar(x))

				if col == t.version {
					plan.versField = col.fieldName
					plan.argFields = append(plan.argFields, versFieldConst)
				} else {
					// Check if this column is a NOT NULL
					if err := checkForNotNull(elem, col, t); err != nil {
						return bindInstance{}, err
					}
					plan.argFields = append(plan.argFields, col.fieldName)
				}
				x++
			}
		}

		s.WriteString(" where ")
		for y := range t.keys {
			col := t.keys[y]
			if y > 0 {
				s.WriteString(" and ")
			}
			s.WriteString(t.dbmap.Dialect.QuoteField(col.ColumnName))
			s.WriteString("=")
			s.WriteString(t.dbmap.Dialect.BindVar(x))

			plan.argFields = append(plan.argFields, col.fieldName)
			plan.keyFields = append(plan.keyFields, col.fieldName)
			x++
		}
		if plan.versField != "" {
			s.WriteString(" and ")
			s.WriteString(t.dbmap.Dialect.QuoteField(t.version.ColumnName))
			s.WriteString("=")
			s.WriteString(t.dbmap.Dialect.BindVar(x))
			plan.argFields = append(plan.argFields, plan.versField)
		}
		s.WriteString(t.dbmap.Dialect.QuerySuffix())

		plan.query = s.String()
		t.updatePlan = plan
	}

	return plan.createBindInstance(elem, t.dbmap.TypeConverter)
}

func (t *TableMap) bindDelete(elem reflect.Value) (bindInstance, error) {
	plan := t.deletePlan
	if plan.query == "" {

		s := bytes.Buffer{}
		s.WriteString(fmt.Sprintf("delete from %s", t.dbmap.Dialect.QuotedTableForQuery(t.SchemaName, t.TableName)))

		for y := range t.Columns {
			col := t.Columns[y]
			if !col.Transient {
				if col == t.version {
					plan.versField = col.fieldName
				}
			}
		}

		s.WriteString(" where ")
		for x := range t.keys {
			k := t.keys[x]
			if x > 0 {
				s.WriteString(" and ")
			}
			s.WriteString(t.dbmap.Dialect.QuoteField(k.ColumnName))
			s.WriteString("=")
			s.WriteString(t.dbmap.Dialect.BindVar(x))

			plan.keyFields = append(plan.keyFields, k.fieldName)
			plan.argFields = append(plan.argFields, k.fieldName)
		}
		if plan.versField != "" {
			s.WriteString(" and ")
			s.WriteString(t.dbmap.Dialect.QuoteField(t.version.ColumnName))
			s.WriteString("=")
			s.WriteString(t.dbmap.Dialect.BindVar(len(plan.argFields)))

			plan.argFields = append(plan.argFields, plan.versField)
		}
		s.WriteString(t.dbmap.Dialect.QuerySuffix())

		plan.query = s.String()
		t.deletePlan = plan
	}

	return plan.createBindInstance(elem, t.dbmap.TypeConverter)
}

func (t *TableMap) bindGet() bindPlan {
	plan := t.getPlan
	if plan.query == "" {

		s := bytes.Buffer{}
		s.WriteString("select ")

		x := 0
		for _, col := range t.Columns {
			if !col.Transient {
				if x > 0 {
					s.WriteString(",")
				}
				s.WriteString(t.dbmap.Dialect.QuoteField(col.ColumnName))
				plan.argFields = append(plan.argFields, col.fieldName)
				x++
			}
		}
		s.WriteString(" from ")
		s.WriteString(t.dbmap.Dialect.QuotedTableForQuery(t.SchemaName, t.TableName))
		s.WriteString(" where ")
		for x := range t.keys {
			col := t.keys[x]
			if x > 0 {
				s.WriteString(" and ")
			}
			s.WriteString(t.dbmap.Dialect.QuoteField(col.ColumnName))
			s.WriteString("=")
			s.WriteString(t.dbmap.Dialect.BindVar(x))

			plan.keyFields = append(plan.keyFields, col.fieldName)
		}
		s.WriteString(t.dbmap.Dialect.QuerySuffix())

		plan.query = s.String()
		t.getPlan = plan
	}

	return plan
}

// ColumnMap represents a mapping between a Go struct field and a single
// column in a table.
// Unique and MaxSize only inform the
// CreateTables() function and are not used by Insert/Update/Delete/Get.
type ColumnMap struct {
	// Column name in db table
	ColumnName string

	// If true, this column is skipped in generated SQL statements
	Transient bool

	// If true, " unique" is added to create table statements.
	// Not used elsewhere
	Unique bool

	// Passed to Dialect.ToSqlType() to assist in informing the
	// correct column type to map to in CreateTables()
	MaxSize int

	//DbType overrides the conversion from go types to dab types
	DbType string

	// If EnforceNotNull is true then an error will be generated if
	// a zero value for this coumn is inserted/updated into a table
	EnforceNotNull bool

	DefaultValue string

	fieldName  string
	gotype     reflect.Type
	isPK       bool
	isAutoIncr bool
	isNotNull  bool
}

// IndexMap represents the data to create an index
type IndexMap struct {
	// Index name in db table
	IndexName string

	// If true, " unique" is added to the create index statement.
	Unique bool

	// List of fields for the index
	fieldNames []string
	gotype     reflect.Type
}

// This const is used to flag an index whos name should be autogenerated
const autoGenerateIndexname string = "autogenerateindexname"

// Rename allows you to specify the column name in the table
//
// Example:  table.ColMap("Updated").Rename("date_updated")
//
func (c *ColumnMap) Rename(colname string) *ColumnMap {
	c.ColumnName = colname
	return c
}

// SetTransient allows you to mark the column as transient. If true
// this column will be skipped when SQL statements are generated
func (c *ColumnMap) SetTransient(b bool) *ColumnMap {
	c.Transient = b
	return c
}

// SetUnique adds "unique" to the create table statements for this
// column, if b is true.
func (c *ColumnMap) SetUnique(b bool) *ColumnMap {
	c.Unique = b
	return c
}

// SetNotNull adds "not null" to the create table statements for this
// column, if nn is true.
func (c *ColumnMap) SetNotNull(nn bool) *ColumnMap {
	c.isNotNull = nn
	return c
}

// SetMaxSize specifies the max length of values of this column. This is
// passed to the dialect.ToSqlType() function, which can use the value
// to alter the generated type for "create table" statements
func (c *ColumnMap) SetMaxSize(size int) *ColumnMap {
	c.MaxSize = size
	return c
}

// Transaction represents a database transaction.
// Insert/Update/Delete/Get/Exec operations will be run in the context
// of that transaction.  Transactions should be terminated with
// a call to Commit() or Rollback()
type Transaction struct {
	dbmap  *DbMap
	tx     *sql.Tx
	closed bool
}

// Executor exposes the sql.DB and sql.Tx Exec function so that it can be used
// on internal functions that convert named parameters for the Exec function.
type executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// SqlExecutor exposes gorp operations that can be run from Pre/Post
// hooks.  This hides whether the current operation that triggered the
// hook is in a transaction.
//
// See the DbMap function docs for each of the functions below for more
// information.
type SqlExecutor interface {
	Get(i interface{}, keys ...interface{}) (interface{}, error)
	Insert(list ...interface{}) error
	Update(list ...interface{}) (int64, error)
	Delete(list ...interface{}) (int64, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(i interface{}, query string,
		args ...interface{}) ([]interface{}, error)
	SelectInt(query string, args ...interface{}) (int64, error)
	SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error)
	SelectFloat(query string, args ...interface{}) (float64, error)
	SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error)
	SelectStr(query string, args ...interface{}) (string, error)
	SelectNullStr(query string, args ...interface{}) (sql.NullString, error)
	SelectOne(holder interface{}, query string, args ...interface{}) error
	query(query string, args ...interface{}) (*sql.Rows, error)
	queryRow(query string, args ...interface{}) *sql.Row
}

// Compile-time check that DbMap and Transaction implement the SqlExecutor
// interface.
var _, _ SqlExecutor = &DbMap{}, &Transaction{}

type GorpLogger interface {
	Printf(format string, v ...interface{})
}

// TraceOn turns on SQL statement logging for this DbMap.  After this is
// called, all SQL statements will be sent to the logger.  If prefix is
// a non-empty string, it will be written to the front of all logged
// strings, which can aid in filtering log lines.
//
// Use TraceOn if you want to spy on the SQL statements that gorp
// generates.
//
// Note that the base log.Logger type satisfies GorpLogger, but adapters can
// easily be written for other logging packages (e.g., the golang-sanctioned
// glog framework).
func (m *DbMap) TraceOn(prefix string, logger GorpLogger) {
	m.logger = logger
	if prefix == "" {
		m.logPrefix = prefix
	} else {
		m.logPrefix = fmt.Sprintf("%s ", prefix)
	}
}

// TraceOff turns off tracing. It is idempotent.
func (m *DbMap) TraceOff() {
	m.logger = nil
	m.logPrefix = ""
}

// AddTable registers the given interface type with gorp. The table name
// will be given the name of the TypeOf(i).  You must call this function,
// or AddTableWithName, for any struct type you wish to persist with
// the given DbMap.
//
// This operation is idempotent. If i's type is already mapped, the
// existing *TableMap is returned
func (m *DbMap) AddTable(i interface{}) *TableMap {
	return m.AddTableWithName(i, "")
}

// AddTableWithName has the same behavior as AddTable, but sets
// table.TableName to name.
func (m *DbMap) AddTableWithName(i interface{}, name string) *TableMap {
	return m.AddTableWithNameAndSchema(i, "", name)
}

// AddTableWithNameAndSchema has the same behavior as AddTable, but sets
// table.TableName to name.
func (m *DbMap) AddTableWithNameAndSchema(i interface{}, schema string, name string) *TableMap {
	t := reflect.TypeOf(i)
	if name == "" {
		name = t.Name()
	}

	// check if we have a table for this type already
	// if so, update the name and return the existing pointer
	for i := range m.tables {
		table := m.tables[i]
		if table.gotype == t {
			if m.DebugLevel > 3 {
				log.Printf("AddTableWithNameAndSchema changed table name from %s to %s\n", table.TableName, name)
			}

			table.TableName = name
			return table
		}
	}

	if m.DebugLevel > 3 {
		log.Printf("AddTable type %v\n", t)
	}

	tmap := &TableMap{gotype: t, TableName: name, SchemaName: schema, dbmap: m}

	tmap.Columns = m.readStructColumns(t, tmap)

	m.tables = append(m.tables, tmap)
	if m.DebugLevel > 3 {
		for i, testim := range tmap.Indexes {
			log.Printf("tm.Indexes %d: indexname %s\n", i, testim.IndexName)
			for x, testfn := range testim.fieldNames {
				log.Printf("tm.Indexes %d: indexname %s, field %d: %s\n", i, testim.IndexName, x, testfn)
			}
		}
	}

	return tmap
}

func (m *DbMap) readStructColumns(t reflect.Type, tm *TableMap) (cols []*ColumnMap) {

	// Create slice for primary keys - initially empty
	tm.keys = make([]*ColumnMap, 0)

	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			// Recursively add nested fields in embedded structs.
			subcols := m.readStructColumns(f.Type, tm)
			// Don't append nested fields that have the same field
			// name as an already-mapped field.
			for _, subcol := range subcols {
				shouldAppend := true
				for _, col := range cols {
					if !subcol.Transient && subcol.fieldName == col.fieldName {
						shouldAppend = false
						break
					}
				}
				if shouldAppend {
					cols = append(cols, subcol)
				}
			}
		} else {
			// Parse all field tags into a GorpParsedTag
			pt := m.ParseTag(f.Tag)

			if pt.ColumnName == "" {
				pt.ColumnName = strings.Trim(strings.Split(f.Name, ",")[0], " ")
			}

			// Is this field is marked as a relation to a child/detail struct/table?
			if pt.ForeignKey != "" {

				subField := f.Type.Elem()
				if subField.Kind() == reflect.Ptr {
					subField = subField.Elem()
				}

				masterFieldName := f.Name
				subFieldValue := reflect.New(subField)

				if subFieldValue.Kind() == reflect.Ptr {
					subFieldValue = reflect.Indirect(subFieldValue)
				}

				subFieldValueInterface := subFieldValue.Interface()
				rtm := m.AddTable(subFieldValueInterface)
				r := RelationMap{DetailTable: rtm, ForeignKeyFieldName: pt.ForeignKey,
					DetailTableType: subFieldValueInterface, MasterFieldName: masterFieldName}

				tm.Relations = append(tm.Relations, &r)
			}

			gotype := f.Type
			value := reflect.New(gotype).Interface()
			if m.TypeConverter != nil {
				// Make a new pointer to a value of type gotype and
				// pass it to the TypeConverter's FromDb method to see
				// if a different type should be used for the column
				// type during table creation.
				scanner, useHolder := m.TypeConverter.FromDb(value)
				if useHolder {
					value = scanner.Holder
					gotype = reflect.TypeOf(value)
				}
			}
			if typer, ok := value.(SqlTyper); ok {
				gotype = reflect.TypeOf(typer.SqlType())
			} else if valuer, ok := value.(driver.Valuer); ok {
				// Only check for driver.Valuer if SqlTyper wasn't
				// found.
				v, err := valuer.Value()
				if err == nil && v != nil {
					gotype = reflect.TypeOf(v)
				}
			}

			cm := &ColumnMap{
				ColumnName:     pt.ColumnName,
				Transient:      pt.Transient,
				fieldName:      f.Name,
				gotype:         gotype,
				MaxSize:        pt.MaxColumnSize,
				DbType:         pt.DbType,
				isNotNull:      pt.IsNotNull,
				EnforceNotNull: pt.EnforceNotNull,
				Unique:         pt.IsFieldUnique,
				isPK:           pt.IsPk,
				isAutoIncr:     pt.IsAutoIncr,
			}
			// Check for nested fields of the same field name and
			// override them.
			shouldAppend := true
			for index, col := range cols {
				if !col.Transient && col.fieldName == cm.fieldName {
					cols[index] = cm
					shouldAppend = false
					break
				}
			}
			if shouldAppend {
				cols = append(cols, cm)
				// Collect info for Index creation from the current column
				tm.Indexes = m.addIndexForColumn(cm, f.Tag, *tm)

				if pt.IsPk {
					colmap := &ColumnMap{ColumnName: cm.ColumnName, fieldName: cm.fieldName}
					colmap.isPK = pt.IsPk
					colmap.isAutoIncr = pt.IsAutoIncr
					tm.keys = append(tm.keys, colmap)
				}

			}

		}
	}

	tm.ResetSql()

	return
}

// addIndexForColumn adds IndexMaps from field tags for one Column
// If an IndexMap already exists for the IndexName parsed from the field tags for this column,
// only the field is added to the existing IndexMap, or else a new IndexMap is created.
// The existing IndexMaps are taken from the input TableMap (tm)
// Example:
// type Post struct {
//	Id       uint64
//	User     string `db:"index:idx_user, othertag:tagid1"`
//	Err      error `db:"-"` // ignore this field when storing with gorp or gorm
// }
func (m *DbMap) addIndexForColumn(cm *ColumnMap, tag reflect.StructTag, tm TableMap) []*IndexMap {
	var shouldAppend bool
	var indexes []*IndexMap
	indexes = tm.Indexes

	// Parse all field tags into a GorpParsedTag
	pt := m.ParseTag(tag)

	for _, it := range pt.Indexes {

		if it.IndexName == "" {
			// No index tag found, early return
			continue
		}

		var im *IndexMap

		// Get all params from tagstring
		if it.IndexName == autoGenerateIndexname {
			// Index name not set, create one from the tablename + fieldname
			// The table name is necessary for postgres, as the indexnames are global inside one schema
			it.IndexName = "ix_gorp_autoindex_" + strings.ToLower(tm.TableName) + "_" + strings.ToLower(cm.fieldName)
		}
		shouldAppend = true

		fn := cm.ColumnName
		if fn == "" {
			fn = cm.fieldName
		}
		if fn == "" {
			// Unknown fieldname
			panic("func addIndexForColumn: No columnname found for index " + it.IndexName)
			break
		}

		// Check if indexName already exists in array of IndexMaps
		for _, im = range indexes {
			if im.IndexName == it.IndexName {
				im.fieldNames = append(im.fieldNames, fn)
				shouldAppend = false

				if m.DebugLevel > 3 {
					log.Println("addIndexForColumn append: index: " + it.IndexName + " field: cm.fieldName: " + cm.fieldName)
					log.Println("addIndexForColumn append: index: " + it.IndexName + " field: cm.ColumnName: " + cm.ColumnName)
				}
				break
			}
		}

		// Index not found, append it to array of IndexMaps
		if shouldAppend {

			im = &IndexMap{
				IndexName:  it.IndexName,
				Unique:     it.IsIndexUnique,
				fieldNames: []string{fn},
			}
			indexes = append(indexes, im)

			if m.DebugLevel > 3 {
				log.Println("addIndexForColumn new: index: " + it.IndexName + " field: cm.fieldName: " + cm.fieldName)
				log.Println("addIndexForColumn new: index: " + it.IndexName + " field: cm.ColumnName: " + cm.ColumnName)
			}
		}
	}
	return indexes
}

// CreateTables iterates through TableMaps registered to this DbMap and
// executes "create table" statements against the database for each.
//
// This is particularly useful in unit tests where you want to create
// and destroy the schema automatically.
func (m *DbMap) CreateTables() error {
	return m.createTables(false)
}

// CreateTablesIfNotExists is similar to CreateTables, but starts
// each statement with "create table if not exists" so that existing
// tables do not raise errors
func (m *DbMap) CreateTablesIfNotExists() error {
	return m.createTables(true)
}

func (m *DbMap) createTables(ifNotExists bool) error {
	var err error
	for i := range m.tables {
		table := m.tables[i]

		s := bytes.Buffer{}

		if strings.TrimSpace(table.SchemaName) != "" {
			schemaCreate := "create schema"
			if ifNotExists {
				s.WriteString(m.Dialect.IfSchemaNotExists(schemaCreate, table.SchemaName))
			} else {
				s.WriteString(schemaCreate)
			}
			s.WriteString(fmt.Sprintf(" %s;", table.SchemaName))
		}

		tableCreate := "create table"
		if ifNotExists {
			s.WriteString(m.Dialect.IfTableNotExists(tableCreate, table.SchemaName, table.TableName))
		} else {
			s.WriteString(tableCreate)
		}
		s.WriteString(fmt.Sprintf(" %s (", m.Dialect.QuotedTableForQuery(table.SchemaName, table.TableName)))

		x := 0
		for _, col := range table.Columns {
			if !col.Transient {
				if x > 0 {
					s.WriteString(", ")
				}
				// Check if the db type has been overriden
				stype := col.DbType
				if stype == "" {
					stype = m.Dialect.ToSqlType(col.gotype, col.MaxSize, col.isAutoIncr)
				}
				s.WriteString(fmt.Sprintf("%s %s", m.Dialect.QuoteField(col.ColumnName), stype))

				if col.isPK || col.isNotNull {
					s.WriteString(" not null")
				}
				if col.isPK && len(table.keys) == 1 {
					s.WriteString(" primary key")
				}
				if col.Unique {
					s.WriteString(" unique")
				}
				if col.isAutoIncr {
					s.WriteString(fmt.Sprintf(" %s", m.Dialect.AutoIncrStr()))
				}

				x++
			}
		}
		if len(table.keys) > 1 {
			s.WriteString(", primary key (")
			for x := range table.keys {
				if x > 0 {
					s.WriteString(", ")
				}
				s.WriteString(m.Dialect.QuoteField(table.keys[x].ColumnName))
			}
			s.WriteString(")")
		}
		if len(table.uniqueTogether) > 0 {
			for _, columns := range table.uniqueTogether {
				s.WriteString(", unique (")
				for i, column := range columns {
					if i > 0 {
						s.WriteString(", ")
					}
					s.WriteString(m.Dialect.QuoteField(column))
				}
				s.WriteString(")")
			}
		}

		s.WriteString(") ")
		s.WriteString(m.Dialect.CreateTableSuffix())
		s.WriteString(m.Dialect.QuerySuffix())

		_, err = m.Exec(s.String())
		if err != nil {
			break
		}
	}

	return err
}

// Creates indexes from a list of IndexMaps in the TableMap
// If the index already exists it is checked if the index in the database has all
// the same fields as in the TableMap
func (m *DbMap) CreateIndexes() error {
	return m.createIndexes(false)
}

// Creates indexes from a list of IndexMaps in the TableMap
// If the index already exists nothing is done, that means no checking is done
// if the index is correct. This is the suggested method for production databases,
// as index creation can have a heavy performance impact and you dont want to
// have a drop and recreate index unexpectedly on a large table.
func (m *DbMap) CreateIndexesIfNotExists() error {
	return m.createIndexes(true)
}

// Creates indexes from a list of IndexMaps in the TableMap
// if ifNotExists == false a checking is done if all fields of the IndexMap
// match with the database index
func (m *DbMap) createIndexes(ifNotExists bool) error {
	var err error

	for _, table := range m.tables {
		for _, index := range table.Indexes {
			var exists bool
			var matches bool

			exists, matches, err = m.checkIfIndexMatches(table, index)
			if err != nil {
				err = errors.New("checkIfIndexMatches for index " + index.IndexName + " failed: " + err.Error())
				return err
			}

			// DEBUG
			if m.DebugLevel > 3 {
				log.Printf("Index: %s, on table %s, exists %t, matches: %t\n", index.IndexName, table.TableName, exists, matches)
			}

			if exists {
				if matches {
					continue
				} else {
					if ifNotExists {
						continue
					} else {
						_, err = m.Exec(m.Dialect.DropIndex(table, index.IndexName))
						if err != nil {
							err = errors.New("drop index failed on index " + index.IndexName)
							break
						}
					}
				}
			}

			// Build the create index sql string
			var indexCreate string
			if index.Unique {
				indexCreate = "create unique index "
			} else {
				indexCreate = "create index "
			}

			s := bytes.Buffer{}
			s.WriteString(indexCreate)
			s.WriteString(strings.Trim(fmt.Sprintf(" %s ", m.Dialect.BuildIndexName(table.TableName, index.IndexName)), " "))
			s.WriteString(fmt.Sprintf(" on %s (", m.Dialect.QuotedTableForQuery(table.SchemaName, table.TableName)))

			sep := ""
			for _, field := range index.fieldNames {
				s.WriteString(sep + m.Dialect.QuoteField(field))
				sep = ","
			}
			s.WriteString(")")
			_, err = m.Exec(s.String())
			if err != nil {
				err = errors.New("Create index " + index.IndexName + " failed: " + err.Error())
				break
			}
		}
	}

	return err
}

// Tests if an index already exists for a table
// and if the fields in the index matches the given input IndexMap
func (m *DbMap) checkIfIndexMatches(table *TableMap, index *IndexMap) (exists bool, matches bool, err error) {
	var rows *sql.Rows
	var columnList []string
	var columnName string

	sql := m.Dialect.IfIndexExists(table.TableName, index.IndexName, table.SchemaName)

	rows, err = m.query(sql)
	if err != nil {
		err = errors.New("Getting index info failed: SQL: " + sql + " Error: " + err.Error())
		return
	}
	defer rows.Close()

	// Loop over rows and fill the columnList
	for rows.Next() {
		err = rows.Scan(&columnName)
		if err != nil {
			err = errors.New("Scanning query results for getting index info failed: " + err.Error())
			return
		}
		exists = true
		columnList = append(columnList, columnName)
	}

	matches = true
	var columnFound bool

	for _, i := range index.fieldNames {

		columnFound = false

		// Test if columnName is part of the IndexMap fields list
		for _, c := range columnList {
			if strings.ToLower(c) == strings.ToLower(i) {
				columnFound = true
			}
		}
		if !columnFound {
			matches = false
			break
		}
	}

	return
}

// DropIndex drops an individual index.  Will throw an error
// if the index does not exist.
func (m *DbMap) DropIndex(table *TableMap, index string) error {

	_, err := m.Exec(m.Dialect.DropIndex(table, index))
	return err
}

// DropTable drops an individual table.  Will throw an error
// if the table does not exist.
func (m *DbMap) DropTable(table interface{}) error {
	t := reflect.TypeOf(table)
	return m.dropTable(t, false)
}

// DropTable drops an individual table.  Will NOT throw an error
// if the table does not exist.
func (m *DbMap) DropTableIfExists(table interface{}) error {
	t := reflect.TypeOf(table)
	return m.dropTable(t, true)
}

// DropTables iterates through TableMaps registered to this DbMap and
// executes "drop table" statements against the database for each.
func (m *DbMap) DropTables() error {
	return m.dropTables(false)
}

// DropTablesIfExists is the same as DropTables, but uses the "if exists" clause to
// avoid errors for tables that do not exist.
func (m *DbMap) DropTablesIfExists() error {
	return m.dropTables(true)
}

// Goes through all the registered tables, dropping them one by one.
// If an error is encountered, then it is returned and the rest of
// the tables are not dropped.
func (m *DbMap) dropTables(addIfExists bool) (err error) {
	for _, table := range m.tables {
		err = m.dropTableImpl(table, addIfExists)
		if err != nil {
			return
		}
	}
	return err
}

// Implementation of dropping a single table.
func (m *DbMap) dropTable(t reflect.Type, addIfExists bool) error {
	table := tableOrNil(m, t)
	if table == nil {
		return errors.New(fmt.Sprintf("table %s was not registered!", table.TableName))
	}

	return m.dropTableImpl(table, addIfExists)
}

func (m *DbMap) dropTableImpl(table *TableMap, ifExists bool) (err error) {
	tableDrop := "drop table"
	if ifExists {
		tableDrop = m.Dialect.IfTableExists(tableDrop, table.SchemaName, table.TableName)
	}
	_, err = m.Exec(fmt.Sprintf("%s %s;", tableDrop, m.Dialect.QuotedTableForQuery(table.SchemaName, table.TableName)))
	return err
}

// TruncateTables iterates through TableMaps registered to this DbMap and
// executes "truncate table" statements against the database for each, or in the case of
// sqlite, a "delete from" with no "where" clause, which uses the truncate optimization
// (http://www.sqlite.org/lang_delete.html)
func (m *DbMap) TruncateTables() error {
	var err error
	for i := range m.tables {
		table := m.tables[i]
		_, e := m.Exec(fmt.Sprintf("%s %s;", m.Dialect.TruncateClause(), m.Dialect.QuotedTableForQuery(table.SchemaName, table.TableName)))
		if e != nil {
			err = e
		}
	}
	return err
}

// Insert runs a SQL INSERT statement for each element in list.
// List items must be pointers.
//
// Any interface whose TableMap has an auto-increment primary key will
// have its last insert id bound to the PK field on the struct.
//
// The hook functions PreInsert() and/or PostInsert() will be executed
// before/after the INSERT statement if the interface defines them.
//
// Panics if any interface in the list has not been registered with AddTable
func (m *DbMap) Insert(list ...interface{}) error {
	return insert(m, m, false, list...)
}

// InsertWithChilds runs a SQL INSERT statement for each element in list.
// If nested structures exist in one of the elements in list, they are
// inserted, too.
// The nested structs to insert are read from the relations list in the
// tablemap, which are populated during m.AddTable
//
// Any interface whose TableMap has an auto-increment primary key will
// have its last insert id bound to the PK field on the struct.
//
// The hook functions PreInsert() and/or PostInsert() will be executed
// before/after the INSERT statement if the interface defines them.
//
// Panics if any interface in the list has not been registered with AddTable
func (m *DbMap) InsertWithChilds(list ...interface{}) error {
	return insert(m, m, true, list...)
}

/*
// Store checks for each element in the list if it is already present in the
// database by checking on the primary key. If not present an SQL INSERT is done,
// else an SQL UPDATE
//
// Any interface whose TableMap has an auto-increment primary key will
// have its last insert id bound to the PK field on the struct.
//
// The hook functions PreInsert() and/or PostInsert() will be executed
// before/after the INSERT statement if the interface defines them.
//
// The hook functions PreUpdate() and/or PostUpdate() will be executed
// before/after the UPDATE statement if the interface defines them.
//
// Panics if any interface in the list has not been registered with AddTable
func (m *DbMap) Store(list ...interface{}) error {
	var err error
	for _, ptr := range list {
		// Check if a pointer to reflect.Value has been passed
		if reflect.TypeOf(ptr).String() == "*reflect.Value" {
			// Indirect from Pointer to Value
			ptr = *ptr.(*reflect.Value)
		}

		if reflect.TypeOf(ptr).String() == "reflect.Value" {
			fmt.Printf("+++++++++++++++ Store type %s\n", reflect.TypeOf(ptr))
			fmt.Printf("+++++++++++++++ Store type2 '%s'\n", reflect.TypeOf(ptr).String())
			fmt.Printf("+++++++++++++++ Store name %v\n", reflect.TypeOf(ptr).Name())
			err = m.InsertFromValue(m, ptr.(reflect.Value))

			if err != nil {
				break
			}

		}

	}
	return err
}
*/

/*
func (m *DbMap) InsertFromValue(exec SqlExecutor, value reflect.Value) error {

	table, err := m.TableFor(value.Type(), true)
	//table, elem, err := m.tableForPointer(ptr, false)
	elem := value
	if err != nil {
		return err
	}

	if m.DebugLevel > 2 {
		fmt.Printf("GORP InsertFromValue table '%s', elem %v\n", table.TableName, elem)
	}

	eval := elem.Addr().Interface()
	if v, ok := eval.(HasPreInsert); ok {
		err := v.PreInsert(exec)
		if err != nil {
			return err
		}
	}

	bi, err := table.bindInsert(elem)
	if err != nil {
		return err
	}

	if bi.autoIncrIdx > -1 {
		f := elem.FieldByName(bi.autoIncrFieldName)
		switch inserter := m.Dialect.(type) {
		case IntegerAutoIncrInserter:
			id, err := inserter.InsertAutoIncr(exec, bi.query, bi.args...)
			if err != nil {
				return err
			}
			k := f.Kind()
			if (k == reflect.Int) || (k == reflect.Int16) || (k == reflect.Int32) || (k == reflect.Int64) {
				f.SetInt(id)
			} else if (k == reflect.Uint) || (k == reflect.Uint16) || (k == reflect.Uint32) || (k == reflect.Uint64) {
				f.SetUint(uint64(id))
			} else {
				return fmt.Errorf("gorp: Cannot set autoincrement value on non-Int field. SQL=%s  autoIncrIdx=%d autoIncrFieldName=%s", bi.query, bi.autoIncrIdx, bi.autoIncrFieldName)
			}
		case TargetedAutoIncrInserter:
			err := inserter.InsertAutoIncrToTarget(exec, bi.query, f.Addr().Interface(), bi.args...)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("gorp: Cannot use autoincrement fields on dialects that do not implement an autoincrementing interface")
		}
	} else {
		_, err := exec.Exec(bi.query, bi.args...)
		if err != nil {
			return err
		}
	}

	// Insert child records if present
	for _, r := range table.Relations {

		if m.DebugLevel > 2 {
			fmt.Printf("******GORP RELATION '%s'\n", r.String())
		}

		//err = m.Insert(r.DetailTableType)
		//if err != nil {
		//	return err
		//}
	}

	if v, ok := eval.(HasPostInsert); ok {
		err := v.PostInsert(exec)
		if err != nil {
			return err
		}
	}

	return nil
}
*/

// Update runs a SQL UPDATE statement for each element in list.  List
// items must be pointers.
//
// The hook functions PreUpdate() and/or PostUpdate() will be executed
// before/after the UPDATE statement if the interface defines them.
//
// Returns the number of rows updated.
//
// Returns an error if SetKeys has not been called on the TableMap
// Panics if any interface in the list has not been registered with AddTable
func (m *DbMap) Update(list ...interface{}) (int64, error) {
	return update(m, m, false, list...)
}

// UpdateWithChilds runs a SQL UPDATE statement for each element in list.
// If nested structures exist in one of the elements in list, they are
// inserted or updated, too.
// The nested structs to update are read from the relations list in the
// tablemap, which are populated during m.AddTable
//
// The hook functions PreUpdate() and/or PostUpdate() will be executed
// before/after the UPDATE statement if the interface defines them.
//
// Returns the number of rows updated and the number of childs updated
// Panics if any interface in the list has not been registered with AddTable
func (m *DbMap) UpdateWithChilds(list ...interface{}) (int64, error) {
	return update(m, m, true, list...)
}

// Delete runs a SQL DELETE statement for each element in list.  List
// items must be pointers.
//
// The hook functions PreDelete() and/or PostDelete() will be executed
// before/after the DELETE statement if the interface defines them.
//
// Returns the number of rows deleted.
//
// Returns an error if SetKeys has not been called on the TableMap
// Panics if any interface in the list has not been registered with AddTable
func (m *DbMap) Delete(list ...interface{}) (int64, error) {
	return delete(m, m, list...)
}

// Get runs a SQL SELECT to fetch a single row from the table based on the
// primary key(s)
//
// i should be an empty value for the struct to load.  keys should be
// the primary key value(s) for the row to load.  If multiple keys
// exist on the table, the order should match the column order
// specified in SetKeys() when the table mapping was defined.
//
// The hook function PostGet() will be executed after the SELECT
// statement if the interface defines them.
//
// Returns a pointer to a struct that matches or nil if no row is found.
//
// Returns an error if SetKeys has not been called on the TableMap
// Panics if any interface in the list has not been registered with AddTable
func (m *DbMap) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	return get(m, m, i, false, 0, 0, keys...)
}

// GetWithChilds runs a SQL SELECT to fetch a single row from the table based on the
// primary key(s). All child records are fetched if a RelationMap exists for this table
//
// i should be an empty value for the struct to load.  keys should be
// the primary key value(s) for the row to load.  If multiple keys
// exist on the table, the order should match the column order
// specified in SetKeys() when the table mapping was defined.
//
// The hook function PostGet() will be executed after the SELECT
// statement if the interface defines them.
//
// Returns a pointer to a struct that matches or nil if no row is found.
//
// Returns an error if SetKeys has not been called on the TableMap
// Panics if any interface in the list has not been registered with AddTable
func (m *DbMap) GetWithChilds(i interface{}, ChildLimit int64, ChildOffset int64, keys ...interface{}) (interface{}, error) {
	return get(m, m, i, true, ChildLimit, ChildOffset, keys...)
}

// Select runs an arbitrary SQL query, binding the columns in the result
// to fields on the struct specified by i.  args represent the bind
// parameters for the SQL statement.
//
// Column names on the SELECT statement should be aliased to the field names
// on the struct i. Returns an error if one or more columns in the result
// do not match.  It is OK if fields on i are not part of the SQL
// statement.
//
// The hook function PostGet() will be executed after the SELECT
// statement if the interface defines them.
//
// Values are returned in one of two ways:
// 1. If i is a struct or a pointer to a struct, returns a slice of pointers to
// matching rows of type i.
// 2. If i is a pointer to a slice, the results will be appended to that slice
// and nil returned.
//
// i does NOT need to be registered with AddTable()
func (m *DbMap) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return hookedselect(m, m, i, query, args...)
}

// Exec runs an arbitrary SQL statement.  args represent the bind parameters.
// This is equivalent to running:  Exec() using database/sql
func (m *DbMap) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.logger != nil {
		now := time.Now()
		defer m.trace(now, query, args...)
	}
	return exec(m, query, args...)
}

// SelectInt is a convenience wrapper around the gorp.SelectInt function
func (m *DbMap) SelectInt(query string, args ...interface{}) (int64, error) {
	return SelectInt(m, query, args...)
}

// SelectNullInt is a convenience wrapper around the gorp.SelectNullInt function
func (m *DbMap) SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error) {
	return SelectNullInt(m, query, args...)
}

// SelectFloat is a convenience wrapper around the gorp.SelectFloat function
func (m *DbMap) SelectFloat(query string, args ...interface{}) (float64, error) {
	return SelectFloat(m, query, args...)
}

// SelectNullFloat is a convenience wrapper around the gorp.SelectNullFloat function
func (m *DbMap) SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error) {
	return SelectNullFloat(m, query, args...)
}

// SelectStr is a convenience wrapper around the gorp.SelectStr function
func (m *DbMap) SelectStr(query string, args ...interface{}) (string, error) {
	return SelectStr(m, query, args...)
}

// SelectNullStr is a convenience wrapper around the gorp.SelectNullStr function
func (m *DbMap) SelectNullStr(query string, args ...interface{}) (sql.NullString, error) {
	return SelectNullStr(m, query, args...)
}

// SelectOne is a convenience wrapper around the gorp.SelectOne function
func (m *DbMap) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return SelectOne(m, m, holder, query, args...)
}

// Begin starts a gorp Transaction
func (m *DbMap) Begin() (*Transaction, error) {
	if m.logger != nil {
		now := time.Now()
		defer m.trace(now, "begin;")
	}
	tx, err := m.Db.Begin()
	if err != nil {
		return nil, err
	}
	return &Transaction{m, tx, false}, nil
}

// TableFor returns the *TableMap corresponding to the given Go Type
// If no table is mapped to that type an error is returned.
// If checkPK is true and the mapped table has no registered PKs, an error is returned.
func (m *DbMap) TableFor(t reflect.Type, checkPK bool) (*TableMap, error) {
	table := tableOrNil(m, t)
	if table == nil {
		return nil, errors.New(fmt.Sprintf("No table found for type: %v", t))
	}

	if checkPK && len(table.keys) < 1 {
		e := fmt.Sprintf("gorp: No keys defined for table: %s",
			table.TableName)
		return nil, errors.New(e)
	}

	return table, nil
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
// This is equivalent to running:  Prepare() using database/sql
func (m *DbMap) Prepare(query string) (*sql.Stmt, error) {
	if m.logger != nil {
		now := time.Now()
		defer m.trace(now, query, nil)
	}
	return m.Db.Prepare(query)
}

func tableOrNil(m *DbMap, t reflect.Type) *TableMap {
	for i := range m.tables {
		table := m.tables[i]
		if table.gotype == t {
			return table
		}
	}
	return nil
}

func (m *DbMap) tableForPointer(ptr interface{}, checkPK bool) (*TableMap, reflect.Value, error) {
	ptrv := reflect.ValueOf(ptr)
	if ptrv.Kind() != reflect.Ptr {
		e := fmt.Sprintf("gorp: passed non-pointer: %v (kind=%v)", ptr,
			ptrv.Kind())
		return nil, reflect.Value{}, errors.New(e)
	}
	elem := ptrv.Elem()
	etype := reflect.TypeOf(elem.Interface())
	t, err := m.TableFor(etype, checkPK)
	if err != nil {
		return nil, reflect.Value{}, err
	}

	return t, elem, nil
}

func (m *DbMap) queryRow(query string, args ...interface{}) *sql.Row {
	if m.logger != nil {
		now := time.Now()
		defer m.trace(now, query, args...)
	}
	return m.Db.QueryRow(query, args...)
}

func (m *DbMap) query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.logger != nil {
		now := time.Now()
		defer m.trace(now, query, args...)
	}
	return m.Db.Query(query, args...)
}

func (m *DbMap) trace(started time.Time, query string, args ...interface{}) {
	if m.logger != nil {
		var margs = argsString(args...)
		m.logger.Printf("%s%s [%s] (%v)", m.logPrefix, query, margs, (time.Now().Sub(started)))
	}
}

func argsString(args ...interface{}) string {
	var margs string
	for i, a := range args {
		var v interface{} = a
		if x, ok := v.(driver.Valuer); ok {
			y, err := x.Value()
			if err == nil {
				v = y
			}
		}
		switch v.(type) {
		case string:
			v = fmt.Sprintf("%q", v)
		default:
			v = fmt.Sprintf("%v", v)
		}
		margs += fmt.Sprintf("%d:%s", i+1, v)
		if i+1 < len(args) {
			margs += " "
		}
	}
	return margs
}

// GorpParsedTag contains info about a column parsed from a tag
type GorpParsedTag struct {
	ColumnName     string
	IsFieldUnique  bool
	Indexes        []GorpParsedIndexTag
	MaxColumnSize  int
	DbType         string
	IsNotNull      bool
	EnforceNotNull bool
	IsAutoIncr     bool
	IsPk           bool
	Transient      bool
	ForeignKey     string
}

func (pt GorpParsedTag) String() string {
	return "ColumnName " + pt.ColumnName + "\n" +
		"DbType " + pt.DbType + "\n"
}

// GorpParsedIndexTag contains index info parsed from a tag
type GorpParsedIndexTag struct {
	ColumnName    string
	IndexName     string
	IsIndexUnique bool
	ForeignKey    string
}

// ParseTag extracts all field tags from input param tag and returns all found options
// Tag key can be ether "db" (the legacy default) or "gorp"
// "gorp" has been added in this fork only, the intent is to avoid namespace conflicts
// with other database packages for go
/* Example:
type Post struct {
	Id           uint64    `db:"notnull, PID, primarykey, autoincrement"`
	SecondTestID int       `db:"notnull, name: SID, uniqueindex:idx_unique_sid"`
	Created      time.Time `db:"notnull, primarykey"`
	PostDate     time.Time `db:"notnull"`
	Site         string    `db:"name: PostSite, notnull, size:50, index:idx_site"`
	PostId       string    `db:"notnull, size:32, unique"`
	Score        int       `db:"notnull"`
	Title        string    `db:"notnull, size:1024"`
	Url          string    `db:"notnull"`
	User         string    `db:"index:idx_user, size:64"`
	PostSub      string    `db:"index:idx_user, size:128"`
	UserIP       string    `db:"notnull, size:16"`
	BodyType     string    `db:"notnull, size:64"`
	Body         string    `db:"name:PostBody, type:mediumtext"`
	Err          error     `db:"-"` // ignore this field when storing with gorp
}
*/
func (m *DbMap) ParseTag(tag reflect.StructTag) (pt GorpParsedTag) {

	// Get the tags from the struct field ether by tagname "gorp" or "db"
	ts := tag.Get("gorp")
	if ts == "" {
		ts = tag.Get("db")
	}
	if ts == "" {
		// not tags found, exit early
		return
	}

	if ts == "-" {
		// Ignore this column
		pt.ColumnName = "-"
		pt.Transient = true
	} else {

		// Get all params from tagstring
		tags := strings.Split(ts, ",")
		for _, tag := range tags {
			o := strings.Split(tag, ":")
			o[0] = strings.ToLower(strings.Trim(o[0], " "))

			switch o[0] {
			case "name":
				pt.ColumnName = strings.Trim(o[1], " ")
			case "index":
				it := GorpParsedIndexTag{}
				it.IndexName = strings.Trim(o[1], " ")
				if it.IndexName == "" {
					it.IndexName = autoGenerateIndexname
				}
				it.IsIndexUnique = false
				pt.Indexes = append(pt.Indexes, it)
			case "uniqueindex":
				it := GorpParsedIndexTag{}
				it.IndexName = strings.Trim(o[1], " ")
				if it.IndexName == "" {
					it.IndexName = autoGenerateIndexname
				}
				it.IsIndexUnique = true
				pt.Indexes = append(pt.Indexes, it)
			case "size":
				var ErrAtoi error
				pt.MaxColumnSize, ErrAtoi = strconv.Atoi(strings.Trim(o[1], " "))
				if ErrAtoi != nil {
					panic(fmt.Sprintf("Int conversion for tag 'size:%s' failed: %s", o[1], ErrAtoi.Error()))
				}
			case "type":
				pt.DbType = strings.Trim(o[1], " ")
			case "notnull":
				pt.IsNotNull = true
			case "enforcenotnull":
				pt.IsNotNull = true
				pt.EnforceNotNull = true
			case "unique":
				pt.IsFieldUnique = true
			case "autoincrement":
				pt.IsAutoIncr = true
			case "primarykey":
				pt.IsPk = true
			case "relation":
				pt.Transient = true
				pt.ForeignKey = strings.Trim(o[1], " ")
			case "ignorefield":
				pt.Transient = true

			default:
				// Fallback to traditional gorp tags - use it as a fieldname if it is none of the tags above
				if len(o) == 1 {
					pt.ColumnName = o[0]
				}
			}
		}
	}

	return
}

// Insert has the same behavior as DbMap.Insert(), but runs in a transaction.
func (t *Transaction) Insert(list ...interface{}) error {
	return insert(t.dbmap, t, false, list...)
}

// Update had the same behavior as DbMap.Update(), but runs in a transaction.
func (t *Transaction) Update(list ...interface{}) (int64, error) {
	return update(t.dbmap, t, false, list...)
}

// Delete has the same behavior as DbMap.Delete(), but runs in a transaction.
func (t *Transaction) Delete(list ...interface{}) (int64, error) {
	return delete(t.dbmap, t, list...)
}

// Get has the same behavior as DbMap.Get(), but runs in a transaction.
func (t *Transaction) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	return get(t.dbmap, t, i, false, 0, 0, keys...)
}

// Select has the same behavior as DbMap.Select(), but runs in a transaction.
func (t *Transaction) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return hookedselect(t.dbmap, t, i, query, args...)
}

// Exec has the same behavior as DbMap.Exec(), but runs in a transaction.
func (t *Transaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	if t.dbmap.logger != nil {
		now := time.Now()
		defer t.dbmap.trace(now, query, args...)
	}
	return exec(t, query, args...)
}

// SelectInt is a convenience wrapper around the gorp.SelectInt function.
func (t *Transaction) SelectInt(query string, args ...interface{}) (int64, error) {
	return SelectInt(t, query, args...)
}

// SelectNullInt is a convenience wrapper around the gorp.SelectNullInt function.
func (t *Transaction) SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error) {
	return SelectNullInt(t, query, args...)
}

// SelectFloat is a convenience wrapper around the gorp.SelectFloat function.
func (t *Transaction) SelectFloat(query string, args ...interface{}) (float64, error) {
	return SelectFloat(t, query, args...)
}

// SelectNullFloat is a convenience wrapper around the gorp.SelectNullFloat function.
func (t *Transaction) SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error) {
	return SelectNullFloat(t, query, args...)
}

// SelectStr is a convenience wrapper around the gorp.SelectStr function.
func (t *Transaction) SelectStr(query string, args ...interface{}) (string, error) {
	return SelectStr(t, query, args...)
}

// SelectNullStr is a convenience wrapper around the gorp.SelectNullStr function.
func (t *Transaction) SelectNullStr(query string, args ...interface{}) (sql.NullString, error) {
	return SelectNullStr(t, query, args...)
}

// SelectOne is a convenience wrapper around the gorp.SelectOne function.
func (t *Transaction) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return SelectOne(t.dbmap, t, holder, query, args...)
}

// Commit commits the underlying database transaction.
func (t *Transaction) Commit() error {
	if !t.closed {
		t.closed = true
		if t.dbmap.logger != nil {
			now := time.Now()
			defer t.dbmap.trace(now, "commit;")
		}
		return t.tx.Commit()
	}

	return sql.ErrTxDone
}

// Rollback rolls back the underlying database transaction.
func (t *Transaction) Rollback() error {
	if !t.closed {
		t.closed = true
		if t.dbmap.logger != nil {
			now := time.Now()
			defer t.dbmap.trace(now, "rollback;")
		}
		return t.tx.Rollback()
	}

	return sql.ErrTxDone
}

// Savepoint creates a savepoint with the given name. The name is interpolated
// directly into the SQL SAVEPOINT statement, so you must sanitize it if it is
// derived from user input.
func (t *Transaction) Savepoint(name string) error {
	query := "savepoint " + t.dbmap.Dialect.QuoteField(name)
	if t.dbmap.logger != nil {
		now := time.Now()
		defer t.dbmap.trace(now, query, nil)
	}
	_, err := t.tx.Exec(query)
	return err
}

// RollbackToSavepoint rolls back to the savepoint with the given name. The
// name is interpolated directly into the SQL SAVEPOINT statement, so you must
// sanitize it if it is derived from user input.
func (t *Transaction) RollbackToSavepoint(savepoint string) error {
	query := "rollback to savepoint " + t.dbmap.Dialect.QuoteField(savepoint)
	if t.dbmap.logger != nil {
		now := time.Now()
		defer t.dbmap.trace(now, query, nil)
	}
	_, err := t.tx.Exec(query)
	return err
}

// ReleaseSavepint releases the savepoint with the given name. The name is
// interpolated directly into the SQL SAVEPOINT statement, so you must sanitize
// it if it is derived from user input.
func (t *Transaction) ReleaseSavepoint(savepoint string) error {
	query := "release savepoint " + t.dbmap.Dialect.QuoteField(savepoint)
	if t.dbmap.logger != nil {
		now := time.Now()
		defer t.dbmap.trace(now, query, nil)
	}
	_, err := t.tx.Exec(query)
	return err
}

// Prepare has the same behavior as DbMap.Prepare(), but runs in a transaction.
func (t *Transaction) Prepare(query string) (*sql.Stmt, error) {
	if t.dbmap.logger != nil {
		now := time.Now()
		defer t.dbmap.trace(now, query, nil)
	}
	return t.tx.Prepare(query)
}

func (t *Transaction) queryRow(query string, args ...interface{}) *sql.Row {
	if t.dbmap.logger != nil {
		now := time.Now()
		defer t.dbmap.trace(now, query, args...)
	}
	return t.tx.QueryRow(query, args...)
}

func (t *Transaction) query(query string, args ...interface{}) (*sql.Rows, error) {
	if t.dbmap.logger != nil {
		now := time.Now()
		defer t.dbmap.trace(now, query, args...)
	}
	return t.tx.Query(query, args...)
}

///////////////

// SelectInt executes the given query, which should be a SELECT statement for a single
// integer column, and returns the value of the first row returned.  If no rows are
// found, zero is returned.
func SelectInt(e SqlExecutor, query string, args ...interface{}) (int64, error) {
	var h int64
	err := selectVal(e, &h, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return h, nil
}

// SelectNullInt executes the given query, which should be a SELECT statement for a single
// integer column, and returns the value of the first row returned.  If no rows are
// found, the empty sql.NullInt64 value is returned.
func SelectNullInt(e SqlExecutor, query string, args ...interface{}) (sql.NullInt64, error) {
	var h sql.NullInt64
	err := selectVal(e, &h, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return h, err
	}
	return h, nil
}

// SelectFloat executes the given query, which should be a SELECT statement for a single
// float column, and returns the value of the first row returned. If no rows are
// found, zero is returned.
func SelectFloat(e SqlExecutor, query string, args ...interface{}) (float64, error) {
	var h float64
	err := selectVal(e, &h, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return h, nil
}

// SelectNullFloat executes the given query, which should be a SELECT statement for a single
// float column, and returns the value of the first row returned. If no rows are
// found, the empty sql.NullInt64 value is returned.
func SelectNullFloat(e SqlExecutor, query string, args ...interface{}) (sql.NullFloat64, error) {
	var h sql.NullFloat64
	err := selectVal(e, &h, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return h, err
	}
	return h, nil
}

// SelectStr executes the given query, which should be a SELECT statement for a single
// char/varchar column, and returns the value of the first row returned.  If no rows are
// found, an empty string is returned.
func SelectStr(e SqlExecutor, query string, args ...interface{}) (string, error) {
	var h string
	err := selectVal(e, &h, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return h, nil
}

// SelectNullStr executes the given query, which should be a SELECT
// statement for a single char/varchar column, and returns the value
// of the first row returned.  If no rows are found, the empty
// sql.NullString is returned.
func SelectNullStr(e SqlExecutor, query string, args ...interface{}) (sql.NullString, error) {
	var h sql.NullString
	err := selectVal(e, &h, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return h, err
	}
	return h, nil
}

// SelectOne executes the given query (which should be a SELECT statement)
// and binds the result to holder, which must be a pointer.
//
// If no row is found, an error (sql.ErrNoRows specifically) will be returned
//
// If more than one row is found, an error will be returned.
//
func SelectOne(m *DbMap, e SqlExecutor, holder interface{}, query string, args ...interface{}) error {
	t := reflect.TypeOf(holder)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	} else {
		return fmt.Errorf("gorp: SelectOne holder must be a pointer, but got: %t", holder)
	}

	// Handle pointer to pointer
	isptr := false
	if t.Kind() == reflect.Ptr {
		isptr = true
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		var nonFatalErr error

		list, err := hookedselect(m, e, holder, query, args...)
		if err != nil {
			if !NonFatalError(err) {
				return err
			}
			nonFatalErr = err
		}

		dest := reflect.ValueOf(holder)
		if isptr {
			dest = dest.Elem()
		}

		if list != nil && len(list) > 0 {
			// check for multiple rows
			if len(list) > 1 {
				return fmt.Errorf("gorp: multiple rows returned for: %s - %v", query, args)
			}

			// Initialize if nil
			if dest.IsNil() {
				dest.Set(reflect.New(t))
			}

			// only one row found
			src := reflect.ValueOf(list[0])
			dest.Elem().Set(src.Elem())
		} else {
			// No rows found, return a proper error.
			return sql.ErrNoRows
		}

		return nonFatalErr
	}

	return selectVal(e, holder, query, args...)
}

func selectVal(e SqlExecutor, holder interface{}, query string, args ...interface{}) error {
	if len(args) == 1 {
		switch m := e.(type) {
		case *DbMap:
			query, args = maybeExpandNamedQuery(m, query, args)
		case *Transaction:
			query, args = maybeExpandNamedQuery(m.dbmap, query, args)
		}
	}

	rows, err := e.query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}

	return rows.Scan(holder)
}

///////////////

func hookedselect(m *DbMap, exec SqlExecutor, i interface{}, query string,
	args ...interface{}) ([]interface{}, error) {

	var nonFatalErr error

	list, err := rawselect(m, exec, i, query, args...)
	if err != nil {
		if !NonFatalError(err) {
			if m.DebugLevel > 0 {
				log.Printf("[gorp] hookedselect error: %d\n", err.Error())
			}
			return nil, err
		}
		nonFatalErr = err
		if m.DebugLevel > 0 {
			log.Printf("[gorp] hookedselect nonFatalErr: %d\n", nonFatalErr.Error())
		}
	}

	// Determine where the results are: written to i, or returned in list
	if t, _ := toSliceType(i); t == nil {
		for _, v := range list {
			if v, ok := v.(HasPostGet); ok {
				err := v.PostGet(exec)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		resultsValue := reflect.Indirect(reflect.ValueOf(i))
		for i := 0; i < resultsValue.Len(); i++ {
			if v, ok := resultsValue.Index(i).Interface().(HasPostGet); ok {
				err := v.PostGet(exec)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return list, nonFatalErr
}

func rawselect(m *DbMap, exec SqlExecutor, i interface{}, query string,
	args ...interface{}) ([]interface{}, error) {
	var (
		appendToSlice   = false // Write results to i directly?
		intoStruct      = true  // Selecting into a struct?
		pointerElements = true  // Are the slice elements pointers (vs values)?
	)

	var nonFatalErr error

	// get type for i, verifying it's a supported destination
	t, err := toType(i)
	if err != nil {
		var err2 error
		if t, err2 = toSliceType(i); t == nil {
			if err2 != nil {
				return nil, err2
			}
			return nil, err
		}
		pointerElements = t.Kind() == reflect.Ptr
		if pointerElements {
			t = t.Elem()
		}
		appendToSlice = true
		intoStruct = t.Kind() == reflect.Struct
	}

	// If the caller supplied a single struct/map argument, assume a "named
	// parameter" query.  Extract the named arguments from the struct/map, create
	// the flat arg slice, and rewrite the query to use the dialect's placeholder.
	if len(args) == 1 {
		query, args = maybeExpandNamedQuery(m, query, args)
	}

	if m.DebugLevel > 3 {
		log.Printf("[gorp] rawselect start\n")
	}

	// Run the query
	rows, err := exec.query(query, args...)
	if err != nil {
		if m.DebugLevel > 0 {
			log.Printf("[gorp] rawselect exec.query error: %d\n", err.Error())
		}

		return nil, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		if m.DebugLevel > 2 {
			log.Printf("[gorp] rawselect rows error: %d\n", rows.Err())
		}
	}

	// Fetch the column names as returned from db
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	if !intoStruct && len(cols) > 1 {
		return nil, fmt.Errorf("gorp: select into non-struct slice requires 1 column, got %d", len(cols))
	}

	var colToFieldIndex [][]int
	if intoStruct {
		// TODO - try to cache the columnToFieldIndex map
		colToFieldIndex, err = columnToFieldIndex(m, t, cols)
		if err != nil {
			if !NonFatalError(err) {
				return nil, err
			}
			nonFatalErr = err
		}
	}

	conv := m.TypeConverter

	// Add results to one of these two slices.
	var (
		list       = make([]interface{}, 0)
		sliceValue = reflect.Indirect(reflect.ValueOf(i))
	)

	for {
		if !rows.Next() {
			// if error occured return rawselect
			if rows.Err() != nil {
				return nil, rows.Err()
			}
			// time to exit from outer "for" loop
			break
		}

		v := reflect.New(t)
		dest := make([]interface{}, len(cols))

		custScan := make([]CustomScanner, 0)

		for x := range cols {
			f := v.Elem()
			if intoStruct {
				index := colToFieldIndex[x]
				if index == nil {
					// this field is not present in the struct, so create a dummy
					// value for rows.Scan to scan into
					var dummy sql.RawBytes
					dest[x] = &dummy
					continue
				}
				f = f.FieldByIndex(index)
			}
			target := f.Addr().Interface()
			if conv != nil {
				scanner, ok := conv.FromDb(target)
				if ok {
					target = scanner.Holder
					custScan = append(custScan, scanner)
				}
			}
			dest[x] = target
		}

		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}

		for _, c := range custScan {
			err = c.Bind()
			if err != nil {
				return nil, err
			}
		}

		if appendToSlice {
			if !pointerElements {
				v = v.Elem()
			}
			sliceValue.Set(reflect.Append(sliceValue, v))
		} else {
			list = append(list, v.Interface())
		}
	}

	if appendToSlice && sliceValue.IsNil() {
		sliceValue.Set(reflect.MakeSlice(sliceValue.Type(), 0, 0))
	}

	return list, nonFatalErr
}

// Calls the Exec function on the executor, but attempts to expand any eligible named
// query arguments first.
func exec(e SqlExecutor, query string, args ...interface{}) (sql.Result, error) {
	var dbMap *DbMap
	var executor executor
	switch m := e.(type) {
	case *DbMap:
		executor = m.Db
		dbMap = m
	case *Transaction:
		executor = m.tx
		dbMap = m.dbmap
	}

	if len(args) == 1 {
		query, args = maybeExpandNamedQuery(dbMap, query, args)
	}

	return executor.Exec(query, args...)
}

// maybeExpandNamedQuery checks the given arg to see if it's eligible to be used
// as input to a named query.  If so, it rewrites the query to use
// dialect-dependent bindvars and instantiates the corresponding slice of
// parameters by extracting data from the map / struct.
// If not, returns the input values unchanged.
func maybeExpandNamedQuery(m *DbMap, query string, args []interface{}) (string, []interface{}) {
	var (
		arg    = args[0]
		argval = reflect.ValueOf(arg)
	)
	if argval.Kind() == reflect.Ptr {
		argval = argval.Elem()
	}

	if argval.Kind() == reflect.Map && argval.Type().Key().Kind() == reflect.String {
		return expandNamedQuery(m, query, func(key string) reflect.Value {
			return argval.MapIndex(reflect.ValueOf(key))
		})
	}
	if argval.Kind() != reflect.Struct {
		return query, args
	}
	if _, ok := arg.(time.Time); ok {
		// time.Time is driver.Value
		return query, args
	}
	if _, ok := arg.(driver.Valuer); ok {
		// driver.Valuer will be converted to driver.Value.
		return query, args
	}

	return expandNamedQuery(m, query, argval.FieldByName)
}

var keyRegexp = regexp.MustCompile(`:[[:word:]]+`)

// expandNamedQuery accepts a query with placeholders of the form ":key", and a
// single arg of Kind Struct or Map[string].  It returns the query with the
// dialect's placeholders, and a slice of args ready for positional insertion
// into the query.
func expandNamedQuery(m *DbMap, query string, keyGetter func(key string) reflect.Value) (string, []interface{}) {
	var (
		n    int
		args []interface{}
	)
	return keyRegexp.ReplaceAllStringFunc(query, func(key string) string {
		val := keyGetter(key[1:])
		if !val.IsValid() {
			return key
		}
		args = append(args, val.Interface())
		newVar := m.Dialect.BindVar(n)
		n++
		return newVar
	}), args
}

// columnToFieldIndex
func columnToFieldIndex(m *DbMap, t reflect.Type, cols []string) ([][]int, error) {
	colToFieldIndex := make([][]int, len(cols))

	// check if type t is a mapped table - if so we'll
	// check the table for column aliasing below
	tableMapped := false
	table := tableOrNil(m, t)
	if table != nil {
		tableMapped = true
	}

	// Loop over column names and find field in t to bind to
	// based on column name. all returned columns must match
	// a field in the t struct
	missingColNames := []string{}
	for x := range cols {
		colName := strings.ToLower(cols[x])

		field, found := t.FieldByNameFunc(func(fieldName string) bool {
			field, _ := t.FieldByName(fieldName)

			// Parse all field tags into a GorpParsedTag
			pt := m.ParseTag(field.Tag)

			if m.DebugLevel > 3 {
				// DEBUG
				log.Printf("columnToFieldIndex LOOKING FOR: %s\n", colName)
				log.Printf("columnToFieldIndex Name: %s\n", field.Name)
				log.Printf("columnToFieldIndex PkgPath: %s\n", field.PkgPath)
				log.Printf("columnToFieldIndex Tag: %s\n", field.Tag)
				log.Printf("columnToFieldIndex pt.ColumnName: %s\n", pt.ColumnName)
				log.Println("----- columnToFieldIndex END -----------")
			}

			if pt.Transient {
				return false
			} else if pt.ColumnName == "" {
				pt.ColumnName = field.Name
			}
			if tableMapped {
				colMap := colMapOrNil(table, pt.ColumnName)
				if colMap != nil {

					if m.DebugLevel > 3 {
						// DEBUG
						log.Printf("Changed ColumnName from %s to %s\n", pt.ColumnName, colMap.ColumnName)
					}
					pt.ColumnName = colMap.ColumnName
				}
			}

			ColMatches := (colName == strings.ToLower(pt.ColumnName))

			if m.DebugLevel > 3 {
				// DEBUG
				if ColMatches {
					log.Println("--!!!! YES MATCHES !!!!-----------")
				} else {
					log.Println("----------------------------------")
				}
			}

			return ColMatches
		})
		if found {
			colToFieldIndex[x] = field.Index
		}
		if colToFieldIndex[x] == nil {
			missingColNames = append(missingColNames, colName)
		}

		if m.DebugLevel > 3 {
			// DEBUG
			log.Printf("colToFieldIndex[x]: %v\n ", colToFieldIndex[x])
			log.Println("columnToFieldIndex colName: " + colName)
			log.Println("columnToFieldIndex fieldName: " + field.Name)
		}
	}
	if len(missingColNames) > 0 {
		return colToFieldIndex, &NoFieldInTypeError{
			TypeName:        t.Name(),
			MissingColNames: missingColNames,
		}
	}
	return colToFieldIndex, nil
}

func fieldByName(val reflect.Value, fieldName string) *reflect.Value {
	// try to find field by exact match
	f := val.FieldByName(fieldName)

	if f != zeroVal {
		return &f
	}

	// try to find by case insensitive match - only the Postgres driver
	// seems to require this - in the case where columns are aliased in the sql
	fieldNameL := strings.ToLower(fieldName)
	fieldCount := val.NumField()
	t := val.Type()
	for i := 0; i < fieldCount; i++ {
		sf := t.Field(i)
		if strings.ToLower(sf.Name) == fieldNameL {
			f := val.Field(i)
			return &f
		}
	}

	return nil
}

// toSliceType returns the element type of the given object, if the object is a
// "*[]*Element" or "*[]Element". If not, returns nil.
// err is returned if the user was trying to pass a pointer-to-slice but failed.
func toSliceType(i interface{}) (reflect.Type, error) {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr {
		// If it's a slice, return a more helpful error message
		if t.Kind() == reflect.Slice {
			return nil, fmt.Errorf("gorp: Cannot SELECT into a non-pointer slice: %v", t)
		}
		return nil, nil
	}
	if t = t.Elem(); t.Kind() != reflect.Slice {
		return nil, nil
	}
	return t.Elem(), nil
}

func toType(i interface{}) (reflect.Type, error) {
	t := reflect.TypeOf(i)

	// If a Pointer to a type, follow
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("gorp: Cannot SELECT into this type: %v", reflect.TypeOf(i))
	}
	return t, nil
}

func get(m *DbMap, exec SqlExecutor, i interface{}, getChilds bool, ChildLimit int64, ChildOffset int64,
	keys ...interface{}) (interface{}, error) {

	t, err := toType(i)
	if err != nil {
		return nil, err
	}

	table, err := m.TableFor(t, true)
	if err != nil {
		return nil, err
	}

	plan := table.bindGet()

	v := reflect.New(t)
	dest := make([]interface{}, len(plan.argFields))

	conv := m.TypeConverter
	custScan := make([]CustomScanner, 0)

	for x, fieldName := range plan.argFields {
		f := v.Elem().FieldByName(fieldName)
		target := f.Addr().Interface()
		if conv != nil {
			scanner, ok := conv.FromDb(target)
			if ok {
				target = scanner.Holder
				custScan = append(custScan, scanner)
			}
		}
		dest[x] = target
	}

	row := exec.queryRow(plan.query, keys...)
	err = row.Scan(dest...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return nil, err
	}

	for _, c := range custScan {
		err = c.Bind()
		if err != nil {
			return nil, err
		}
	}

	if getChilds {
		// Get the primaty key for this table
		// Use the first PK found, multiple PKs are not supported
		// by now and will yield an error
		elem := v.Elem()
		PkId, _, err := m.GetPrimaryKey(*table, elem)
		if err != nil {
			return nil, errors.New("GetPrimaryKey failed in table '" + table.TableName + ": " + err.Error())
		}

		// Get child records if present - using the RelationMaps for this TableMap
		for _, r := range table.Relations {

			// Get the slice field in the master where the details will be stored into
			fv := elem.FieldByName(r.MasterFieldName)
			if !fv.IsValid() {
				return nil, errors.New("Field '" + r.MasterFieldName + "' not found in " + v.Kind().String() + " '" + elem.String() + "'")
			}
			if fv.Kind() == reflect.Slice {

				sql := fmt.Sprintf("select * from %s where %s = %d",
					m.Dialect.QuotedTableForQuery(table.SchemaName, r.DetailTable.TableName),
					m.Dialect.QuoteField(r.ForeignKeyFieldName), PkId)

				if (ChildLimit > -1) && (ChildOffset > -1) {
					sql = fmt.Sprintf(sql+" limit %d offset %d", ChildLimit, ChildOffset)
				} else {
					if (ChildLimit < 0) && (ChildOffset > 0) {
						ChildLimit = 999999999999999999
						sql = fmt.Sprintf(sql+" limit %d offset %d", ChildLimit, ChildOffset)
					}
				}

				_, err = m.Select(fv.Addr().Interface(), sql)
				if err != nil {
					return nil, errors.New("Get child relation " + r.DetailTable.TableName + " failed: " + err.Error())
				}
			} else {
				return nil, errors.New("Get child relation failed: Type " + fv.Type().Name() + " is not a slice")
			}
		}
	}

	if v, ok := v.Interface().(HasPostGet); ok {
		err := v.PostGet(exec)
		if err != nil {
			return nil, err
		}
	}

	return v.Interface(), nil
}

func delete(m *DbMap, exec SqlExecutor, list ...interface{}) (int64, error) {
	count := int64(0)
	for _, ptr := range list {
		table, elem, err := m.tableForPointer(ptr, true)
		if err != nil {
			return -1, err
		}

		eval := elem.Addr().Interface()
		if v, ok := eval.(HasPreDelete); ok {
			err = v.PreDelete(exec)
			if err != nil {
				return -1, err
			}
		}

		bi, err := table.bindDelete(elem)
		if err != nil {
			return -1, err
		}

		res, err := exec.Exec(bi.query, bi.args...)
		if err != nil {
			return -1, err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return -1, err
		}

		if rows == 0 && bi.existingVersion > 0 {
			return lockError(m, exec, table.TableName,
				bi.existingVersion, elem, bi.keys...)
		}

		count += rows

		if v, ok := eval.(HasPostDelete); ok {
			err := v.PostDelete(exec)
			if err != nil {
				return -1, err
			}
		}
	}

	return count, nil
}

func update(m *DbMap, exec SqlExecutor, updateChilds bool, list ...interface{}) (int64, error) {
	var table *TableMap
	var elem reflect.Value
	var err error

	count := int64(0)
	for _, ptr := range list {

		// Check if a pointer to reflect.Value has been passed
		if reflect.TypeOf(ptr).String() == "*reflect.Value" {
			// Indirect from Pointer to Value
			ptr = *ptr.(*reflect.Value)
		}

		// Check if a reflect.Value has been passed
		if reflect.TypeOf(ptr).String() == "reflect.Value" {
			elem = ptr.(reflect.Value)
			table, err = m.TableFor(elem.Type(), true)
			if err != nil {
				return -1, err
			}
		} else {

			table, elem, err = m.tableForPointer(ptr, true)
			if err != nil {
				return -1, err
			}
		}

		eval := elem.Addr().Interface()
		if v, ok := eval.(HasPreUpdate); ok {
			err = v.PreUpdate(exec)
			if err != nil {
				return -1, err
			}
		}
		// Get the primaty key for this table
		// Use the first PK found, multiple PKs are not supported
		// by now and will yield an error
		PkId, _, err := m.GetPrimaryKey(*table, elem)
		if err != nil {
			return count, errors.New("GetPrimaryKey failed: " + err.Error())
		}
		if m.DebugLevel > 2 {
			fmt.Printf("Update table %s with primarykey %d\n", table.TableName, PkId)
		}

		var bi bindInstance
		var rows int64
		if PkId == 0 {
			err = insert(m, exec, false, ptr)
			//bi, err = table.bindInsert(elem)
			if err != nil {
				return -1, err
			}
			rows = 1
			if m.DebugLevel > 2 {
				fmt.Printf("Update table %s has empty primary key, doing insert\n", table.TableName)
			}
		} else {
			bi, err = table.bindUpdate(elem)
			if err != nil {
				return -1, err
			}
			res, err := exec.Exec(bi.query, bi.args...)
			if err != nil {
				return -1, fmt.Errorf("gorp: update failed for table '%s': %s", table.TableName, err.Error())
			}
			rows, err = res.RowsAffected()
			if m.DebugLevel > 2 {
				fmt.Printf("Update RowsAffected %d\n", rows)
			}
			if err != nil {
				return -1, err
			}
		}

		if rows == 0 && bi.existingVersion > 0 {
			return lockError(m, exec, table.TableName,
				bi.existingVersion, elem, bi.keys...)
		}

		if bi.versField != "" {
			elem.FieldByName(bi.versField).SetInt(bi.existingVersion + 1)
		}

		count += rows
		// Store info about this update operation
		m.LastOpInfo.Type = Update
		m.LastOpInfo.BindPlanUsed = &table.updatePlan
		m.LastOpInfo.RowCount = count

		if updateChilds {
			// Get the primaty key for this table
			// Use the first PK found, multiple PKs are not supported
			// by now and will yield an error
			PkId, _, err = m.GetPrimaryKey(*table, elem)
			if err != nil {
				return count, errors.New("GetPrimaryKey failed: " + err.Error())
			}
			if PkId == 0 {
				return -1, errors.New(fmt.Sprintf("Update childs of table %s failed, primary key is zero\n", table.TableName))
			}
			//m.SetPrimaryKey(*table, elem, PkId)
			if m.DebugLevel > 2 {
				fmt.Printf("Update childs of table %s with primarykey %d\n", table.TableName, PkId)
			}
			// Update child records if present
			for _, r := range table.Relations {

				fv := elem.FieldByName(r.MasterFieldName)

				if fv.Kind() == reflect.Slice {
					updatecount, insertcount, err := m.UpdateDetailsFromSlice(elem, r, PkId)

					if err != nil {
						return count, errors.New("Update child relation on table '" + r.DetailTable.TableName + "' failed: " + err.Error())
					}
					m.LastOpInfo.ChildUpdateRowCount += updatecount
					m.LastOpInfo.ChildInsertRowCount += insertcount
				}
			}
		}

		if v, ok := eval.(HasPostUpdate); ok {
			err = v.PostUpdate(exec)
			if err != nil {
				return -1, err
			}
		}
	}
	return count, nil
}

// GetPrimaryKey returns the value(PkId) and the name (PkName) of a primary key from a table, if it exists
func (m *DbMap) GetPrimaryKey(table TableMap, elem reflect.Value) (PkId uint64, PkName string, err error) {
	// Get the primaty key for this table, multiple PKs are not supported by now
	for _, c := range table.Columns {
		if c.isPK {
			if PkName == "" {
				PkName = c.fieldName
			} else {
				err = fmt.Errorf("unsupported multiple primarykeys found in table '%s'", table.TableName)
				return
			}
		}
	}

	if PkName == "" {
		err = fmt.Errorf("No primary key found in table '%s'", table.TableName)
		return
	}
	// Get the value of the primary key
	f := elem.FieldByName(PkName)
	if !f.IsValid() {
		err = fmt.Errorf("Field '%s' not found in table '%s'", PkName, table.TableName)
		return
	}
	k := f.Kind()
	if (k == reflect.Int) || (k == reflect.Int16) || (k == reflect.Int32) || (k == reflect.Int64) {
		PkId = uint64(f.Int())
	} else if (k == reflect.Uint) || (k == reflect.Uint16) || (k == reflect.Uint32) || (k == reflect.Uint64) {
		PkId = f.Uint()
	} else {
		err = fmt.Errorf("Primary key '%s' in table '%s' is not of type int or uint", PkName, table.TableName)
		return
	}
	return
}

// SetPrimaryKey sets the primary key of a inmemory table if it exists
func (m *DbMap) SetPrimaryKey(table TableMap, elem reflect.Value, pk int64) (err error) {

	var PkName string
	// Get the primaty key for this table, multiple PKs are not supported by now
	for _, c := range table.Columns {
		if c.isPK {
			if PkName == "" {
				PkName = c.fieldName
			} else {
				err = fmt.Errorf("unsupported multiple primarykeys found in table '%s'", table.TableName)
				return
			}
		}
	}

	if PkName == "" {
		err = fmt.Errorf("No primary key found in table '%s'", table.TableName)
		return
	}
	// Get the value of the primary key
	f := elem.FieldByName(PkName)
	if !f.IsValid() {
		err = fmt.Errorf("Field '%s' not found in table '%s'", PkName, table.TableName)
		return
	}
	k := f.Kind()
	if (k == reflect.Int) || (k == reflect.Int16) || (k == reflect.Int32) || (k == reflect.Int64) {
		f.SetInt(pk)
	} else if (k == reflect.Uint) || (k == reflect.Uint16) || (k == reflect.Uint32) || (k == reflect.Uint64) {
		f.SetUint(uint64(pk))
	} else {
		err = fmt.Errorf("Primary key '%s' in table '%s' is not of type int or uint", PkName, table.TableName)
		return
	}
	return
}

func insert(m *DbMap, exec SqlExecutor, insertChilds bool, list ...interface{}) error {

	var table *TableMap
	var elem reflect.Value
	var err error

	for _, ptr := range list {

		// Check if a pointer to reflect.Value has been passed
		if reflect.TypeOf(ptr).String() == "*reflect.Value" {
			// Indirect from Pointer to Value
			ptr = *ptr.(*reflect.Value)
		}

		// Check if a reflect.Value has been passed
		if reflect.TypeOf(ptr).String() == "reflect.Value" {
			elem = ptr.(reflect.Value)
			table, err = m.TableFor(elem.Type(), true)
			if err != nil {
				return err
			}
		} else {
			table, elem, err = m.tableForPointer(ptr, false)
			if err != nil {
				return err
			}
		}

		eval := elem.Addr().Interface()
		if v, ok := eval.(HasPreInsert); ok {
			err := v.PreInsert(exec)
			if err != nil {
				return err
			}
		}

		bi, err := table.bindInsert(elem)
		if err != nil {
			return err
		}

		if bi.autoIncrIdx > -1 {
			f := elem.FieldByName(bi.autoIncrFieldName)
			switch inserter := m.Dialect.(type) {
			case IntegerAutoIncrInserter:
				id, err := inserter.InsertAutoIncr(exec, bi.query, bi.args...)
				if err != nil {
					return fmt.Errorf("gorp: insert failed for table '%s': %s", table.TableName, err.Error())
				}
				k := f.Kind()
				if (k == reflect.Int) || (k == reflect.Int16) || (k == reflect.Int32) || (k == reflect.Int64) {
					f.SetInt(id)
				} else if (k == reflect.Uint) || (k == reflect.Uint16) || (k == reflect.Uint32) || (k == reflect.Uint64) {
					f.SetUint(uint64(id))
				} else {
					return fmt.Errorf("gorp: Cannot set autoincrement value on non-Int field. SQL=%s  autoIncrIdx=%d autoIncrFieldName=%s", bi.query, bi.autoIncrIdx, bi.autoIncrFieldName)
				}
			case TargetedAutoIncrInserter:
				err := inserter.InsertAutoIncrToTarget(exec, bi.query, f.Addr().Interface(), bi.args...)
				if err != nil {
					return fmt.Errorf("gorp: insert failed for table '%s': %s", table.TableName, err.Error())
				}
			default:
				return fmt.Errorf("gorp: Cannot use autoincrement fields on dialects that do not implement an autoincrementing interface")
			}
		} else {
			_, err := exec.Exec(bi.query, bi.args...)
			if err != nil {
				return fmt.Errorf("gorp: Exec failed: %s", err.Error())
			}
		}

		// Store info about this update operation
		m.LastOpInfo.Type = Insert
		m.LastOpInfo.BindPlanUsed = &table.insertPlan
		m.LastOpInfo.RowCount++

		if insertChilds {
			// Get the primaty key for this table
			// Use the first PK found, multiple PKs are not supported
			// by now and will yield an error
			PkId, _, err := m.GetPrimaryKey(*table, elem)
			if err != nil {
				return errors.New("Insert child relation failed: " + err.Error())
			}

			// Insert child records if present
			for _, r := range table.Relations {

				fv := elem.FieldByName(r.MasterFieldName)

				if fv.Kind() == reflect.Slice {
					var count int64
					count, err = m.InsertDetailsFromSlice(elem, r, PkId)
					if err != nil {
						return errors.New("Insert child relation into table '" + r.DetailTable.TableName + "' failed: " + err.Error())

					}
					m.LastOpInfo.ChildInsertRowCount += count
				}

			}
		}

		if v, ok := eval.(HasPostInsert); ok {
			err := v.PostInsert(exec)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// InsertDetailsFromSlice inserts embedded structs described by the RelationMap r
// and sets the foreign key into each slice element from PK
// The master table is described by "elem"
func (m *DbMap) InsertDetailsFromSlice(elem reflect.Value, r *RelationMap, PK uint64) (int64, error) {

	var err error
	var count int64

	fv := elem.FieldByName(r.MasterFieldName)
	if fv.Kind() != reflect.Slice {
		return 0, fmt.Errorf("InsertFromSlice failed because Field '%s' in '%s' is not of type slice\n", r.MasterFieldName, elem.String())
	}

	// Buffer the rowcount info struct
	bufferLastOpInfo := m.LastOpInfo
	defer func() {
		m.LastOpInfo = bufferLastOpInfo
	}()

	for sliceIndex := 0; sliceIndex < fv.Len(); sliceIndex++ {

		fv0 := fv.Index(sliceIndex)
		// if the slice holds pointers, get the Element (do indirection)
		if fv0.Kind() == reflect.Ptr {
			fv0 = fv0.Elem()
		}

		// Get the foreign key name for this detail struct
		fd := fv0.FieldByName(r.ForeignKeyFieldName)

		// Set the foreign key of this detail
		if fd.Kind() == reflect.Uint32 || fd.Kind() == reflect.Uint64 {
			fd.SetUint(PK)
		} else if fd.Kind() == reflect.Int32 || fd.Kind() == reflect.Int64 {
			fd.SetInt(int64(PK))
		} else {
			return count, errors.New("InsertDetailsFromSlice failed: Unable to get ForeignKey: " + fd.Kind().String())
		}

		err = m.Insert(fv0)
		if err != nil {
			return count, errors.New("InsertDetailsFromSlice failed: " + err.Error())
		}
		count++
	}

	return count, err
}

// UpdateDetailsFromSlice updates embedded structs described by the RelationMap r in TableMap t
// If the primary key of the child/detail struct is set an update is executed, else an insert is done.
// The master table is the reflect value "elem"
func (m *DbMap) UpdateDetailsFromSlice(elem reflect.Value, r *RelationMap, PK uint64) (ucount int64, icount int64, err error) {
	var detailtable *TableMap

	fv := elem.FieldByName(r.MasterFieldName)
	if fv.Kind() != reflect.Slice {
		err = fmt.Errorf("UpdateDetailsFromSlice failed because Field '%s' in '%s' is not of type slice\n", r.MasterFieldName, elem.String())
		return
	}

	// Buffer the rowcount info struct
	bufferLastOpInfo := m.LastOpInfo
	defer func() {
		m.LastOpInfo = bufferLastOpInfo
	}()

	for sliceIndex := 0; sliceIndex < fv.Len(); sliceIndex++ {

		fv0 := fv.Index(sliceIndex)
		// if the slice holds pointers, get the Element (do indirection)
		if fv0.Kind() == reflect.Ptr {
			fv0 = fv0.Elem()
		}

		// Get the foreign key name for this detail struct
		fd := fv0.FieldByName(r.ForeignKeyFieldName)
		if !fd.IsValid() {
			err = fmt.Errorf("Field '%s' not found in '%s'", r.ForeignKeyFieldName, fv0.Type().String())
			return
		}

		detailtable, err = m.TableFor(fv0.Type(), true)
		if err != nil {
			err = fmt.Errorf("TableFor %s failed: %s", fv0.Type().String(), err.Error())
			return
		}

		// Get the primaty key for this detail table
		// Use the first PK found, multiple PKs are not supported
		// by now and will yield an error
		var detailPkId uint64
		detailPkId, _, err = m.GetPrimaryKey(*detailtable, fv0)
		if err != nil {
			err = errors.New("GetPrimaryKey in UpdateDetailsFromSlice failed: " + err.Error())
			return
		}

		// Check if the foreign key of the detail matches with the primary key of the master table
		if fd.Uint() == PK {
			if m.DebugLevel > 3 {
				log.Printf("r.ForeignKeyFieldName %s matches: %d, %d\n", r.ForeignKeyFieldName, fd.Uint(), PK)
			}
		} else {

			if m.DebugLevel > 3 {
				log.Printf("r.ForeignKeyFieldName %s does not match: %d, %d\n", r.ForeignKeyFieldName, fd.Uint(), PK)
			}
			// Set the foreign key of this detail
			if fd.Kind() == reflect.Uint32 || fd.Kind() == reflect.Uint64 {
				fd.SetUint(PK)
			} else if fd.Kind() == reflect.Int32 || fd.Kind() == reflect.Int64 {
				fd.SetInt(int64(PK))
			} else {
				err = errors.New(fmt.Sprintf("UpdateDetailsFromSlice failed: ForeignKey '%s' has incorrect type: '%s'", r.ForeignKeyFieldName, fd.Kind().String()))
				return
			}
			if m.DebugLevel > 3 {
				log.Printf("r.ForeignKeyFieldName %s has been updated to %d\n", r.ForeignKeyFieldName, fd.Uint())
			}
		}

		if detailPkId == 0 {

			err = m.Insert(fv0)

			if err != nil {
				err = errors.New("InsertFromSlice insert failed: " + err.Error())
				return
			}
			icount++

		} else {
			var affected int64
			affected, err = m.Update(fv0)
			if err != nil {
				err = errors.New(fmt.Sprintf("UpdateDetailsFromSlice failed for detailPkId: %d: %s", detailPkId, err.Error()))
				return
			}
			if (affected == 0) && (m.CheckAffectedRows == true) {
				err = errors.New(fmt.Sprintf("UpdateDetailsFromSlice affected 0 records for detailPkId: %d\n", detailPkId))
				return
			}
			ucount += affected
		}

	}

	return
}

func checkForNotNull(elem reflect.Value, col *ColumnMap, table *TableMap) (err error) {
	var isNull bool
	if col.EnforceNotNull && col.isNotNull {
		val := elem.FieldByName(col.fieldName)
		if val.Kind() == reflect.String {
			valstring := val.String()
			// DEBUG
			if valstring == "" {
				isNull = true
			}
		} else if val.Kind() == reflect.Int || val.Kind() == reflect.Int8 || val.Kind() == reflect.Int16 ||
			val.Kind() == reflect.Int32 || val.Kind() == reflect.Int64 {
			valint := val.Int()
			if valint == 0 {
				isNull = true
			}
		} else if val.Kind() == reflect.Uint || val.Kind() == reflect.Uint8 || val.Kind() == reflect.Uint16 ||
			val.Kind() == reflect.Uint32 || val.Kind() == reflect.Uint64 {
			valint := val.Uint()
			if valint == 0 {
				isNull = true
			}
		}
	}
	if isNull {
		err = errors.New(fmt.Sprintf("Trying to insert a zero value into field '%s.%s', which is NOT NULL and EnforceNotNull is true", table.TableName, col.ColumnName))
	}
	return
}

func lockError(m *DbMap, exec SqlExecutor, tableName string,
	existingVer int64, elem reflect.Value,
	keys ...interface{}) (int64, error) {

	existing, err := get(m, exec, elem.Interface(), false, 0, 0, keys...)
	if err != nil {
		return -1, err
	}

	ole := OptimisticLockError{tableName, keys, true, existingVer}
	if existing == nil {
		ole.RowExists = false
	}
	return -1, ole
}

// PostUpdate() will be executed after the GET statement.
type HasPostGet interface {
	PostGet(SqlExecutor) error
}

// PostUpdate() will be executed after the DELETE statement
type HasPostDelete interface {
	PostDelete(SqlExecutor) error
}

// PostUpdate() will be executed after the UPDATE statement
type HasPostUpdate interface {
	PostUpdate(SqlExecutor) error
}

// PostInsert() will be executed after the INSERT statement
type HasPostInsert interface {
	PostInsert(SqlExecutor) error
}

// PreDelete() will be executed before the DELETE statement.
type HasPreDelete interface {
	PreDelete(SqlExecutor) error
}

// PreUpdate() will be executed before UPDATE statement.
type HasPreUpdate interface {
	PreUpdate(SqlExecutor) error
}

// PreInsert() will be executed before INSERT statement.
type HasPreInsert interface {
	PreInsert(SqlExecutor) error
}