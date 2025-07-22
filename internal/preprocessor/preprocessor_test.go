package preprocessor

import (
	"testing"
)

func TestPreprocessor_New(t *testing.T) {
	verbosity := 2
	p := New(verbosity)
	
	if p == nil {
		t.Fatal("Expected preprocessor to be created, got nil")
	}
	
	if p.verbosity != verbosity {
		t.Errorf("Expected verbosity %d, got %d", verbosity, p.verbosity)
	}
}

func TestPreprocessor_Analyze_QTI12to21(t *testing.T) {
	qti12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
		<metadata>
			<qtimetadata>
				<interactiontype>multiple_choice</interactiontype>
			</qtimetadata>
		</metadata>
		<presentation>
			<material>
				<mattext texttype="text/html">What is 2 + 2?</mattext>
			</material>
			<response_lid ident="RESPONSE" rcardinality="single">
				<render_choice shuffle="yes">
					<response_label ident="A">
						<material><mattext>3</mattext></material>
					</response_label>
					<response_label ident="B">
						<material><mattext>4</mattext></material>
					</response_label>
				</render_choice>
			</response_lid>
		</presentation>
		<resprocessing>
			<respcondition continue="yes">
				<conditionvar>
					<varequal respident="RESPONSE">B</varequal>
				</conditionvar>
				<setvar action="set" varname="SCORE">1</setvar>
			</respcondition>
		</resprocessing>
	</item>
</questestinterop>`

	p := New(2)
	report, err := p.Analyze([]byte(qti12XML), "1.2", "2.1")
	
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}
	
	if report.SourceVersion != "1.2" {
		t.Errorf("Expected source version '1.2', got '%s'", report.SourceVersion)
	}
	
	if report.TargetVersion != "2.1" {
		t.Errorf("Expected target version '2.1', got '%s'", report.TargetVersion)
	}
	
	if report.TotalItems != 1 {
		t.Errorf("Expected 1 total item, got %d", report.TotalItems)
	}
	
	if report.CompatibleItems != 1 {
		t.Errorf("Expected 1 compatible item, got %d", report.CompatibleItems)
	}
	
	if report.HasErrors() {
		t.Errorf("Expected no fatal errors, got %d errors", len(report.Errors))
	}
	
	// Check that migration details are generated
	found := false
	for _, detail := range report.MigrationDetails {
		if detail.Action == "transform" && detail.OldValue == `shuffle="yes"` {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected migration detail for shuffle attribute transformation")
	}
}

func TestPreprocessor_Analyze_QTI21to30(t *testing.T) {
	qti21XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="2.1">
	<item ident="q001" title="Test Question">
		<itemBody>
			<p>What is 2 + 2?</p>
		</itemBody>
	</item>
</questestinterop>`

	p := New(1)
	report, err := p.Analyze([]byte(qti21XML), "2.1", "3.0")
	
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}
	
	if !report.HasErrors() {
		t.Error("Expected fatal error for unsupported migration, but got none")
	}
	
	if report.IncompatibleItems != 1 {
		t.Errorf("Expected 1 incompatible item, got %d", report.IncompatibleItems)
	}
	
	// Check that error message is appropriate
	found := false
	for _, err := range report.Errors {
		if err.Fatal && err.Message == "QTI 2.1 to 3.0 migration is not yet implemented" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected fatal error message for unsupported migration")
	}
}

func TestPreprocessor_Analyze_UnsupportedMigrationPath(t *testing.T) {
	qti12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
	</item>
</questestinterop>`

	p := New(1)
	_, err := p.Analyze([]byte(qti12XML), "1.2", "3.0")
	
	if err == nil {
		t.Error("Expected error for unsupported migration path, but got none")
	}
	
	if err.Error() != "unsupported migration path: 1.2 to 3.0" {
		t.Errorf("Expected unsupported migration path error, got: %v", err)
	}
}

func TestPreprocessor_Analyze_InvalidSourceVersion(t *testing.T) {
	qti12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
	</item>
</questestinterop>`

	p := New(1)
	_, err := p.Analyze([]byte(qti12XML), "invalid", "2.1")
	
	if err == nil {
		t.Error("Expected error for invalid source version, but got none")
	}
}

func TestPreprocessor_Analyze_InvalidXML(t *testing.T) {
	invalidXML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
</questestinterop`  // Missing closing tags

	p := New(1)
	_, err := p.Analyze([]byte(invalidXML), "1.2", "2.1")
	
	if err == nil {
		t.Error("Expected error for invalid XML, but got none")
	}
}

func TestPreprocessor_AnalyzeItem12to21_ComplexScenario(t *testing.T) {
	qti12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
		<metadata>
			<qtimetadata>
				<interactiontype>invalid_type</interactiontype>
			</qtimetadata>
		</metadata>
		<presentation>
			<material>
				<mattext texttype="text/html"><![CDATA[<p>What is 2 + 2?</p>]]></mattext>
				<matimage uri="test.jpg"/>
			</material>
			<response_lid ident="RESPONSE" rcardinality="multiple">
				<render_choice shuffle="no" maxnumber="2">
					<response_label ident="A">
						<material><mattext>3</mattext></material>
					</response_label>
				</render_choice>
			</response_lid>
		</presentation>
		<resprocessing>
			<respcondition continue="no">
				<conditionvar>
					<varequal respident="RESPONSE">A</varequal>
				</conditionvar>
				<setvar action="add" varname="SCORE">1</setvar>
			</respcondition>
		</resprocessing>
	</item>
</questestinterop>`

	p := New(3) // High verbosity to get all details
	report, err := p.Analyze([]byte(qti12XML), "1.2", "2.1")
	
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}
	
	// Should have warnings for invalid interaction type
	warningFound := false
	for _, warning := range report.Warnings {
		if warning.ItemID == "q001" && 
		   warning.Message == "Interaction type 'invalid_type' may need adjustment for QTI 2.1" {
			warningFound = true
			break
		}
	}
	if !warningFound {
		t.Error("Expected warning for invalid interaction type")
	}
	
	// Should have warning for missing image type
	imageWarningFound := false
	for _, warning := range report.Warnings {
		if warning.Message == "Image type not specified" {
			imageWarningFound = true
			break
		}
	}
	if !imageWarningFound {
		t.Error("Expected warning for missing image type")
	}
	
	// Should have migration detail for HTML content
	htmlDetailFound := false
	for _, detail := range report.MigrationDetails {
		if detail.Action == "validate" && 
		   detail.OldValue == "text/html content" {
			htmlDetailFound = true
			break
		}
	}
	if !htmlDetailFound {
		t.Error("Expected migration detail for HTML content validation")
	}
}

func TestPreprocessor_AnalyzeMaterial12to21(t *testing.T) {
	// Test with different verbosity levels
	testCases := []struct {
		name      string
		verbosity int
		expectDetails bool
	}{
		{"Low verbosity", 1, false},
		{"High verbosity", 2, true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			qti12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
		<presentation>
			<material>
				<mattext texttype="text/html">HTML content</mattext>
				<matimage uri="test.jpg"/>
			</material>
		</presentation>
	</item>
</questestinterop>`

			p := New(tc.verbosity)
			report, err := p.Analyze([]byte(qti12XML), "1.2", "2.1")
			
			if err != nil {
				t.Fatalf("Analysis failed: %v", err)
			}
			
			htmlDetailFound := false
			for _, detail := range report.MigrationDetails {
				if detail.Action == "validate" {
					htmlDetailFound = true
					break
				}
			}
			
			if tc.expectDetails && !htmlDetailFound {
				t.Error("Expected HTML validation details with high verbosity")
			} else if !tc.expectDetails && htmlDetailFound {
				t.Error("Did not expect HTML validation details with low verbosity")
			}
		})
	}
}

func TestAnalysisReport_HasErrors(t *testing.T) {
	report := &AnalysisReport{}
	
	// No errors
	if report.HasErrors() {
		t.Error("Expected HasErrors to return false for empty error list")
	}
	
	// Non-fatal error
	report.Errors = append(report.Errors, Error{
		Message: "Non-fatal error",
		Fatal:   false,
	})
	
	if report.HasErrors() {
		t.Error("Expected HasErrors to return false for non-fatal errors")
	}
	
	// Fatal error
	report.Errors = append(report.Errors, Error{
		Message: "Fatal error",
		Fatal:   true,
	})
	
	if !report.HasErrors() {
		t.Error("Expected HasErrors to return true when fatal errors are present")
	}
}

func TestIsValidQTI21InteractionType(t *testing.T) {
	validTypes := []string{
		"choiceInteraction",
		"orderInteraction",
		"associateInteraction",
		"matchInteraction",
		"textEntryInteraction",
		"extendedTextInteraction",
		"hotspotInteraction",
		"sliderInteraction",
		"uploadInteraction",
		"customInteraction",
	}
	
	invalidTypes := []string{
		"multiple_choice",
		"fill_in_blank",
		"invalid_type",
		"",
		"CHOICEINTERACTION",  // Case sensitive
	}
	
	for _, validType := range validTypes {
		if !isValidQTI21InteractionType(validType) {
			t.Errorf("Expected '%s' to be a valid QTI 2.1 interaction type", validType)
		}
	}
	
	for _, invalidType := range invalidTypes {
		if isValidQTI21InteractionType(invalidType) {
			t.Errorf("Expected '%s' to be an invalid QTI 2.1 interaction type", invalidType)
		}
	}
}

func BenchmarkPreprocessor_Analyze(b *testing.B) {
	qti12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
		<presentation>
			<material>
				<mattext texttype="text/plain">What is 2 + 2?</mattext>
			</material>
			<response_lid ident="RESPONSE" rcardinality="single">
				<render_choice shuffle="no">
					<response_label ident="A">
						<material><mattext>3</mattext></material>
					</response_label>
					<response_label ident="B">
						<material><mattext>4</mattext></material>
					</response_label>
				</render_choice>
			</response_lid>
		</presentation>
	</item>
</questestinterop>`

	p := New(1)
	data := []byte(qti12XML)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.Analyze(data, "1.2", "2.1")
		if err != nil {
			b.Fatalf("Analysis failed: %v", err)
		}
	}
}