package pq

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// AlterTable <TODO>
func alterTableConstraint(dbConn *sql.DB, t string, c string) error {
	alterQuery := fmt.Sprintf("ALTER TABLE \"%s\" DROP CONSTRAINT IF EXISTS \"%s\" CASCADE", t, c)
	log.Debugf(alterQuery)
	_, err := dbConn.Query(alterQuery)
	if err != nil {
		log.Error(err)
	}

	return nil
}

// DropConstraints <TODO>
func DropConstraints(dbConn *sql.DB, t string) error {
	constraintQuery := fmt.Sprintf("ALTER TABLE \"%s\" DISABLE trigger ALL;SELECT constraint_name FROM information_schema.constraint_table_usage WHERE table_name = '%s'", t, t)
	log.Debugf(constraintQuery)
	rows, err := dbConn.Query(constraintQuery)
	if err != nil {
		return fmt.Errorf("Failed to query Constraints: %s - %v", t, err)
	}

	for rows.Next() {
		var constraintName string

		if err := rows.Scan(&constraintName); err != nil {
			return fmt.Errorf("Unable to get contraints for Table: %s - %v", t, err)
		}

		if err := alterTableConstraint(dbConn, t, constraintName); err != nil {
			return fmt.Errorf("Unable to drop contraints for Table: %s - %v", t, err)
		}
	}

	return nil
}

// DeleteTable <TODO>
func DeleteTable(dbConn *sql.DB, t string, isCascade bool) error {
	var deleteTableQuery string
	if isCascade {
		deleteTableQuery = fmt.Sprintf("TRUNCATE \"%v\" CASCADE", t)
	} else {
		deleteTableQuery = fmt.Sprintf("TRUNCATE \"%v\"", t)
	}
	log.Debugf(deleteTableQuery)
	_, err := dbConn.Query(deleteTableQuery)
	if err != nil {
		return fmt.Errorf("Table Deletion error for %s: %v", t, err)
	}

	return nil
}

// UpdateAllTableColumn update all row values in a column
func UpdateAllTableColumn(dbConn *sql.DB, t string, col string, val string) error {
	updateQuery := fmt.Sprintf("UPDATE \"%s\" SET \"%s\" = '%s'", t, col, val)
	log.Debugf(updateQuery)
	_, err := dbConn.Query(updateQuery)
	if err != nil {
		return fmt.Errorf("Unable to Update Table Column in %s error for %s: %v", t, col, err)
	}

	return nil
}

// UpdateRowTableColumn increments a value for each row in a column
func UpdateRowTableColumn(dbConn *sql.DB, t string, col string, val string) error {
	var count int64
	count++

	updateQuery := fmt.Sprintf("UPDATE \"%s\" SET \"%s\" = '%s'", t, col, val)
	log.Debugf(updateQuery)
	_, err := dbConn.Query(updateQuery)
	if err != nil {
		return fmt.Errorf("Unable to Update Table Column in %s error for %s: %v", t, col, err)
	}

	return nil
}
