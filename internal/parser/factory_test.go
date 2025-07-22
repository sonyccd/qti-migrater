package parser

import (
	"strings"
	"testing"
)

func TestGetParser_QTI12(t *testing.T) {
	testCases := []string{"1.2", "1.2.0", "1.2.1"}
	
	for _, version := range testCases {
		parser, err := GetParser(version)
		if err != nil {
			t.Errorf("Failed to get parser for version %s: %v", version, err)
			continue
		}
		
		if parser.Version() != "1.2" {
			t.Errorf("Expected parser version '1.2', got '%s' for input version '%s'", 
				parser.Version(), version)
		}
	}
}

func TestGetParser_QTI21(t *testing.T) {
	testCases := []string{"2.1", "2.1.0", "2.1.1", "2.2", "2.2.0"}
	
	for _, version := range testCases {
		parser, err := GetParser(version)
		if err != nil {
			t.Errorf("Failed to get parser for version %s: %v", version, err)
			continue
		}
		
		if parser.Version() != "2.1" {
			t.Errorf("Expected parser version '2.1', got '%s' for input version '%s'", 
				parser.Version(), version)
		}
	}
}

func TestGetParser_UnsupportedVersion(t *testing.T) {
	unsupportedVersions := []string{"1.0", "1.1", "3.0", "4.0", "invalid", ""}
	
	for _, version := range unsupportedVersions {
		parser, err := GetParser(version)
		
		if err == nil {
			t.Errorf("Expected error for unsupported version '%s', but got parser: %v", 
				version, parser)
			continue
		}
		
		if !strings.Contains(err.Error(), "unsupported QTI version") {
			t.Errorf("Expected 'unsupported QTI version' error for version '%s', got: %v", 
				version, err)
		}
		
		if parser != nil {
			t.Errorf("Expected nil parser for unsupported version '%s', got: %v", 
				version, parser)
		}
	}
}

func TestGetParser_WhitespaceHandling(t *testing.T) {
	testCases := map[string]string{
		"  1.2  ":   "1.2",
		"\t2.1\n":   "2.1", 
		" 2.2.0 ":   "2.1",
		"1.2.1\t":   "1.2",
	}
	
	for input, expectedVersion := range testCases {
		parser, err := GetParser(input)
		if err != nil {
			t.Errorf("Failed to get parser for version '%s': %v", input, err)
			continue
		}
		
		if parser.Version() != expectedVersion {
			t.Errorf("Expected parser version '%s', got '%s' for input version '%s'", 
				expectedVersion, parser.Version(), input)
		}
	}
}

func TestGetParser_EdgeCases(t *testing.T) {
	// Test case sensitivity (versions should be case-sensitive)
	_, err := GetParser("1.2.0")
	if err != nil {
		t.Errorf("Expected success for version '1.2.0', got error: %v", err)
	}
	
	// Test that we correctly handle version prefixes
	parser12, err := GetParser("1.2.5")  // Should work as it starts with 1.2
	if err != nil {
		t.Errorf("Expected success for version '1.2.5', got error: %v", err)
	}
	if parser12.Version() != "1.2" {
		t.Errorf("Expected parser version '1.2', got '%s'", parser12.Version())
	}
	
	parser21, err := GetParser("2.1.5")  // Should work as it starts with 2.1
	if err != nil {
		t.Errorf("Expected success for version '2.1.5', got error: %v", err)
	}
	if parser21.Version() != "2.1" {
		t.Errorf("Expected parser version '2.1', got '%s'", parser21.Version())
	}
	
	// Test that we handle 2.2 versions correctly (should use 2.1 parser)
	parser22, err := GetParser("2.2.3")
	if err != nil {
		t.Errorf("Expected success for version '2.2.3', got error: %v", err)
	}
	if parser22.Version() != "2.1" {
		t.Errorf("Expected parser version '2.1' for QTI 2.2 document, got '%s'", parser22.Version())
	}
}

func BenchmarkGetParser_QTI12(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetParser("1.2")
		if err != nil {
			b.Fatalf("GetParser failed: %v", err)
		}
	}
}

func BenchmarkGetParser_QTI21(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetParser("2.1")
		if err != nil {
			b.Fatalf("GetParser failed: %v", err)
		}
	}
}