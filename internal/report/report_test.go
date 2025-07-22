package report

import (
	"strings"
	"testing"

	"github.com/qti-migrator/internal/preprocessor"
)

func TestReporter_New(t *testing.T) {
	verbosity := 2
	r := New(verbosity)
	
	if r == nil {
		t.Fatal("Expected reporter to be created, got nil")
	}
	
	if r.verbosity != verbosity {
		t.Errorf("Expected verbosity %d, got %d", verbosity, r.verbosity)
	}
}

func TestReporter_Generate_BasicReport(t *testing.T) {
	report := &preprocessor.AnalysisReport{
		SourceVersion:     "1.2",
		TargetVersion:     "2.1",
		TotalItems:        5,
		CompatibleItems:   4,
		IncompatibleItems: 1,
		Warnings: []preprocessor.Warning{
			{
				ItemID:     "q001",
				Message:    "Test warning",
				Suggestion: "Test suggestion",
			},
		},
		Errors: []preprocessor.Error{
			{
				ItemID:  "q002",
				Message: "Test error",
				Fatal:   true,
			},
		},
	}
	
	r := New(1)
	output := r.Generate(report)
	
	// Check header
	if !strings.Contains(output, "QTI Migration Analysis Report") {
		t.Error("Expected report header")
	}
	
	if !strings.Contains(output, "QTI 1.2 → QTI 2.1") {
		t.Error("Expected migration path in header")
	}
	
	// Check summary
	if !strings.Contains(output, "Status: BLOCKED") {
		t.Error("Expected BLOCKED status due to fatal error")
	}
	
	if !strings.Contains(output, "Total Items: 5") {
		t.Error("Expected total items count")
	}
	
	if !strings.Contains(output, "Compatible Items: 4") {
		t.Error("Expected compatible items count")
	}
	
	if !strings.Contains(output, "Items Requiring Attention: 1") {
		t.Error("Expected incompatible items count")
	}
	
	if !strings.Contains(output, "Errors: 1") {
		t.Error("Expected error count")
	}
	
	if !strings.Contains(output, "Warnings: 1") {
		t.Error("Expected warning count")
	}
	
	// Check errors section
	if !strings.Contains(output, "ERRORS (Migration Blockers)") {
		t.Error("Expected errors section header")
	}
	
	if !strings.Contains(output, "[Item: q002] Test error") {
		t.Error("Expected error message")
	}
	
	// Check warnings section
	if !strings.Contains(output, "WARNINGS") {
		t.Error("Expected warnings section header")
	}
	
	if !strings.Contains(output, "[Item: q001] Test warning") {
		t.Error("Expected warning message")
	}
	
	if !strings.Contains(output, "→ Test suggestion") {
		t.Error("Expected warning suggestion")
	}
	
	// Check footer
	if !strings.Contains(output, "MIGRATION BLOCKED") {
		t.Error("Expected blocked migration message in footer")
	}
}

func TestReporter_Generate_ReadyReport(t *testing.T) {
	report := &preprocessor.AnalysisReport{
		SourceVersion:     "1.2",
		TargetVersion:     "2.1",
		TotalItems:        3,
		CompatibleItems:   3,
		IncompatibleItems: 0,
		Warnings: []preprocessor.Warning{
			{
				ItemID:     "q001",
				Message:    "Minor warning",
				Suggestion: "Optional improvement",
			},
		},
		Errors: []preprocessor.Error{
			{
				ItemID:  "q002",
				Message: "Non-fatal error",
				Fatal:   false,
			},
		},
	}
	
	r := New(1)
	output := r.Generate(report)
	
	// Should show READY status since no fatal errors
	if !strings.Contains(output, "Status: READY") {
		t.Error("Expected READY status")
	}
	
	// Should show success message in footer
	if !strings.Contains(output, "Migration can proceed. Please review warnings") {
		t.Error("Expected success message with warnings")
	}
}

func TestReporter_Generate_NoIssuesReport(t *testing.T) {
	report := &preprocessor.AnalysisReport{
		SourceVersion:     "1.2",
		TargetVersion:     "2.1",
		TotalItems:        2,
		CompatibleItems:   2,
		IncompatibleItems: 0,
		Warnings:          []preprocessor.Warning{},
		Errors:            []preprocessor.Error{},
	}
	
	r := New(1)
	output := r.Generate(report)
	
	if !strings.Contains(output, "Status: READY") {
		t.Error("Expected READY status")
	}
	
	if !strings.Contains(output, "Migration can proceed without issues") {
		t.Error("Expected clean success message")
	}
}

func TestReporter_Generate_VerbosityLevels(t *testing.T) {
	report := &preprocessor.AnalysisReport{
		SourceVersion:     "1.2",
		TargetVersion:     "2.1",
		TotalItems:        1,
		CompatibleItems:   1,
		IncompatibleItems: 0,
		Warnings: []preprocessor.Warning{
			{
				ItemID:      "q001",
				ElementPath: "/item[@ident='q001']/presentation",
				Message:     "Warning message",
				Suggestion:  "Suggestion",
			},
		},
		MigrationDetails: []preprocessor.MigrationDetail{
			{
				ItemID:      "q001",
				ElementPath: "/item[@ident='q001']/response",
				OldValue:    "old_value",
				NewValue:    "new_value",
				Action:      "transform",
				Description: "Transform description",
			},
		},
	}
	
	// Test verbosity level 0 (minimal)
	r0 := New(0)
	output0 := r0.Generate(report)
	
	// Should not include warnings at verbosity 0
	if strings.Contains(output0, "WARNINGS") {
		t.Error("Expected no warnings section at verbosity 0")
	}
	
	// Test verbosity level 1 (normal)
	r1 := New(1)
	output1 := r1.Generate(report)
	
	// Should include warnings but not detailed paths
	if !strings.Contains(output1, "WARNINGS") {
		t.Error("Expected warnings section at verbosity 1")
	}
	
	if strings.Contains(output1, "Path: /item") {
		t.Error("Did not expect detailed paths at verbosity 1")
	}
	
	// Should not include migration details at verbosity 1
	if strings.Contains(output1, "MIGRATION DETAILS") {
		t.Error("Did not expect migration details at verbosity 1")
	}
	
	// Test verbosity level 2 (detailed)
	r2 := New(2)
	output2 := r2.Generate(report)
	
	// Should include warnings with paths
	if !strings.Contains(output2, "Path: /item[@ident='q001']/presentation") {
		t.Error("Expected detailed paths at verbosity 2")
	}
	
	// Should include migration details
	if !strings.Contains(output2, "MIGRATION DETAILS") {
		t.Error("Expected migration details at verbosity 2")
	}
	
	if !strings.Contains(output2, "Transform Actions (1):") {
		t.Error("Expected grouped migration details by action")
	}
	
	// Test verbosity level 3 (debug)
	r3 := New(3)
	output3 := r3.Generate(report)
	
	// Should include old/new values in migration details
	if !strings.Contains(output3, "Old: old_value") {
		t.Error("Expected old values at verbosity 3")
	}
	
	if !strings.Contains(output3, "New: new_value") {
		t.Error("Expected new values at verbosity 3")
	}
}

func TestReporter_Generate_VerbosityHint(t *testing.T) {
	report := &preprocessor.AnalysisReport{
		SourceVersion:     "1.2",
		TargetVersion:     "2.1",
		TotalItems:        1,
		CompatibleItems:   1,
		IncompatibleItems: 0,
		Warnings: []preprocessor.Warning{
			{Message: "Test warning"},
		},
		MigrationDetails: []preprocessor.MigrationDetail{
			{Description: "Test detail", Action: "test"},
		},
	}
	
	// At low verbosity, should show hint
	r1 := New(1)
	output1 := r1.Generate(report)
	
	if !strings.Contains(output1, "Use -v 2 or -v 3 for more detailed information") {
		t.Error("Expected verbosity hint at low verbosity level")
	}
	
	// At high verbosity, should not show hint
	r3 := New(3)
	output3 := r3.Generate(report)
	
	if strings.Contains(output3, "Use -v 2 or -v 3") {
		t.Error("Did not expect verbosity hint at high verbosity level")
	}
}

func TestReporter_GroupDetailsByAction(t *testing.T) {
	r := New(2)
	
	details := []preprocessor.MigrationDetail{
		{Action: "transform", Description: "Transform 1"},
		{Action: "validate", Description: "Validate 1"},
		{Action: "transform", Description: "Transform 2"},
		{Action: "convert", Description: "Convert 1"},
		{Action: "validate", Description: "Validate 2"},
	}
	
	grouped := r.groupDetailsByAction(details)
	
	if len(grouped) != 3 {
		t.Errorf("Expected 3 action groups, got %d", len(grouped))
	}
	
	if len(grouped["transform"]) != 2 {
		t.Errorf("Expected 2 transform actions, got %d", len(grouped["transform"]))
	}
	
	if len(grouped["validate"]) != 2 {
		t.Errorf("Expected 2 validate actions, got %d", len(grouped["validate"]))
	}
	
	if len(grouped["convert"]) != 1 {
		t.Errorf("Expected 1 convert action, got %d", len(grouped["convert"]))
	}
}

func TestReporter_TruncateValue(t *testing.T) {
	r := New(2)
	
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "short value",
			expected: "short value",
		},
		{
			input:    "this is a very long value that should be truncated because it exceeds the maximum length",
			expected: "this is a very long value that should be trunca...",
		},
		{
			input:    strings.Repeat("a", 50),
			expected: strings.Repeat("a", 50),
		},
		{
			input:    strings.Repeat("a", 51),
			expected: strings.Repeat("a", 47) + "...",
		},
	}
	
	for i, tc := range testCases {
		result := r.truncateValue(tc.input)
		if result != tc.expected {
			t.Errorf("Test case %d: expected '%s', got '%s'", i+1, tc.expected, result)
		}
	}
}

func TestReporter_Generate_EmptyElementPaths(t *testing.T) {
	report := &preprocessor.AnalysisReport{
		SourceVersion:     "1.2",
		TargetVersion:     "2.1",
		TotalItems:        1,
		CompatibleItems:   1,
		IncompatibleItems: 0,
		Warnings: []preprocessor.Warning{
			{
				ItemID:      "",  // Empty item ID
				ElementPath: "",  // Empty element path
				Message:     "Global warning",
				Suggestion:  "Global suggestion",
			},
		},
		Errors: []preprocessor.Error{
			{
				ItemID:      "",  // Empty item ID
				ElementPath: "",  // Empty element path
				Message:     "Global error",
				Fatal:       true,
			},
		},
	}
	
	r := New(2)
	output := r.Generate(report)
	
	// Should handle empty item IDs gracefully
	if !strings.Contains(output, "Global warning") {
		t.Error("Expected global warning message")
	}
	
	if !strings.Contains(output, "Global error") {
		t.Error("Expected global error message")
	}
	
	// Should not show empty brackets for empty item IDs
	if strings.Contains(output, "[Item: ]") {
		t.Error("Did not expect empty item ID brackets")
	}
}

func BenchmarkReporter_Generate_SmallReport(b *testing.B) {
	report := &preprocessor.AnalysisReport{
		SourceVersion:     "1.2",
		TargetVersion:     "2.1",
		TotalItems:        1,
		CompatibleItems:   1,
		IncompatibleItems: 0,
		Warnings: []preprocessor.Warning{
			{ItemID: "q001", Message: "Test warning"},
		},
	}
	
	r := New(1)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Generate(report)
	}
}

func BenchmarkReporter_Generate_LargeReport(b *testing.B) {
	// Create a large report with many warnings and details
	var warnings []preprocessor.Warning
	var details []preprocessor.MigrationDetail
	
	for i := 0; i < 100; i++ {
		warnings = append(warnings, preprocessor.Warning{
			ItemID:     "q" + string(rune(i)),
			Message:    "Warning message " + string(rune(i)),
			Suggestion: "Suggestion " + string(rune(i)),
		})
		
		details = append(details, preprocessor.MigrationDetail{
			ItemID:      "q" + string(rune(i)),
			Action:      "transform",
			Description: "Detail " + string(rune(i)),
		})
	}
	
	report := &preprocessor.AnalysisReport{
		SourceVersion:     "1.2",
		TargetVersion:     "2.1",
		TotalItems:        100,
		CompatibleItems:   100,
		IncompatibleItems: 0,
		Warnings:          warnings,
		MigrationDetails:  details,
	}
	
	r := New(3)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Generate(report)
	}
}