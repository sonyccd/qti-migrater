package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of qti-migrator",
	Long:  `All software has versions. This is qti-migrator's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("qti-migrator v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}