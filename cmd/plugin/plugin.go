package plugin

import (
	"github.com/spf13/cobra"
)

// PluginOptions command options for sanitize
type PluginOptions struct {
}

var (
	pluginShort   = "Run Plugins that perform custom logic for PostgreSQL DB Table Sanitization"
	pluginExample = "usage: px-db plugin [practice-pii|some-plugin]"
)

// NewPluginCmd Sanitize a PostgreSQL DB
func NewPluginCmd() *cobra.Command {
	options := PluginOptions{}
	cmd := &cobra.Command{
		Use:     "plugin",
		Args:    cobra.MinimumNArgs(1),
		Short:   pluginShort,
		Example: pluginExample,
	}

	cmd.AddCommand(NewPracticePIICmd(options))
	return cmd
}
