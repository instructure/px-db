package pq

import (
	"database/sql"
	"fmt"

	// To avoid name collision
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// DBConnectionOptions PostgreSQL Connection Options
type DBConnectionOptions struct {
	Endpoint string
	Name     string
	Password string
	Port     int64
	User     string
	SSLMode  bool
}

// NewDBConnection create a postgres connection
func NewDBConnection(d *DBConnectionOptions) (*sql.DB, error) {
	var sslmode string

	if d.SSLMode {
		sslmode = "verify-full"
	} else {
		sslmode = "disable"
	}

	connParams := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		d.Endpoint, d.Port, d.User, d.Password, d.Name, sslmode)
	connParamsSanitized := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=<sanitized> dbname=%s sslmode=%s",
		d.Endpoint, d.Port, d.User, d.Name, sslmode)

	log.Debugf(connParamsSanitized)
	db, err := sql.Open("postgres", connParams)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error connection to Database: %v", err))
	}

	return db, nil
}
