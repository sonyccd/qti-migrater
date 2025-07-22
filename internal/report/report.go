package report

import (
	"fmt"
	"strings"

	"github.com/qti-migrator/internal/preprocessor"
)

type Reporter struct {
	verbosity int
}

func New(verbosity int) *Reporter {
	return &Reporter{
		verbosity: verbosity,
	}
}

func (r *Reporter) Generate(report *preprocessor.AnalysisReport) string {
	var builder strings.Builder

	builder.WriteString(r.generateHeader(report))
	builder.WriteString(r.generateSummary(report))

	if len(report.Errors) > 0 {
		builder.WriteString(r.generateErrors(report))
	}

	if len(report.Warnings) > 0 && r.verbosity >= 1 {
		builder.WriteString(r.generateWarnings(report))
	}

	if len(report.MigrationDetails) > 0 && r.verbosity >= 2 {
		builder.WriteString(r.generateMigrationDetails(report))
	}

	builder.WriteString(r.generateFooter(report))

	return builder.String()
}

func (r *Reporter) generateHeader(report *preprocessor.AnalysisReport) string {
	return fmt.Sprintf(`
================================================================================
                          QTI Migration Analysis Report
================================================================================
Migration Path: QTI %s → QTI %s
================================================================================

`, report.SourceVersion, report.TargetVersion)
}

func (r *Reporter) generateSummary(report *preprocessor.AnalysisReport) string {
	status := "READY"
	if report.HasErrors() {
		status = "BLOCKED"
	}

	summary := fmt.Sprintf(`SUMMARY
-------
Status: %s
Total Items: %d
Compatible Items: %d
Items Requiring Attention: %d
Errors: %d
Warnings: %d

`, status, report.TotalItems, report.CompatibleItems, report.IncompatibleItems,
		len(report.Errors), len(report.Warnings))

	return summary
}

func (r *Reporter) generateErrors(report *preprocessor.AnalysisReport) string {
	var builder strings.Builder

	builder.WriteString("ERRORS (Migration Blockers)\n")
	builder.WriteString("--------------------------\n")

	for i, err := range report.Errors {
		builder.WriteString(fmt.Sprintf("%d. ", i+1))
		if err.ItemID != "" {
			builder.WriteString(fmt.Sprintf("[Item: %s] ", err.ItemID))
		}
		builder.WriteString(fmt.Sprintf("%s\n", err.Message))
		
		if r.verbosity >= 2 && err.ElementPath != "" {
			builder.WriteString(fmt.Sprintf("   Path: %s\n", err.ElementPath))
		}
		
		if err.Fatal {
			builder.WriteString("   ⚠️  This error must be resolved before migration can proceed.\n")
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func (r *Reporter) generateWarnings(report *preprocessor.AnalysisReport) string {
	var builder strings.Builder

	builder.WriteString("WARNINGS\n")
	builder.WriteString("--------\n")

	for i, warning := range report.Warnings {
		builder.WriteString(fmt.Sprintf("%d. ", i+1))
		if warning.ItemID != "" {
			builder.WriteString(fmt.Sprintf("[Item: %s] ", warning.ItemID))
		}
		builder.WriteString(fmt.Sprintf("%s\n", warning.Message))
		
		if r.verbosity >= 2 && warning.ElementPath != "" {
			builder.WriteString(fmt.Sprintf("   Path: %s\n", warning.ElementPath))
		}
		
		if warning.Suggestion != "" {
			builder.WriteString(fmt.Sprintf("   → %s\n", warning.Suggestion))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func (r *Reporter) generateMigrationDetails(report *preprocessor.AnalysisReport) string {
	var builder strings.Builder

	builder.WriteString("MIGRATION DETAILS\n")
	builder.WriteString("-----------------\n")

	groupedDetails := r.groupDetailsByAction(report.MigrationDetails)

	for action, details := range groupedDetails {
		builder.WriteString(fmt.Sprintf("\n%s Actions (%d):\n", strings.Title(action), len(details)))
		builder.WriteString(strings.Repeat("-", len(action)+15) + "\n")

		for i, detail := range details {
			builder.WriteString(fmt.Sprintf("%d. ", i+1))
			if detail.ItemID != "" {
				builder.WriteString(fmt.Sprintf("[Item: %s] ", detail.ItemID))
			}
			builder.WriteString(fmt.Sprintf("%s\n", detail.Description))
			
			if r.verbosity >= 3 {
				builder.WriteString(fmt.Sprintf("   Path: %s\n", detail.ElementPath))
				if detail.OldValue != "" {
					builder.WriteString(fmt.Sprintf("   Old: %s\n", r.truncateValue(detail.OldValue)))
				}
				if detail.NewValue != "" {
					builder.WriteString(fmt.Sprintf("   New: %s\n", r.truncateValue(detail.NewValue)))
				}
			}
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func (r *Reporter) groupDetailsByAction(details []preprocessor.MigrationDetail) map[string][]preprocessor.MigrationDetail {
	grouped := make(map[string][]preprocessor.MigrationDetail)
	
	for _, detail := range details {
		grouped[detail.Action] = append(grouped[detail.Action], detail)
	}
	
	return grouped
}

func (r *Reporter) truncateValue(value string) string {
	if len(value) > 50 {
		return value[:47] + "..."
	}
	return value
}

func (r *Reporter) generateFooter(report *preprocessor.AnalysisReport) string {
	var builder strings.Builder

	builder.WriteString("\n")
	builder.WriteString("================================================================================\n")
	
	if report.HasErrors() {
		builder.WriteString("⚠️  MIGRATION BLOCKED: Please resolve the errors listed above before proceeding.\n")
	} else if len(report.Warnings) > 0 {
		builder.WriteString("✓ Migration can proceed. Please review warnings for potential issues.\n")
	} else {
		builder.WriteString("✓ Migration can proceed without issues.\n")
	}
	
	if r.verbosity < 3 && (len(report.Warnings) > 0 || len(report.MigrationDetails) > 0) {
		builder.WriteString("\nTip: Use -v 2 or -v 3 for more detailed information.\n")
	}
	
	builder.WriteString("================================================================================\n")

	return builder.String()
}