package pq

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// BaseAlterTableProperties <TODO>
type BaseAlterTableProperties struct {
	Table string
}

// UpdateTableByColumnProperties <TODO>
type UpdateTableByColumnProperties struct {
	BaseAlterTableProperties
	Column string
	NewVal string
}

// UpdateTableByColumnUniqueProperties <TODO>
type UpdateTableByColumnUniqueProperties struct {
	BaseAlterTableProperties
	Column            string
	IncrementByColumn string
	NewValPrefix      string
}

// UpdateTableByColumnUniqueFmtProperties <TODO>
type UpdateTableByColumnUniqueFmtProperties struct {
	UpdateTableByColumnUniqueProperties
	NewValSuffix string
}

// AlterTable <TODO>
func alterTableConstraint(dbConn *sql.DB, t string, constraintName string) error {
	alterQuery := fmt.Sprintf("ALTER TABLE \"%s\" DROP CONSTRAINT IF EXISTS \"%s\" CASCADE", t, constraintName)
	log.Debugf(alterQuery)
	_, err := dbConn.Query(alterQuery)
	if err != nil {
		log.Error(err)
	}

	return nil
}

// DropConstraints <TODO>
func DropConstraints(dbConn *sql.DB, props *BaseAlterTableProperties) error {
	t := props.Table
	constraintQuery := fmt.Sprintf("ALTER TABLE \"%s\" DISABLE trigger ALL;"+
		"SELECT constraint_name FROM information_schema.constraint_table_usage WHERE table_name = '%s'", t, t)
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
func DeleteTable(dbConn *sql.DB, props *BaseAlterTableProperties, isCascade bool) error {
	t := props.Table
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

// UpdateTableByColumn update a column for all rows in a table with the same new value
func UpdateTableByColumn(dbConn *sql.DB, props *UpdateTableByColumnProperties) error {
	t := props.Table
	col := props.Column
	val := props.NewVal

	updateQuery := fmt.Sprintf("UPDATE \"%s\" SET \"%s\" = '%s'", t, col, val)
	log.Debugf(updateQuery)
	_, err := dbConn.Query(updateQuery)
	if err != nil {
		return fmt.Errorf("Unable to Update Table Column in %s error for %s: %v", t, col, err)
	}

	return nil
}

// UpdateTableByColumnUnique increments onto a new prefix value for a column in each row within a table
// Leverages an IncrementByColumn to append a unique row value for that column onto the NewValPrefix
func UpdateTableByColumnUnique(dbConn *sql.DB, props *UpdateTableByColumnUniqueProperties) error {
	t := props.Table
	col := props.Column
	newValPrefix := props.NewValPrefix
	incrementColumn := props.IncrementByColumn

	updateQuery := fmt.Sprintf("UPDATE \"%s\" t "+
		"SET \"%s\" = '%s' || (SELECT \"%s\" FROM \"%s\" WHERE \"%s\" = t.\"%s\")",
		t, col, newValPrefix, incrementColumn, t, incrementColumn, incrementColumn)
	log.Debugf(updateQuery)
	_, err := dbConn.Query(updateQuery)
	if err != nil {
		return fmt.Errorf("Unable to Update Table Column in %s error for %s: %v", t, col, err)
	}

	return nil
}

// UpdateTableByColumnUniqueFmt increments and formats a new unique value via a prefix & suffix value for a column in each row within a table
// Leverages an IncrementByColumn to enclose a unique value in between the supplied NewValPrefix and NewValSuffix
func UpdateTableByColumnUniqueFmt(dbConn *sql.DB, props *UpdateTableByColumnUniqueFmtProperties) error {
	t := props.Table
	col := props.Column
	newValPrefix := props.NewValPrefix
	newValSuffix := props.NewValSuffix
	incrementColumn := props.IncrementByColumn

	updateQuery := fmt.Sprintf("UPDATE \"%s\" t "+
		"SET \"%s\" = '%s' || (SELECT \"%s\" FROM \"%s\" WHERE \"%s\" = t.\"%s\") || '%s'",
		t, col, newValPrefix, incrementColumn, t, incrementColumn, incrementColumn, newValSuffix)
	log.Debugf(updateQuery)
	_, err := dbConn.Query(updateQuery)
	if err != nil {
		return fmt.Errorf("Unable to Update Table Column in %s error for %s: %v", t, col, err)
	}

	return nil
}
