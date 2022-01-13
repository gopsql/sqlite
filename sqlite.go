package sqlite

import (
	"database/sql"
	"strings"

	"github.com/gopsql/db"
	"github.com/gopsql/standard"
	_ "modernc.org/sqlite"
)

type (
	sqliteDB struct {
		standard.DB
	}
)

func (d *sqliteDB) FieldDataType(fieldName, fieldType string) (dataType string) {
	return FieldDataType(fieldName, fieldType)
}

var _ db.DB = (*sqliteDB)(nil)

// MustOpen is like Open but panics if connect operation fails.
func MustOpen(conn string) *sqliteDB {
	c, err := Open(conn)
	if err != nil {
		panic(err)
	}
	return c
}

// Open creates and establishes one connection to database.
func Open(conn string) (*sqliteDB, error) {
	c, err := sql.Open("sqlite", conn)
	if err != nil {
		return nil, err
	}
	if err := c.Ping(); err != nil {
		return nil, err
	}
	return &sqliteDB{standard.DB{c}}, nil
}

// Generate data type based on struct's field name and type.
func FieldDataType(fieldName, fieldType string) (dataType string) {
	if strings.ToLower(fieldName) == "id" && strings.Contains(fieldType, "int") {
		dataType = "INTEGER PRIMARY KEY AUTOINCREMENT"
		return
	}
	var null bool
	if strings.HasPrefix(fieldType, "*") {
		fieldType = strings.TrimPrefix(fieldType, "*")
		null = true
	}
	var defValue string
	switch fieldType {
	case "int8", "int16", "int32", "uint8", "uint16", "uint32":
		dataType = "integer"
		defValue = "0"
	case "int64", "uint64", "int", "uint":
		dataType = "bigint"
		defValue = "0"
	case "time.Time":
		dataType = "timestamp"
		defValue = "CURRENT_TIMESTAMP"
	case "float32", "float64":
		dataType = "decimal(10, 2)"
		defValue = "0.0"
	case "decimal.Decimal":
		dataType = "decimal(10, 2)"
		defValue = "0.0"
	case "bool":
		dataType = "boolean"
		defValue = "false"
	default:
		dataType = "text"
		defValue = "''"
	}
	dataType += " DEFAULT " + defValue
	if !null {
		dataType += " NOT NULL"
	}
	return
}
