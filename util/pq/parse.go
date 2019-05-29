package pq

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// DisplayRows Parse PostgreSQL Rows in a dynamic fashion
func DisplayRows(dbConn *sql.DB, t string) error {

	displayQuery := fmt.Sprintf("SELECT * FROM \"%s\"", t)
	log.Debugf(displayQuery)

	rows, err := dbConn.Query(displayQuery)
	if err != nil {
		return err
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		fmt.Println(m)
	}
	return nil
}
