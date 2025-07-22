package migrator

import (
	"strings"
	"testing"
)

func TestMigratorService_New(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("Expected migrator service to be created, got nil")
	}
}

func TestMigratorService_Migrate_QTI12to21(t *testing.T) {
	qti12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
		<presentation>
			<material>
				<mattext texttype="text/plain">What is 2 + 2?</mattext>
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
			<outcomes>
				<decvar varname="SCORE" vartype="decimal" defaultval="0"/>
			</outcomes>
			<respcondition continue="no">
				<conditionvar>
					<varequal respident="RESPONSE">B</varequal>
				</conditionvar>
				<setvar action="set" varname="SCORE">1</setvar>
			</respcondition>
		</resprocessing>
	</item>
</questestinterop>`

	m := New()
	result, err := m.Migrate([]byte(qti12XML), "1.2", "2.1")
	
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}
	
	if len(result) == 0 {
		t.Fatal("Expected migration result to have content")
	}
	
	resultStr := string(result)
	
	// Check that XML header is present
	if !strings.Contains(resultStr, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("Expected XML header in result")
	}
	
	// Check version was updated
	if !strings.Contains(resultStr, `version="2.1"`) {
		t.Error("Expected version to be updated to 2.1")
	}
	
	// Check that itemBody element is present (QTI 2.1 structure)
	if !strings.Contains(resultStr, `<itemBody>`) {
		t.Error("Expected itemBody element in QTI 2.1 output")
	}
	
	// Check that responseDeclaration is present
	if !strings.Contains(resultStr, `<responseDeclaration`) {
		t.Error("Expected responseDeclaration in QTI 2.1 output")
	}
	
	// Check that outcomeDeclaration is present
	if !strings.Contains(resultStr, `<outcomeDeclaration`) {
		t.Error("Expected outcomeDeclaration in QTI 2.1 output")
	}
}

func TestMigratorService_Migrate_UnsupportedSourceVersion(t *testing.T) {
	qti12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
	</item>
</questestinterop>`

	m := New()
	_, err := m.Migrate([]byte(qti12XML), "invalid", "2.1")
	
	if err == nil {
		t.Error("Expected error for invalid source version, but got none")
	}
	
	if !strings.Contains(err.Error(), "failed to get parser") {
		t.Errorf("Expected parser error, got: %v", err)
	}
}

func TestMigratorService_Migrate_UnsupportedMigrationPath(t *testing.T) {
	qti12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
	</item>
</questestinterop>`

	m := New()
	_, err := m.Migrate([]byte(qti12XML), "1.2", "3.0")
	
	if err == nil {
		t.Error("Expected error for unsupported migration path, but got none")
	}
	
	if !strings.Contains(err.Error(), "unsupported migration path") {
		t.Errorf("Expected unsupported migration path error, got: %v", err)
	}
}

func TestMigratorService_Migrate_QTI21to30_NotImplemented(t *testing.T) {
	qti21XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="2.1">
	<item ident="q001" title="Test Question">
		<itemBody>
			<p>What is 2 + 2?</p>
		</itemBody>
	</item>
</questestinterop>`

	m := New()
	_, err := m.Migrate([]byte(qti21XML), "2.1", "3.0")
	
	if err == nil {
		t.Error("Expected error for unimplemented migration, but got none")
	}
	
	if !strings.Contains(err.Error(), "QTI 2.1 to 3.0 migration not yet implemented") {
		t.Errorf("Expected not implemented error, got: %v", err)
	}
}

func TestMigratorService_Migrate_InvalidXML(t *testing.T) {
	invalidXML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
</questestinterop`  // Missing closing tags

	m := New()
	_, err := m.Migrate([]byte(invalidXML), "1.2", "2.1")
	
	if err == nil {
		t.Error("Expected error for invalid XML, but got none")
	}
	
	if !strings.Contains(err.Error(), "failed to parse source document") {
		t.Errorf("Expected parsing error, got: %v", err)
	}
}

func TestMigratorService_Migrate_ComplexDocument(t *testing.T) {
	complexQTI12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<assessment ident="test001" title="Sample Test">
		<section ident="sec001" title="Section 1">
			<item ident="q001" title="Multiple Choice Question">
				<metadata>
					<qtimetadata>
						<interactiontype>multiple_choice</interactiontype>
					</qtimetadata>
				</metadata>
				<presentation>
					<material>
						<mattext texttype="text/html"><![CDATA[<p>What is the capital of France?</p>]]></mattext>
						<matimage imagetype="image/jpeg" uri="france_map.jpg" width="300" height="200"/>
					</material>
					<response_lid ident="RESPONSE" rcardinality="single">
						<render_choice shuffle="no" maxnumber="1">
							<response_label ident="A">
								<material><mattext>London</mattext></material>
							</response_label>
							<response_label ident="B">
								<material><mattext>Berlin</mattext></material>
							</response_label>
							<response_label ident="C">
								<material><mattext>Paris</mattext></material>
							</response_label>
						</render_choice>
					</response_lid>
				</presentation>
				<resprocessing scoremodel="SumOfScores">
					<outcomes>
						<decvar varname="SCORE" vartype="decimal" defaultval="0" maxvalue="1"/>
					</outcomes>
					<respcondition continue="no">
						<conditionvar>
							<varequal respident="RESPONSE">C</varequal>
						</conditionvar>
						<setvar action="set" varname="SCORE">1</setvar>
					</respcondition>
				</resprocessing>
				<itemfeedback ident="correct_fb" title="Correct">
					<material>
						<mattext>Correct! Paris is the capital of France.</mattext>
					</material>
				</itemfeedback>
			</item>
		</section>
	</assessment>
</questestinterop>`

	m := New()
	result, err := m.Migrate([]byte(complexQTI12XML), "1.2", "2.1")
	
	if err != nil {
		t.Fatalf("Migration of complex document failed: %v", err)
	}
	
	resultStr := string(result)
	
	// Check assessment structure is preserved
	if !strings.Contains(resultStr, `<assessment`) {
		t.Error("Expected assessment element to be preserved")
	}
	
	// Check section structure is preserved
	if !strings.Contains(resultStr, `<section`) {
		t.Error("Expected section element to be preserved")
	}
	
	// Check metadata is preserved
	if !strings.Contains(resultStr, `<metadata>`) {
		t.Error("Expected metadata to be preserved")
	}
	
	// Check QTI 2.1 specific elements
	if !strings.Contains(resultStr, `<choiceInteraction`) {
		t.Error("Expected choiceInteraction in QTI 2.1 output")
	}
	
	if !strings.Contains(resultStr, `<simpleChoice`) {
		t.Error("Expected simpleChoice elements in QTI 2.1 output")
	}
	
	// Check that shuffle attribute is handled (omitempty when false)
	if strings.Contains(resultStr, `shuffle=`) {
		t.Error("Expected shuffle attribute to be omitted when false (omitempty)")
	}
	
	// Check feedback is preserved
	if !strings.Contains(resultStr, `<itemfeedback`) {
		t.Error("Expected itemfeedback to be preserved")
	}
}

func TestMigratorService_Migrate_EmptyDocument(t *testing.T) {
	emptyQTI12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
</questestinterop>`

	m := New()
	result, err := m.Migrate([]byte(emptyQTI12XML), "1.2", "2.1")
	
	if err != nil {
		t.Fatalf("Migration of empty document failed: %v", err)
	}
	
	resultStr := string(result)
	
	// Should still have basic structure
	if !strings.Contains(resultStr, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("Expected XML header")
	}
	
	if !strings.Contains(resultStr, `<questestinterop`) {
		t.Error("Expected questestinterop root element")
	}
	
	if !strings.Contains(resultStr, `version="2.1"`) {
		t.Error("Expected version to be updated to 2.1")
	}
}

func BenchmarkMigratorService_Migrate_SimpleDocument(b *testing.B) {
	qti12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
		<presentation>
			<material>
				<mattext texttype="text/plain">What is 2 + 2?</mattext>
			</material>
		</presentation>
	</item>
</questestinterop>`

	m := New()
	data := []byte(qti12XML)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.Migrate(data, "1.2", "2.1")
		if err != nil {
			b.Fatalf("Migration failed: %v", err)
		}
	}
}

func BenchmarkMigratorService_Migrate_ComplexDocument(b *testing.B) {
	complexQTI12XML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<assessment ident="test001" title="Sample Test">
		<section ident="sec001" title="Section 1">
			<item ident="q001" title="Question 1">
				<presentation>
					<material>
						<mattext texttype="text/html">Question content</mattext>
					</material>
					<response_lid ident="RESPONSE" rcardinality="single">
						<render_choice shuffle="yes">
							<response_label ident="A"><material><mattext>Option A</mattext></material></response_label>
							<response_label ident="B"><material><mattext>Option B</mattext></material></response_label>
						</render_choice>
					</response_lid>
				</presentation>
				<resprocessing>
					<outcomes><decvar varname="SCORE" vartype="decimal" defaultval="0"/></outcomes>
					<respcondition continue="no">
						<conditionvar><varequal respident="RESPONSE">A</varequal></conditionvar>
						<setvar action="set" varname="SCORE">1</setvar>
					</respcondition>
				</resprocessing>
			</item>
		</section>
	</assessment>
</questestinterop>`

	m := New()
	data := []byte(complexQTI12XML)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.Migrate(data, "1.2", "2.1")
		if err != nil {
			b.Fatalf("Migration failed: %v", err)
		}
	}
}