package models

import "encoding/xml"

// QTI 2.1/2.2 specific structures
// Note: QTI 2.1 and 2.2 share the same structure with minor differences in behavior

type QTIDocument21 struct {
	XMLName    xml.Name    `xml:"questestinterop"`
	Version    string      `xml:"version,attr"`
	Items      []Item21    `xml:"item"`
	Assessment *Assessment21 `xml:"assessment,omitempty"`
	Metadata   *Metadata   `xml:"metadata,omitempty"`
}

type Assessment21 struct {
	XMLName     xml.Name     `xml:"assessment"`
	Title       string       `xml:"title,attr"`
	Ident       string       `xml:"ident,attr"`
	Sections    []Section21  `xml:"section"`
	Metadata    *Metadata    `xml:"metadata,omitempty"`
	Objectives  []Objective  `xml:"objectives>objective,omitempty"`
	RubricBlock *RubricBlock `xml:"rubricBlock,omitempty"`
}

type Section21 struct {
	XMLName  xml.Name  `xml:"section"`
	Title    string    `xml:"title,attr"`
	Ident    string    `xml:"ident,attr"`
	Items    []Item21  `xml:"item"`
	Metadata *Metadata `xml:"metadata,omitempty"`
}

// QTI 2.1/2.2 uses a hybrid structure that supports both 1.2 legacy and newer elements
type Item21 struct {
	XMLName        xml.Name        `xml:"item"`
	Title          string          `xml:"title,attr"`
	Ident          string          `xml:"ident,attr"`
	MaxAttempts    int             `xml:"maxattempts,attr,omitempty"`
	Metadata       *Metadata       `xml:"metadata,omitempty"`
	// Legacy 1.2 structures still supported
	Presentation   *Presentation   `xml:"presentation,omitempty"`
	ResponseProc   *ResponseProc   `xml:"resprocessing,omitempty"`
	// Newer 2.x structures
	ItemBody       *ItemBody21     `xml:"itemBody,omitempty"`
	ResponseDecl   []ResponseDecl21  `xml:"responseDeclaration,omitempty"`
	OutcomeDecl    []OutcomeDecl21   `xml:"outcomeDeclaration,omitempty"`
	TemplateDecl   []TemplateDecl21  `xml:"templateDeclaration,omitempty"`
	Feedback       []Feedback21    `xml:"itemfeedback,omitempty"`
	RubricBlock    *RubricBlock    `xml:"rubricBlock,omitempty"`
}

// QTI 2.1/2.2 ItemBody structures

type ItemBody21 struct {
	XMLName     xml.Name     `xml:"itemBody"`
	P           []P21        `xml:"p,omitempty"`
	Div         []Div21      `xml:"div,omitempty"`
	ChoiceInteraction []ChoiceInteraction21 `xml:"choiceInteraction,omitempty"`
	TextEntryInteraction []TextEntryInteraction21 `xml:"textEntryInteraction,omitempty"`
	ExtendedTextInteraction []ExtendedTextInteraction21 `xml:"extendedTextInteraction,omitempty"`
}

type P21 struct {
	XMLName xml.Name `xml:"p"`
	Content string   `xml:",innerxml"`
}

type Div21 struct {
	XMLName xml.Name `xml:"div"`
	Class   string   `xml:"class,attr,omitempty"`
	Content string   `xml:",innerxml"`
}

type ChoiceInteraction21 struct {
	XMLName         xml.Name        `xml:"choiceInteraction"`
	ResponseIdent   string          `xml:"responseIdentifier,attr"`
	Shuffle         bool            `xml:"shuffle,attr,omitempty"`
	MaxChoices      int             `xml:"maxChoices,attr,omitempty"`
	MinChoices      int             `xml:"minChoices,attr,omitempty"`
	Prompt          *Prompt21       `xml:"prompt,omitempty"`
	SimpleChoice    []SimpleChoice21 `xml:"simpleChoice"`
}

type SimpleChoice21 struct {
	XMLName     xml.Name `xml:"simpleChoice"`
	Identifier  string   `xml:"identifier,attr"`
	Fixed       bool     `xml:"fixed,attr,omitempty"`
	Content     string   `xml:",innerxml"`
}

type Prompt21 struct {
	XMLName xml.Name `xml:"prompt"`
	Content string   `xml:",innerxml"`
}

type TextEntryInteraction21 struct {
	XMLName         xml.Name `xml:"textEntryInteraction"`
	ResponseIdent   string   `xml:"responseIdentifier,attr"`
	ExpectedLength  int      `xml:"expectedLength,attr,omitempty"`
	PatternMask     string   `xml:"patternMask,attr,omitempty"`
	PlaceholderText string   `xml:"placeholderText,attr,omitempty"`
}

type ExtendedTextInteraction21 struct {
	XMLName         xml.Name `xml:"extendedTextInteraction"`
	ResponseIdent   string   `xml:"responseIdentifier,attr"`
	MinStrings      int      `xml:"minStrings,attr,omitempty"`
	MaxStrings      int      `xml:"maxStrings,attr,omitempty"`
	ExpectedLines   int      `xml:"expectedLines,attr,omitempty"`
	ExpectedLength  int      `xml:"expectedLength,attr,omitempty"`
	Prompt          *Prompt21 `xml:"prompt,omitempty"`
}

// QTI 2.1/2.2 Declaration structures

type ResponseDecl21 struct {
	XMLName       xml.Name      `xml:"responseDeclaration"`
	Identifier    string        `xml:"identifier,attr"`
	Cardinality   string        `xml:"cardinality,attr"`
	BaseType      string        `xml:"baseType,attr,omitempty"`
	CorrectResponse *CorrectResponse21 `xml:"correctResponse,omitempty"`
	Mapping       *Mapping21    `xml:"mapping,omitempty"`
}

type CorrectResponse21 struct {
	XMLName xml.Name `xml:"correctResponse"`
	Value   []string `xml:"value"`
}

type Mapping21 struct {
	XMLName        xml.Name    `xml:"mapping"`
	LowerBound     float64     `xml:"lowerBound,attr,omitempty"`
	UpperBound     float64     `xml:"upperBound,attr,omitempty"`
	DefaultValue   float64     `xml:"defaultValue,attr,omitempty"`
	MapEntry       []MapEntry21 `xml:"mapEntry"`
}

type MapEntry21 struct {
	XMLName    xml.Name `xml:"mapEntry"`
	MapKey     string   `xml:"mapKey,attr"`
	MappedValue float64 `xml:"mappedValue,attr"`
}

type OutcomeDecl21 struct {
	XMLName         xml.Name        `xml:"outcomeDeclaration"`
	Identifier      string          `xml:"identifier,attr"`
	Cardinality     string          `xml:"cardinality,attr"`
	BaseType        string          `xml:"baseType,attr,omitempty"`
	DefaultValue    *DefaultValue21 `xml:"defaultValue,omitempty"`
}

type DefaultValue21 struct {
	XMLName xml.Name `xml:"defaultValue"`
	Value   string   `xml:"value"`
}

type TemplateDecl21 struct {
	XMLName      xml.Name     `xml:"templateDeclaration"`
	Identifier   string       `xml:"identifier,attr"`
	Cardinality  string       `xml:"cardinality,attr"`
	BaseType     string       `xml:"baseType,attr,omitempty"`
	ParamVariable bool        `xml:"paramVariable,attr,omitempty"`
	DefaultValue *DefaultValue21 `xml:"defaultValue,omitempty"`
}

// QTI 2.1/2.2 Feedback structures
// Note: Can use both legacy format and newer format

type Feedback21 struct {
	XMLName        xml.Name       `xml:"itemfeedback"`
	Ident          string         `xml:"ident,attr"`
	Title          string         `xml:"title,attr,omitempty"`
	FlowMat        []FlowMat      `xml:"flow_mat,omitempty"` // Legacy 1.2 style
	Material       *Material      `xml:"material,omitempty"`  // Can be used directly
}