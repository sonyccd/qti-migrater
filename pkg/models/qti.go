package models

import (
	"encoding/xml"
)

// QTIDocument is a generic structure for backward compatibility
// It's a union of all version-specific fields to support existing code
// New code should use version-specific structures (QTIDocument12, QTIDocument21, QTIDocument30)
type QTIDocument struct {
	XMLName    xml.Name    `xml:"questestinterop"`
	Version    string      `xml:"version,attr"`
	Items      []Item      `xml:"item"`
	Assessment *Assessment `xml:"assessment,omitempty"`
	Metadata   *Metadata   `xml:"metadata,omitempty"`
}

type Assessment struct {
	XMLName     xml.Name     `xml:"assessment"`
	Title       string       `xml:"title,attr"`
	Ident       string       `xml:"ident,attr"`
	Sections    []Section    `xml:"section"`
	Metadata    *Metadata    `xml:"metadata,omitempty"`
	Objectives  []Objective  `xml:"objectives>objective,omitempty"`
	RubricBlock *RubricBlock `xml:"rubricBlock,omitempty"`
}

type Section struct {
	XMLName  xml.Name  `xml:"section"`
	Title    string    `xml:"title,attr"`
	Ident    string    `xml:"ident,attr"`
	Items    []Item    `xml:"item"`
	Metadata *Metadata `xml:"metadata,omitempty"`
}

// Item is a generic structure combining all version fields for backward compatibility
// Use Item12, Item21, or Item30 for version-specific parsing
type Item struct {
	XMLName        xml.Name        `xml:"item"`
	Title          string          `xml:"title,attr"`
	Ident          string          `xml:"ident,attr"`
	MaxAttempts    int             `xml:"maxattempts,attr,omitempty"`
	Metadata       *Metadata       `xml:"metadata,omitempty"`
	// QTI 1.2/2.1 fields
	Presentation   *Presentation   `xml:"presentation,omitempty"`
	ResponseProc   *ResponseProc   `xml:"resprocessing,omitempty"`
	// QTI 2.1/3.0 fields
	ItemBody       *ItemBody       `xml:"itemBody,omitempty"`
	ResponseDecl   []ResponseDecl  `xml:"responseDeclaration,omitempty"`
	OutcomeDecl    []OutcomeDecl   `xml:"outcomeDeclaration,omitempty"`
	TemplateDecl   []TemplateDecl  `xml:"templateDeclaration,omitempty"`
	Feedback       []Feedback      `xml:"itemfeedback,omitempty"`
	RubricBlock    *RubricBlock    `xml:"rubricBlock,omitempty"`
}

// Legacy types kept for backward compatibility
// These are generic versions that combine fields from multiple QTI versions
// For version-specific parsing, use the types in qti12.go, qti21.go, or qti30.go

// Legacy types for backward compatibility
// Note: Import these from version-specific files when needed
// For now, we'll define minimal types here to avoid conflicts

// Generic types that combine fields from multiple versions
type Response = Response12
type Flow = Flow12
type ItemBody = ItemBody21
type P = P21
type Div = Div21
type ChoiceInteraction = ChoiceInteraction21
type SimpleChoice = SimpleChoice21
type Prompt = Prompt21
type TextEntryInteraction = TextEntryInteraction21
type ExtendedTextInteraction = ExtendedTextInteraction21
type ResponseDecl = ResponseDecl21
type CorrectResponse = CorrectResponse21
type Mapping = Mapping21
type MapEntry = MapEntry21
type OutcomeDecl = OutcomeDecl21
type DefaultValue = DefaultValue21
type TemplateDecl = TemplateDecl21
type Feedback = Feedback21