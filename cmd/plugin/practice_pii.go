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
	practicePIIShort   = "Custom logic for updating OAuthClient table with 'development' defaults and Sanitizes User PII in the Users Table"
	practicePIIExample = "usage: px-db plugin practice-pii"
)

// PracticePIIOptions Control options flags around password table updates
type PracticePIIOptions struct {
	*PluginOptions
}

// NewPracticePIICmd Updates OAuthClient table with 'development' defaults and Sanitizes User PIIs in the Users Table
func NewPracticePIICmd(pluginOptions PluginOptions) *cobra.Command {
	/*
		tableOptions := DeleteTableOptions{
			SanitizeOptions: &sanitizeOptions,
		}*/

	cmd := &cobra.Command{
		TraverseChildren: true,
		SilenceUsage:     true,
		Use:              "practice-pii",
		Short:            practicePIIShort,
		Example:          practicePIIExample,
		RunE:             practicePII,
	}

	return cmd
}

func practicePII(cmd *cobra.Command, args []string) error {
	logContext := "[Plugins Practice PII]"
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
	baseOAuthProps := pq.BaseAlterTableProperties{Table: "OAuthClient"}
	if err := pq.UpdateTableByColumn(dbConn, &pq.UpdateTableByColumnProperties{
		BaseAlterTableProperties: baseOAuthProps,
		Column:                   "secretHash",
		NewVal:                   oauthPassword,
	}); err != nil {
		return errors.Wrap(err, logContext)
	}

	if err := pq.UpdateTableByColumnUnique(dbConn, &pq.UpdateTableByColumnUniqueProperties{
		BaseAlterTableProperties: baseOAuthProps,
		Column:                   "key",
		NewValPrefix:             password.APIKey,
		IncrementByColumn:        "id",
	}); err != nil {
		return errors.Wrap(err, logContext)
	}

	// Update the Users' Passwords in User Table
	baseUserProps := pq.BaseAlterTableProperties{Table: "User"}
	if err := pq.UpdateTableByColumn(dbConn, &pq.UpdateTableByColumnProperties{
		BaseAlterTableProperties: baseUserProps,
		Column:                   "passwordHash",
		NewVal:                   userPassword,
	}); err != nil {
		return errors.Wrap(err, logContext)
	}

	if err := pq.UpdateTableByColumnUnique(dbConn, &pq.UpdateTableByColumnUniqueProperties{
		BaseAlterTableProperties: baseUserProps,
		Column:                   "name",
		NewValPrefix:             "px-automation",
		IncrementByColumn:        "id",
	}); err != nil {
		return errors.Wrap(err, logContext)
	}

	// Update the Users' emails in the UserEmail Table
	baseUserEmailProps := pq.BaseAlterTableProperties{Table: "UserEmail"}
	if err := pq.UpdateTableByColumnUniqueFmt(dbConn, &pq.UpdateTableByColumnUniqueFmtProperties{
		UpdateTableByColumnUniqueProperties: pq.UpdateTableByColumnUniqueProperties{
			BaseAlterTableProperties: baseUserEmailProps,
			Column:                   "address",
			IncrementByColumn:        "id",
			NewValPrefix:             "px-automation+",
		},
		NewValSuffix: "@instructure.com",
	}); err != nil {
		return errors.Wrap(err, logContext)
	}

	return nil
}
