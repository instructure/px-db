package display

import (
	"github.com/instructure/px-db/util/pq"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var (
	displayTableShort   = "Display data from a Table in PostgreSQL DB"
	displayTableExample = "usage: px-db display table --name my_table_name"
)

// DisplayTableOptions command options to display data in PostgreSQL
type DisplayTableOptions struct {
	*DisplayOptions
	tableName string
}

// NewDisplayTableCmd Display data from a PostgreSQL DB
func NewDisplayTableCmd(displayOptions DisplayOptions) *cobra.Command {
	options := DisplayTableOptions{
		DisplayOptions: &displayOptions,
	}

	cmd := &cobra.Command{
		TraverseChildren: true,
		SilenceUsage:     true,
		Use:              "table",
		Short:            displayTableShort,
		Example:          displayTableExample,
		RunE:             displayTable,
	}

	cmd.PersistentFlags().StringVar(&options.tableName, "name", "", "A PostgreSQL Table to retrieve data from and print (e.g. table1)")

	return cmd
}

func displayTable(cmd *cobra.Command, args []string) error {
	logContext := "[Display Table]"
	log.Infof("%s Printing Table data out to STDOUT...", logContext)

	port, _ := cmd.Flags().GetInt64("db-port")
	isSSL, _ := cmd.Flags().GetBool("db-ssl-mode")

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

	table, _ := cmd.Flags().GetString("name")
	log.Debugf("%s Printing Table contents: %v", logContext, table)
	if table == "" {
		return errors.Wrap(errors.New("Please use flag --name to pass in a table name (e.g. table1)"), logContext)
	}

	if err := pq.DisplayRows(dbConn, table); err != nil {
		return errors.Wrap(err, logContext)
	}

	return nil
}
