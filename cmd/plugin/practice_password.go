package plugin

import (
	"github.com/spf13/cobra"
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
	return nil
}
