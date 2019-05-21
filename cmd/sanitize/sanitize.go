package sanitize

import (
	"github.com/spf13/cobra"
)

// SanitizeOptions command options for sanitize
type SanitizeOptions struct {
}

var (
	sanitizeShort   = "Sanitize a PostgreSQL DB"
	sanitizeExample = "usage: px-db sanitize [tables|some command]"
)

// NewSanitizeCmd Sanitize a PostgreSQL DB
func NewSanitizeCmd() *cobra.Command {
	options := SanitizeOptions{}
	cmd := &cobra.Command{
		Use:     "sanitize",
		Args:    cobra.MinimumNArgs(1),
		Short:   sanitizeShort,
		Example: sanitizeExample,
	}

	cmd.AddCommand(NewDeleteTablesCmd(options))
	return cmd
}
