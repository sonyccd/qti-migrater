package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/qti-migrator/internal/migrator"
	"github.com/qti-migrator/internal/preprocessor"
	"github.com/qti-migrator/internal/report"
)

var (
	inputFile    string
	outputFile   string
	fromVersion  string
	toVersion    string
	previewOnly  bool
	forceOverwrite bool
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate QTI files between versions",
	Long:  `Migrate QTI files from one version to another. Supports migration from QTI 1.2 to 2.1.`,
	RunE:  runMigrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().StringVarP(&inputFile, "input", "i", "-", "Input file path (use '-' for stdin)")
	migrateCmd.Flags().StringVarP(&outputFile, "output", "o", "-", "Output file path (use '-' for stdout)")
	migrateCmd.Flags().StringVarP(&fromVersion, "from", "f", "", "Source QTI version (1.2, 2.1)")
	migrateCmd.Flags().StringVarP(&toVersion, "to", "t", "", "Target QTI version (2.1, 3.0)")
	migrateCmd.Flags().BoolVarP(&previewOnly, "preview", "p", false, "Preview migration without executing")
	migrateCmd.Flags().BoolVarP(&forceOverwrite, "force", "", false, "Force overwrite output file if it exists")

	migrateCmd.MarkFlagRequired("from")
	migrateCmd.MarkFlagRequired("to")
}

func runMigrate(cmd *cobra.Command, args []string) error {
	var input io.Reader
	var output io.Writer
	var err error

	if inputFile == "-" {
		input = os.Stdin
	} else {
		file, err := os.Open(inputFile)
		if err != nil {
			return fmt.Errorf("error opening input file: %w", err)
		}
		defer file.Close()
		input = file
	}

	if !previewOnly {
		if outputFile == "-" {
			output = os.Stdout
		} else {
			if !forceOverwrite {
				if _, err := os.Stat(outputFile); err == nil {
					return fmt.Errorf("output file already exists: %s (use --force to overwrite)", outputFile)
				}
			}
			file, err := os.Create(outputFile)
			if err != nil {
				return fmt.Errorf("error creating output file: %w", err)
			}
			defer file.Close()
			output = file
		}
	}

	content, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	processor := preprocessor.New(verbosity)
	analysisReport, err := processor.Analyze(content, fromVersion, toVersion)
	if err != nil {
		return fmt.Errorf("error analyzing file: %w", err)
	}

	reporter := report.New(verbosity)
	reportOutput := reporter.Generate(analysisReport)
	
	if verbosity >= 1 || previewOnly {
		fmt.Fprintln(os.Stderr, reportOutput)
	}

	if previewOnly {
		return nil
	}

	if analysisReport.HasErrors() {
		return fmt.Errorf("migration cannot proceed due to errors. See report above for details")
	}

	m := migrator.New()
	result, err := m.Migrate(content, fromVersion, toVersion)
	if err != nil {
		return fmt.Errorf("error during migration: %w", err)
	}

	if output != nil {
		_, err = output.Write(result)
		if err != nil {
			return fmt.Errorf("error writing output: %w", err)
		}
	}

	if verbosity >= 1 && outputFile != "-" {
		fmt.Fprintf(os.Stderr, "Migration completed successfully. Output written to: %s\n", outputFile)
	}

	return nil
}