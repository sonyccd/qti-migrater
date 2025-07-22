package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	verbosity int
)

var rootCmd = &cobra.Command{
	Use:   "qti-migrator",
	Short: "QTI file migration tool",
	Long: `QTI Migrater is a CLI tool for migrating QTI (Question and Test Interoperability) 
files between different versions. Currently supports migration from QTI 1.2 to 2.1, 
with future support planned for QTI 2.1 to 3.0 (JSON format).`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.qti-migrator.yaml)")
	rootCmd.PersistentFlags().IntVarP(&verbosity, "verbosity", "v", 1, "verbosity level (0=minimal, 1=normal, 2=detailed, 3=debug)")

	if err := viper.BindPFlag("verbosity", rootCmd.PersistentFlags().Lookup("verbosity")); err != nil {
		panic(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName(".qti-migrator")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if verbosity >= 2 {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}