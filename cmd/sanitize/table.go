package sanitize

import (
	"strings"

	"github.com/instructure/px-db/util/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	tablesDeleteShort   = "Sanitize Specific Tables via Deletion for PostgreSQL DB"
	tablesDeleteExample = "usage: px-db sanitize delete-tables --db-tables table1,table2,table3"
)

// DeleteTableOptions Control options flags around table sanitization
type DeleteTableOptions struct {
	*SanitizeOptions
	cascadeDelete bool
	dbTables      string
}

// NewDeleteTablesCmd Sanitize a list of tables
func NewDeleteTablesCmd(sanitizeOptions SanitizeOptions) *cobra.Command {
	tableOptions := DeleteTableOptions{
		SanitizeOptions: &sanitizeOptions,
	}

	cmd := &cobra.Command{
		TraverseChildren: true,
		SilenceUsage:     true,
		Use:              "delete-tables",
		Short:            tablesDeleteShort,
		Example:          tablesDeleteExample,
		RunE:             deleteTables,
	}

	cmd.PersistentFlags().BoolVar(&tableOptions.cascadeDelete,
		"cascade-mode", false, "Certain Tables Might need to be grouped for removal and include a CASCADE for TRUNCATE")
	cmd.PersistentFlags().StringVar(&tableOptions.dbTables, "db-tables", "", "Comma separated list of Tables to be wiped (e.g. table1,table2,table3)")

	return cmd
}

func deleteTables(cmd *cobra.Command, args []string) error {
	logContext := "[Sanitize Tables Deletion]"
	log.Infof("%s Sanitizing specific tables via deletion...", logContext)

	port, _ := cmd.Flags().GetInt64("db-port")
	isSSL, _ := cmd.Flags().GetBool("db-ssl-mode")
	isCascade, _ := cmd.Flags().GetBool("cascade-mode")

	dbConn, err := pq.NewDBConnection(&pq.DBConnectionOptions{
		Endpoint: cmd.Flag("db-endpoint").Value.String(),
		Name:     cmd.Flag("db-name").Value.String(),
		Password: viper.GetString("DB_PASSWORD"),
		Port:     port,
		SSLMode:  isSSL,
		User:     cmd.Flag("db-user").Value.String(),
	})
	if err != nil {
		return errors.Wrap(err, logContext)
	}

	tableFmt, _ := cmd.Flags().GetString("db-tables")
	tables := strings.Split(tableFmt, ",")
	log.Debugf("%s Deleting Table contents: %v", logContext, tables)
	log.Debugf("Length: %d Value: %s", len(tables), tables[0])
	if (len(tables)-1) >= 0 && tables[0] == "" {
		return errors.Wrap(errors.New("Please use flag --db-tables to pass in a comma delimited list of tables (e.g. table1,table2,table3)"), logContext)
	}

	// Alter all tables first
	for _, t := range tables {
		if err := pq.DropConstraints(dbConn, &pq.BaseAlterTableProperties{Table: t}); err != nil {
			return errors.Wrap(err, logContext)
		}
	}

	for _, t := range tables {
		log.Infof("%s Deleting Tables: %v", logContext, t)
		if err := pq.DeleteTable(dbConn, &pq.BaseAlterTableProperties{Table: t}, isCascade); err != nil {
			return errors.Wrap(err, logContext)
		}
	}

	return nil
}

/*
	rows, err := dbConn.Query("SELECT * FROM \"Asset\"")
	defer rows.Close()
	if err != nil {
		return fmt.Errorf("Query Error: %v", err)
	}

	cols, _ := rows.Columns()
	if err := pq.ParseRow(rows, cols); err != nil {
		log.Error(err)
	}*/
