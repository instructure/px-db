package cmd

import (
	"os"

	"github.com/instructure/px-db/cmd/plugin"
	"github.com/instructure/px-db/cmd/sanitize"
	"github.com/instructure/px-db/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RootOptions struct {
	configFile string
	dbName     string
	dbUser     string
	dbEndpoint string
	dbSSL      bool
}

func initConfig(options RootOptions) {
	//cfgFile := options.configFile

	viper.AutomaticEnv()
	viper.SetDefault("DB_PASSWORD", "practice")
	/*
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			viper.AddConfigPath(home)
			viper.SetConfigName("config")
		}

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}*/
}

func NewCmdRoot() *cobra.Command {

	options := RootOptions{}
	initConfig(options)

	cmd := &cobra.Command{
		Use:   "px-db",
		Short: "px-db is for sanitizing Practice PostgreSQL Tables",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Usage()
				os.Exit(0)
			}

			if args[0] == "help" || args[0] == "--help" {
				cmd.Usage()
				os.Exit(0)
			}
		},
		Version: version.Print(),
	}

	cmd.PersistentFlags().StringVar(&options.configFile, "config", "", "path to config file (default $HOME/.px-db-sanitizer/config)")
	cmd.PersistentFlags().StringVar(&options.dbEndpoint, "db-endpoint", "", "PostgreSQL Connection String")
	cmd.PersistentFlags().StringVar(&options.dbName, "db-name", "", "PostgreSQL Connection String")
	cmd.PersistentFlags().BoolVar(&options.dbSSL, "db-ssl-mode", true, "PostgreSQL Connection String")
	cmd.PersistentFlags().StringVar(&options.dbUser, "db-user", "", "PostgreSQL Connection String")

	// create subcommands
	cmd.AddCommand(sanitize.NewSanitizeCmd())
	cmd.AddCommand(plugin.NewPluginCmd())
	return cmd
}
