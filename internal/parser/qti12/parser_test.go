package qti12

import (
	"strings"
	"testing"

	"github.com/qti-migrator/pkg/models"
)

func TestParser12_Version(t *testing.T) {
	parser := New()
	if parser.Version() != "1.2" {
		t.Errorf("Expected version '1.2', got '%s'", parser.Version())
	}
}

func TestParser12_Parse_ValidDocument(t *testing.T) {
	validXML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
		<presentation>
			<material>
				<mattext texttype="text/plain">What is 2 + 2?</mattext>
			</material>
			<response_lid ident="RESPONSE" rcardinality="single">
				<render_choice shuffle="no">
					<response_label ident="A">
						<material>
							<mattext texttype="text/plain">3</mattext>
						</material>
					</response_label>
					<response_label ident="B">
						<material>
							<mattext texttype="text/plain">4</mattext>
						</material>
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

	parser := New()
	doc, err := parser.Parse([]byte(validXML))
	
	if err != nil {
		t.Fatalf("Failed to parse valid QTI 1.2 document: %v", err)
	}
	
	if doc.Version != "1.2" {
		t.Errorf("Expected version '1.2', got '%s'", doc.Version)
	}
	
	if len(doc.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(doc.Items))
	}
	
	item := doc.Items[0]
	if item.Ident != "q001" {
		t.Errorf("Expected item ident 'q001', got '%s'", item.Ident)
	}
	
	if item.Title != "Test Question" {
		t.Errorf("Expected item title 'Test Question', got '%s'", item.Title)
	}
}

func TestParser12_Parse_InvalidXML(t *testing.T) {
	invalidXML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
		<presentation>
			<material>
				<mattext texttype="text/plain">What is 2 + 2?</mattext>
			</material>
		</presentation>
	</item>
</questestinterop`  // Missing closing tag

	parser := New()
	_, err := parser.Parse([]byte(invalidXML))
	
	if err == nil {
		t.Error("Expected error for invalid XML, but got none")
	}
	
	if !strings.Contains(err.Error(), "failed to parse QTI 1.2 document") {
		t.Errorf("Expected parsing error message, got: %v", err)
	}
}

func TestParser12_Parse_WrongVersion(t *testing.T) {
	wrongVersionXML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="2.1">
	<item ident="q001" title="Test Question">
	</item>
</questestinterop>`

	parser := New()
	_, err := parser.Parse([]byte(wrongVersionXML))
	
	if err == nil {
		t.Error("Expected error for wrong version, but got none")
	}
	
	if !strings.Contains(err.Error(), "invalid QTI version") {
		t.Errorf("Expected version error message, got: %v", err)
	}
}

func TestParser12_Parse_NoVersion(t *testing.T) {
	noVersionXML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop>
	<item ident="q001" title="Test Question">
	</item>
</questestinterop>`

	parser := New()
	doc, err := parser.Parse([]byte(noVersionXML))
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if doc.Version != "1.2" {
		t.Errorf("Expected default version '1.2', got '%s'", doc.Version)
	}
}

func TestIsValidQTI12Version(t *testing.T) {
	validVersions := []string{"1.2", "1.2.0", "1.2.1"}
	invalidVersions := []string{"1.1", "2.0", "2.1", "3.0", ""}
	
	for _, version := range validVersions {
		if !isValidQTI12Version(version) {
			t.Errorf("Expected version '%s' to be valid for QTI 1.2", version)
		}
	}
	
	for _, version := range invalidVersions {
		if isValidQTI12Version(version) {
			t.Errorf("Expected version '%s' to be invalid for QTI 1.2", version)
		}
	}
}

func TestParser12_Parse_ComplexDocument(t *testing.T) {
	complexXML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<assessment ident="test001" title="Sample Test">
		<section ident="sec001" title="Section 1">
			<item ident="q001" title="Multiple Choice">
				<metadata>
					<qtimetadata>
						<interactiontype>multiple_choice</interactiontype>
					</qtimetadata>
				</metadata>
				<presentation>
					<material>
						<mattext texttype="text/html">What is the capital of France?</mattext>
						<matimage imagetype="image/jpeg" uri="france_map.jpg" width="300" height="200"/>
					</material>
					<response_lid ident="RESPONSE" rcardinality="single">
						<render_choice shuffle="yes" maxnumber="1">
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
				<resprocessing>
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

	parser := New()
	doc, err := parser.Parse([]byte(complexXML))
	
	if err != nil {
		t.Fatalf("Failed to parse complex QTI 1.2 document: %v", err)
	}
	
	// Check assessment structure
	if doc.Assessment == nil {
		t.Fatal("Expected assessment to be present")
	}
	
	if doc.Assessment.Ident != "test001" {
		t.Errorf("Expected assessment ident 'test001', got '%s'", doc.Assessment.Ident)
	}
	
	if len(doc.Assessment.Sections) != 1 {
		t.Errorf("Expected 1 section, got %d", len(doc.Assessment.Sections))
	}
	
	section := doc.Assessment.Sections[0]
	if len(section.Items) != 1 {
		t.Errorf("Expected 1 item in section, got %d", len(section.Items))
	}
	
	item := section.Items[0]
	if item.Metadata == nil || item.Metadata.QTIMetadata == nil {
		t.Fatal("Expected QTI metadata to be present")
	}
	
	if item.Metadata.QTIMetadata.InteractionType != "multiple_choice" {
		t.Errorf("Expected interaction type 'multiple_choice', got '%s'", 
			item.Metadata.QTIMetadata.InteractionType)
	}
	
	// Check presentation structure
	if item.Presentation == nil {
		t.Fatal("Expected presentation to be present")
	}
	
	if len(item.Presentation.Response) != 1 {
		t.Errorf("Expected 1 response, got %d", len(item.Presentation.Response))
	}
	
	response := item.Presentation.Response[0]
	if response.RenderChoice == nil {
		t.Fatal("Expected render_choice to be present")
	}
	
	if len(response.RenderChoice.ResponseLabel) != 3 {
		t.Errorf("Expected 3 response labels, got %d", len(response.RenderChoice.ResponseLabel))
	}
	
	// Check feedback
	if len(item.Feedback) != 1 {
		t.Errorf("Expected 1 feedback item, got %d", len(item.Feedback))
	}
}

func BenchmarkParser12_Parse(b *testing.B) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
	<item ident="q001" title="Test Question">
		<presentation>
			<material>
				<mattext texttype="text/plain">What is 2 + 2?</mattext>
			</material>
		</presentation>
	</item>
</questestinterop>`
	
	parser := New()
	data := []byte(xml)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(data)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}