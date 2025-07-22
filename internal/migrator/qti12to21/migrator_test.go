package qti12to21

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/qti-migrator/pkg/models"
)

func TestMigrator12to21_New(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("Expected migrator to be created, got nil")
	}
}

func TestMigrator12to21_Migrate_SimpleItem(t *testing.T) {
	doc := &models.QTIDocument{
		Version: "1.2",
		Items: []models.Item{
			{
				Ident: "q001",
				Title: "Test Question",
				Presentation: &models.Presentation{
					Material: &models.Material{
						MatText: []models.MatText{
							{Content: "What is 2 + 2?", TextType: "text/plain"},
						},
					},
					Response: []models.Response{
						{
							Ident:        "RESPONSE",
							RCardinality: "single",
							RenderChoice: &models.RenderChoice{
								Shuffle: "no",
								ResponseLabel: []models.ResponseLabel{
									{
										Ident: "A",
										Material: &models.Material{
											MatText: []models.MatText{{Content: "3"}},
										},
									},
									{
										Ident: "B",
										Material: &models.Material{
											MatText: []models.MatText{{Content: "4"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	m := New()
	result, err := m.Migrate(doc)
	
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}
	
	if len(result) == 0 {
		t.Fatal("Expected migration result to have content")
	}
	
	resultStr := string(result)
	
	// Check XML structure
	if !strings.Contains(resultStr, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("Expected XML header")
	}
	
	if !strings.Contains(resultStr, `version="2.1"`) {
		t.Error("Expected version to be updated to 2.1")
	}
	
	if !strings.Contains(resultStr, `<itemBody>`) {
		t.Error("Expected itemBody element")
	}
	
	if !strings.Contains(resultStr, `<choiceInteraction`) {
		t.Error("Expected choiceInteraction element")
	}
	
	if !strings.Contains(resultStr, `<responseDeclaration`) {
		t.Error("Expected responseDeclaration element")
	}
}

func TestMigrator12to21_Migrate_InvalidDocumentType(t *testing.T) {
	m := New()
	_, err := m.Migrate("invalid document type")
	
	if err == nil {
		t.Error("Expected error for invalid document type, but got none")
	}
	
	if !strings.Contains(err.Error(), "invalid document type") {
		t.Errorf("Expected invalid document type error, got: %v", err)
	}
}

func TestMigrator12to21_ConvertPresentationToItemBody(t *testing.T) {
	m := New()
	presentation := &models.Presentation{
		Material: &models.Material{
			MatText: []models.MatText{
				{Content: "Question text", TextType: "text/plain"},
			},
			MatImage: []models.MatImage{
				{URI: "image.jpg", Width: 100, Height: 50},
			},
		},
		Response: []models.Response{
			{
				Ident:        "RESPONSE",
				RCardinality: "single",
				RenderChoice: &models.RenderChoice{
					Shuffle: "yes",
					ResponseLabel: []models.ResponseLabel{
						{Ident: "A", Material: &models.Material{MatText: []models.MatText{{Content: "Option A"}}}},
						{Ident: "B", Material: &models.Material{MatText: []models.MatText{{Content: "Option B"}}}},
					},
				},
			},
		},
	}
	
	itemBody := m.convertPresentationToItemBody(presentation)
	
	if itemBody == nil {
		t.Fatal("Expected itemBody to be created")
	}
	
	// Check paragraphs
	if len(itemBody.P) == 0 {
		t.Error("Expected at least one paragraph")
	}
	
	// Check choice interaction
	if len(itemBody.ChoiceInteraction) != 1 {
		t.Errorf("Expected 1 choice interaction, got %d", len(itemBody.ChoiceInteraction))
	}
	
	choiceInteraction := itemBody.ChoiceInteraction[0]
	if choiceInteraction.ResponseIdent != "RESPONSE" {
		t.Errorf("Expected response identifier 'RESPONSE', got '%s'", choiceInteraction.ResponseIdent)
	}
	
	if !choiceInteraction.Shuffle {
		t.Error("Expected shuffle to be true")
	}
	
	if len(choiceInteraction.SimpleChoice) != 2 {
		t.Errorf("Expected 2 simple choices, got %d", len(choiceInteraction.SimpleChoice))
	}
}

func TestMigrator12to21_ConvertResponseToTextEntry(t *testing.T) {
	m := New()
	response := &models.Response{
		Ident: "TEXT_RESPONSE",
		RenderFib: &models.RenderFib{
			MaxChars: 50,
			Rows:     1,
		},
	}
	
	textEntry := m.convertResponseToTextEntryInteraction(response)
	
	if textEntry.ResponseIdent != "TEXT_RESPONSE" {
		t.Errorf("Expected response identifier 'TEXT_RESPONSE', got '%s'", textEntry.ResponseIdent)
	}
	
	if textEntry.ExpectedLength != 50 {
		t.Errorf("Expected expected length 50, got %d", textEntry.ExpectedLength)
	}
}

func TestMigrator12to21_ConvertResponseToExtendedText(t *testing.T) {
	m := New()
	response := &models.Response{
		Ident: "ESSAY_RESPONSE",
		RenderFib: &models.RenderFib{
			MaxChars: 500,
			Rows:     5,
		},
	}
	
	extText := m.convertResponseToExtendedTextInteraction(response)
	
	if extText.ResponseIdent != "ESSAY_RESPONSE" {
		t.Errorf("Expected response identifier 'ESSAY_RESPONSE', got '%s'", extText.ResponseIdent)
	}
	
	if extText.ExpectedLines != 5 {
		t.Errorf("Expected expected lines 5, got %d", extText.ExpectedLines)
	}
	
	if extText.ExpectedLength != 500 {
		t.Errorf("Expected expected length 500, got %d", extText.ExpectedLength)
	}
}

func TestMigrator12to21_DetermineCardinality(t *testing.T) {
	m := New()
	
	testCases := []struct {
		name           string
		response       models.Response
		expectedCard   string
	}{
		{
			name: "Explicit single",
			response: models.Response{RCardinality: "single"},
			expectedCard: "single",
		},
		{
			name: "Explicit multiple",
			response: models.Response{RCardinality: "multiple"},
			expectedCard: "multiple",
		},
		{
			name: "Explicit ordered",
			response: models.Response{RCardinality: "ordered"},
			expectedCard: "ordered",
		},
		{
			name: "Multiple choice with max > 1",
			response: models.Response{
				RenderChoice: &models.RenderChoice{MaxNumber: 2},
			},
			expectedCard: "multiple",
		},
		{
			name: "Default single",
			response: models.Response{},
			expectedCard: "single",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cardinality := m.determineCardinality(&tc.response)
			if cardinality != tc.expectedCard {
				t.Errorf("Expected cardinality '%s', got '%s'", tc.expectedCard, cardinality)
			}
		})
	}
}

func TestMigrator12to21_DetermineBaseType(t *testing.T) {
	m := New()
	
	testCases := []struct {
		name          string
		response      models.Response
		expectedType  string
	}{
		{
			name: "Choice interaction",
			response: models.Response{
				RenderChoice: &models.RenderChoice{},
			},
			expectedType: "identifier",
		},
		{
			name: "Integer FIB",
			response: models.Response{
				RenderFib: &models.RenderFib{FibType: "integer"},
			},
			expectedType: "integer",
		},
		{
			name: "Decimal FIB",
			response: models.Response{
				RenderFib: &models.RenderFib{FibType: "decimal"},
			},
			expectedType: "float",
		},
		{
			name: "Default string FIB",
			response: models.Response{
				RenderFib: &models.RenderFib{},
			},
			expectedType: "string",
		},
		{
			name: "Default string",
			response: models.Response{},
			expectedType: "string",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseType := m.determineBaseType(&tc.response)
			if baseType != tc.expectedType {
				t.Errorf("Expected base type '%s', got '%s'", tc.expectedType, baseType)
			}
		})
	}
}

func TestMigrator12to21_ExtractCorrectResponse(t *testing.T) {
	m := New()
	responseProc := &models.ResponseProc{
		ResCondition: []models.ResCondition{
			{
				ConditionVar: &models.ConditionVar{
					VarEqual: []models.VarEqual{
						{RespIdent: "RESPONSE", Value: "A"},
					},
				},
				SetVar: []models.SetVar{
					{Action: "set", VarName: "SCORE", Value: "1"},
				},
			},
			{
				ConditionVar: &models.ConditionVar{
					VarEqual: []models.VarEqual{
						{RespIdent: "RESPONSE", Value: "B"},
					},
				},
				SetVar: []models.SetVar{
					{Action: "set", VarName: "SCORE", Value: "0"},
				},
			},
		},
	}
	
	correctResponse := m.extractCorrectResponse("RESPONSE", responseProc)
	
	if correctResponse == nil {
		t.Fatal("Expected correct response to be found")
	}
	
	if len(correctResponse.Value) != 1 {
		t.Errorf("Expected 1 correct value, got %d", len(correctResponse.Value))
	}
	
	if correctResponse.Value[0] != "A" {
		t.Errorf("Expected correct value 'A', got '%s'", correctResponse.Value[0])
	}
}

func TestMigrator12to21_ConvertVarType(t *testing.T) {
	m := New()
	
	testCases := []struct {
		input    string
		expected string
	}{
		{"integer", "integer"},
		{"decimal", "float"},
		{"scientific", "float"},
		{"boolean", "boolean"},
		{"unknown", "float"},
		{"", "float"},
	}
	
	for _, tc := range testCases {
		result := m.convertVarType(tc.input)
		if result != tc.expected {
			t.Errorf("convertVarType(%s): expected '%s', got '%s'", tc.input, tc.expected, result)
		}
	}
}

func TestMigrator12to21_SanitizeHTMLContent(t *testing.T) {
	m := New()
	
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Self-closing br tag",
			input:    "Line 1<br>Line 2",
			expected: "Line 1<br/>Line 2",
		},
		{
			name:     "Self-closing hr tag",
			input:    "Before<hr>After",
			expected: "Before<hr/>After",
		},
		{
			name:     "Img tag without closing",
			input:    `<img src="test.jpg">`,
			expected: `<img src="test.jpg"/>`,
		},
		{
			name:     "Already properly closed",
			input:    `<img src="test.jpg"/>`,
			expected: `<img src="test.jpg"/>`,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := m.sanitizeHTMLContent(tc.input)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestMigrator12to21_ExtractMaterialContent(t *testing.T) {
	m := New()
	material := &models.Material{
		MatText: []models.MatText{
			{Content: "Text content", TextType: "text/plain"},
			{Content: "<p>HTML content</p>", TextType: "text/html"},
		},
		MatImage: []models.MatImage{
			{URI: "image1.jpg", Width: 100, Height: 50},
			{URI: "image2.png"},
		},
	}
	
	content := m.extractMaterialContent(material)
	
	if !strings.Contains(content, "Text content") {
		t.Error("Expected plain text content")
	}
	
	if !strings.Contains(content, "<p>HTML content</p>") {
		t.Error("Expected HTML content")
	}
	
	if !strings.Contains(content, `<img src="image1.jpg" width="100" height="50" />`) {
		t.Error("Expected first image with dimensions")
	}
	
	if !strings.Contains(content, `<img src="image2.png" />`) {
		t.Error("Expected second image without dimensions")
	}
}

func TestMigrator12to21_ComplexDocumentStructure(t *testing.T) {
	doc := &models.QTIDocument{
		Version: "1.2",
		Assessment: &models.Assessment{
			Ident: "test001",
			Title: "Test Assessment",
			Sections: []models.Section{
				{
					Ident: "sec001",
					Title: "Section 1",
					Items: []models.Item{
						{
							Ident: "q001",
							Title: "Question 1",
						},
					},
				},
			},
		},
		Items: []models.Item{
			{
				Ident: "standalone_q001",
				Title: "Standalone Question",
			},
		},
	}
	
	m := New()
	result, err := m.Migrate(doc)
	
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}
	
	resultStr := string(result)
	
	// Check that assessment structure is preserved
	if !strings.Contains(resultStr, `ident="test001"`) {
		t.Error("Expected assessment ident to be preserved")
	}
	
	if !strings.Contains(resultStr, `ident="sec001"`) {
		t.Error("Expected section ident to be preserved")
	}
	
	// Check both standalone items and items in sections are processed
	if !strings.Contains(resultStr, `ident="q001"`) {
		t.Error("Expected section item to be preserved")
	}
	
	if !strings.Contains(resultStr, `ident="standalone_q001"`) {
		t.Error("Expected standalone item to be preserved")
	}
}

func BenchmarkMigrator12to21_Migrate(b *testing.B) {
	doc := &models.QTIDocument{
		Version: "1.2",
		Items: []models.Item{
			{
				Ident: "q001",
				Title: "Test Question",
				Presentation: &models.Presentation{
					Material: &models.Material{
						MatText: []models.MatText{{Content: "What is 2 + 2?"}},
					},
					Response: []models.Response{
						{
							Ident:        "RESPONSE",
							RCardinality: "single",
							RenderChoice: &models.RenderChoice{
								ResponseLabel: []models.ResponseLabel{
									{Ident: "A", Material: &models.Material{MatText: []models.MatText{{Content: "3"}}}},
									{Ident: "B", Material: &models.Material{MatText: []models.MatText{{Content: "4"}}}},
								},
							},
						},
					},
				},
			},
		},
	}
	
	m := New()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.Migrate(doc)
		if err != nil {
			b.Fatalf("Migration failed: %v", err)
		}
	}
}