package qti21

import (
	"strings"
	"testing"
)

func TestParser21_Version(t *testing.T) {
	parser := New()
	if parser.Version() != "2.1" {
		t.Errorf("Expected version '2.1', got '%s'", parser.Version())
	}
}

func TestParser21_Parse_ValidDocument(t *testing.T) {
	validXML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="2.1">
	<item ident="q001" title="Test Question">
		<responseDeclaration identifier="RESPONSE" cardinality="single" baseType="identifier">
			<correctResponse>
				<value>B</value>
			</correctResponse>
		</responseDeclaration>
		<outcomeDeclaration identifier="SCORE" cardinality="single" baseType="float">
			<defaultValue>
				<value>0</value>
			</defaultValue>
		</outcomeDeclaration>
		<itemBody>
			<p>What is 2 + 2?</p>
			<choiceInteraction responseIdentifier="RESPONSE" shuffle="false" maxChoices="1">
				<simpleChoice identifier="A">3</simpleChoice>
				<simpleChoice identifier="B">4</simpleChoice>
			</choiceInteraction>
		</itemBody>
	</item>
</questestinterop>`

	parser := New()
	doc, err := parser.Parse([]byte(validXML))
	
	if err != nil {
		t.Fatalf("Failed to parse valid QTI 2.1 document: %v", err)
	}
	
	if doc.Version != "2.1" {
		t.Errorf("Expected version '2.1', got '%s'", doc.Version)
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

func TestParser21_Parse_InvalidXML(t *testing.T) {
	invalidXML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="2.1">
	<item ident="q001" title="Test Question">
		<itemBody>
			<p>What is 2 + 2?</p>
		</itemBody>
	</item>
</questestinterop`  // Missing closing tag

	parser := New()
	_, err := parser.Parse([]byte(invalidXML))
	
	if err == nil {
		t.Error("Expected error for invalid XML, but got none")
	}
	
	if !strings.Contains(err.Error(), "failed to parse QTI 2.1 document") {
		t.Errorf("Expected parsing error message, got: %v", err)
	}
}

func TestParser21_Parse_WrongVersion(t *testing.T) {
	wrongVersionXML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="1.2">
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

func TestParser21_Parse_NoVersion(t *testing.T) {
	noVersionXML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop>
	<item ident="q001" title="Test Question">
		<itemBody>
			<p>What is 2 + 2?</p>
		</itemBody>
	</item>
</questestinterop>`

	parser := New()
	doc, err := parser.Parse([]byte(noVersionXML))
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if doc.Version != "2.1" {
		t.Errorf("Expected default version '2.1', got '%s'", doc.Version)
	}
}

func TestIsValidQTI21Version(t *testing.T) {
	validVersions := []string{"2.1", "2.1.0", "2.1.1", "2.2", "2.2.0", "2.2.1", "2.2.2", "2.2.3", "2.2.4"}
	invalidVersions := []string{"1.2", "2.0", "3.0", "2.3", ""}
	
	for _, version := range validVersions {
		if !isValidQTI21Version(version) {
			t.Errorf("Expected version '%s' to be valid for QTI 2.1", version)
		}
	}
	
	for _, version := range invalidVersions {
		if isValidQTI21Version(version) {
			t.Errorf("Expected version '%s' to be invalid for QTI 2.1", version)
		}
	}
}

func TestParser21_Parse_ComplexDocument(t *testing.T) {
	complexXML := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="2.1">
	<item ident="q001" title="Complex Question">
		<responseDeclaration identifier="RESPONSE" cardinality="multiple" baseType="identifier">
			<correctResponse>
				<value>A</value>
				<value>C</value>
			</correctResponse>
			<mapping lowerBound="0" upperBound="2" defaultValue="0">
				<mapEntry mapKey="A" mappedValue="1"/>
				<mapEntry mapKey="B" mappedValue="0"/>
				<mapEntry mapKey="C" mappedValue="1"/>
			</mapping>
		</responseDeclaration>
		<outcomeDeclaration identifier="SCORE" cardinality="single" baseType="float">
			<defaultValue>
				<value>0.0</value>
			</defaultValue>
		</outcomeDeclaration>
		<itemBody>
			<div class="question">
				<p>Select all correct answers:</p>
			</div>
			<choiceInteraction responseIdentifier="RESPONSE" shuffle="true" maxChoices="3" minChoices="1">
				<prompt>Which of the following are prime numbers?</prompt>
				<simpleChoice identifier="A" fixed="false">2</simpleChoice>
				<simpleChoice identifier="B" fixed="false">4</simpleChoice>
				<simpleChoice identifier="C" fixed="false">3</simpleChoice>
			</choiceInteraction>
		</itemBody>
	</item>
</questestinterop>`

	parser := New()
	doc, err := parser.Parse([]byte(complexXML))
	
	if err != nil {
		t.Fatalf("Failed to parse complex QTI 2.1 document: %v", err)
	}
	
	item := doc.Items[0]
	
	// Check response declarations
	if len(item.ResponseDecl) != 1 {
		t.Errorf("Expected 1 response declaration, got %d", len(item.ResponseDecl))
	}
	
	responseDecl := item.ResponseDecl[0]
	if responseDecl.Identifier != "RESPONSE" {
		t.Errorf("Expected response identifier 'RESPONSE', got '%s'", responseDecl.Identifier)
	}
	
	if responseDecl.Cardinality != "multiple" {
		t.Errorf("Expected cardinality 'multiple', got '%s'", responseDecl.Cardinality)
	}
	
	if responseDecl.BaseType != "identifier" {
		t.Errorf("Expected baseType 'identifier', got '%s'", responseDecl.BaseType)
	}
	
	// Check correct response
	if responseDecl.CorrectResponse == nil {
		t.Fatal("Expected correct response to be present")
	}
	
	if len(responseDecl.CorrectResponse.Value) != 2 {
		t.Errorf("Expected 2 correct values, got %d", len(responseDecl.CorrectResponse.Value))
	}
	
	// Check mapping
	if responseDecl.Mapping == nil {
		t.Fatal("Expected mapping to be present")
	}
	
	if len(responseDecl.Mapping.MapEntry) != 3 {
		t.Errorf("Expected 3 map entries, got %d", len(responseDecl.Mapping.MapEntry))
	}
	
	// Check outcome declarations
	if len(item.OutcomeDecl) != 1 {
		t.Errorf("Expected 1 outcome declaration, got %d", len(item.OutcomeDecl))
	}
	
	// Check item body
	if item.ItemBody == nil {
		t.Fatal("Expected item body to be present")
	}
	
	if len(item.ItemBody.ChoiceInteraction) != 1 {
		t.Errorf("Expected 1 choice interaction, got %d", len(item.ItemBody.ChoiceInteraction))
	}
	
	choiceInteraction := item.ItemBody.ChoiceInteraction[0]
	if choiceInteraction.ResponseIdent != "RESPONSE" {
		t.Errorf("Expected response identifier 'RESPONSE', got '%s'", choiceInteraction.ResponseIdent)
	}
	
	if !choiceInteraction.Shuffle {
		t.Error("Expected shuffle to be true")
	}
	
	if len(choiceInteraction.SimpleChoice) != 3 {
		t.Errorf("Expected 3 simple choices, got %d", len(choiceInteraction.SimpleChoice))
	}
}

func BenchmarkParser21_Parse(b *testing.B) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<questestinterop version="2.1">
	<item ident="q001" title="Test Question">
		<itemBody>
			<p>What is 2 + 2?</p>
		</itemBody>
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