# sqlite

Support SQLite database for [github.com/gopsql/psql](https://github.com/gopsql/psql).

This package uses [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite), a CGo-free port of SQLite.

## Example

```go
package main

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/gopsql/logger"
	"github.com/gopsql/psql"
	"github.com/gopsql/sqlite"
	"github.com/shopspring/decimal"
)

type Page struct {
	Id        int
	Foo       int64
	Bar       bool
	Title     string
	Cats      strArray
	Price     decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
}

type strArray []string

func (c *strArray) Scan(src interface{}) error {
	if value, ok := src.(string); ok {
		*c = strings.Split(strings.TrimSpace(value), "\n")
	}
	return nil
}

func (c strArray) Value() (driver.Value, error) {
	return strings.Join(c, "\n") + "\n", nil
}

func main() {
	conn := sqlite.MustOpen("mydbfile.sqlite3")
	m := psql.NewModel(Page{}, conn, logger.StandardLogger)
	m.NewSQL(m.Schema()).MustExecute()
	// CREATE TABLE Pages (
	//         Id INTEGER PRIMARY KEY AUTOINCREMENT,
	//         Foo bigint DEFAULT 0 NOT NULL,
	//         Bar boolean DEFAULT false NOT NULL,
	//         Title text DEFAULT '' NOT NULL,
	//         Cats text DEFAULT '' NOT NULL,
	//         Price decimal(10, 2) DEFAULT 0.0 NOT NULL,
	//         CreatedAt timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	//         UpdatedAt timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL
	// );
	m.Insert(
		"Title", "Hello",
		"Price", decimal.RequireFromString("1.2345"),
		"Cats", strArray{"foobar", "hello"},
		"CreatedAt", time.Now(),
	).MustExecute()
	var pages []Page
	m.Find().Where("Cats = $1", strArray{"foobar", "hello"}).MustQuery(&pages)
	fmt.Println(pages)
	conn.Close()
}
```
