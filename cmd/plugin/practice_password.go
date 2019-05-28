package plugin

import (
	"github.com/instructure/px-db/plugins/password"
	"github.com/instructure/px-db/util/pq"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var (
	practicePasswordShort   = "Updates OAuthClient table with 'development' defaults and Sanitizes User Passwords in the Users Table"
	practicePasswordExample = "usage: px-db plugin practice-password"
)

// PracticePasswordOptions Control options flags around password table updates
type PracticePasswordOptions struct {
	*PluginOptions
}

// NewPracticePasswordCmd Updates OAuthClient table with 'development' defaults and Sanitizes User Passwords in the Users Table
func NewPracticePasswordCmd(pluginOptions PluginOptions) *cobra.Command {
	/*
		tableOptions := DeleteTableOptions{
			SanitizeOptions: &sanitizeOptions,
		}*/

	cmd := &cobra.Command{
		TraverseChildren: true,
		SilenceUsage:     true,
		Use:              "practice-password",
		Short:            practicePasswordShort,
		Example:          practicePasswordExample,
		RunE:             practicePassword,
	}

	return cmd
}

func practicePassword(cmd *cobra.Command, args []string) error {
	logContext := "[Plugins Practice Password]"
	log.Infof("%s Sanitizing OAuth and Users tables (passwords, emails, and OAuth Keys)...", logContext)

	dbConn, err := pq.NewDBConnection(&pq.DBConnectionOptions{
		Endpoint: cmd.Flag("db-endpoint").Value.String(),
		Name:     cmd.Flag("db-name").Value.String(),
		Password: viper.GetString("DB_PASSWORD"),
		SSLMode:  viper.GetBool("db-ssl-mode"),
		User:     cmd.Flag("db-user").Value.String(),
	})
	if err != nil {
		return errors.Wrap(err, logContext)
	}

	oauthPassword, err := password.Retrieve(&password.OAuth{})
	if err != nil {
		return errors.Wrap(err, logContext)
	}
	userPassword, err := password.Retrieve(&password.User{})
	if err != nil {
		return errors.Wrap(err, logContext)
	}

	log.Infof("Password: %s and %s key: %s", userPassword, oauthPassword, password.APIKey)

	// Update the OAuth Table
	oauthTable := "OAuthClient"
	if err := pq.UpdateAllTableColumn(dbConn, oauthTable, "secretHash", oauthPassword); err != nil {
		return errors.Wrap(err, logContext)
	}

	if err := pq.UpdateAllTableColumn(dbConn, oauthTable, "key", password.APIKey); err != nil {
		return errors.Wrap(err, logContext)
	}

	// Update the Users Passwords in User Table
	userTable := "User"
	if err := pq.UpdateAllTableColumn(dbConn, userTable, "passwordHash", userPassword); err != nil {
		return errors.Wrap(err, logContext)
	}

	return nil
}
