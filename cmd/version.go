package cmd

import "github.com/spf13/cobra"

var (
	version = "0.8.1"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			RootCmd.Printf("%s version %s\n", appName, version)
		},
	}
)
