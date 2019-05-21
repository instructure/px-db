package pq

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// DBConnectionOptions PostgreSQL Connection Options
type DBConnectionOptions struct {
	Endpoint string
	Name     string
	Password string
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

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", d.User, d.Password, d.Endpoint, d.Name, sslmode)
	log.Debugf("DB Conn String: postgres://%s:<sanitized_password>@%s/%s?sslmode=%s", d.User, d.Endpoint, d.Name, sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error connection to Database: %v", err))
	}

	return db, nil
}
