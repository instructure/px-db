package display

import (
	"github.com/spf13/cobra"
)

// DisplayOptions command options to display data from PostgreSQL
type DisplayOptions struct {
}

var (
	displayShort   = "Display data from a PostgreSQL DB"
	displayExample = "usage: px-db display [table|some-command]"
)

// NewDisplayCmd Display data from a PostgreSQL DB
func NewDisplayCmd() *cobra.Command {
	options := DisplayOptions{}
	cmd := &cobra.Command{
		Use:     "display",
		Args:    cobra.MinimumNArgs(1),
		Short:   displayShort,
		Example: displayExample,
	}

	cmd.AddCommand(NewDisplayTableCmd(options))
	return cmd
}
