package models

import "encoding/xml"

// QTI 3.0 specific structures
// QTI 3.0 represents a major overhaul with cleaner, more modern XML structure

type QTIDocument30 struct {
	XMLName    xml.Name      `xml:"qtiAssessmentItem"`
	Version    string        `xml:"version,attr"`
	Identifier string        `xml:"identifier,attr"`
	Title      string        `xml:"title,attr"`
	TimeDependent bool       `xml:"timeDependent,attr,omitempty"`
	Adaptive   bool          `xml:"adaptive,attr,omitempty"`
	// Declarations come first in QTI 3.0
	ResponseDeclarations []ResponseDecl30 `xml:"responseDeclaration,omitempty"`
	OutcomeDeclarations  []OutcomeDecl30  `xml:"outcomeDeclaration,omitempty"`
	TemplateDeclarations []TemplateDecl30 `xml:"templateDeclaration,omitempty"`
	// Content
	ItemBody       *ItemBody30      `xml:"itemBody"`
	ResponseProcessing *ResponseProcessing30 `xml:"responseProcessing,omitempty"`
	ModalFeedback  []ModalFeedback30 `xml:"modalFeedback,omitempty"`
	// Metadata
	Metadata       *Metadata        `xml:"metadata,omitempty"`
}

// QTI 3.0 uses a different root for assessments
type Assessment30 struct {
	XMLName    xml.Name      `xml:"assessmentTest"`
	Identifier string        `xml:"identifier,attr"`
	Title      string        `xml:"title,attr"`
	TestParts  []TestPart30  `xml:"testPart"`
	Metadata   *Metadata     `xml:"metadata,omitempty"`
}

type TestPart30 struct {
	XMLName        xml.Name          `xml:"testPart"`
	Identifier     string            `xml:"identifier,attr"`
	NavigationMode string            `xml:"navigationMode,attr,omitempty"`
	SubmissionMode string            `xml:"submissionMode,attr,omitempty"`
	Sections       []AssessmentSection30 `xml:"assessmentSection"`
}

type AssessmentSection30 struct {
	XMLName    xml.Name   `xml:"assessmentSection"`
	Identifier string     `xml:"identifier,attr"`
	Title      string     `xml:"title,attr"`
	Visible    bool       `xml:"visible,attr,omitempty"`
	ItemRefs   []ItemRef30 `xml:"assessmentItemRef"`
	Metadata   *Metadata  `xml:"metadata,omitempty"`
}

type ItemRef30 struct {
	XMLName    xml.Name `xml:"assessmentItemRef"`
	Identifier string   `xml:"identifier,attr"`
	Href       string   `xml:"href,attr"`
	Category   []string `xml:"category,attr,omitempty"`
}

// QTI 3.0 ItemBody - cleaner structure than 2.x

type ItemBody30 struct {
	XMLName xml.Name `xml:"itemBody"`
	// Content can include various interactions and block elements
	Content []interface{} `xml:",any"`
}

// QTI 3.0 Interactions - more structured than previous versions

type ChoiceInteraction30 struct {
	XMLName               xml.Name        `xml:"choiceInteraction"`
	ResponseIdentifier    string          `xml:"responseIdentifier,attr"`
	Shuffle               bool            `xml:"shuffle,attr,omitempty"`
	MaxChoices            int             `xml:"maxChoices,attr,omitempty"`
	MinChoices            int             `xml:"minChoices,attr,omitempty"`
	Orientation           string          `xml:"orientation,attr,omitempty"`
	Prompt                *Prompt30       `xml:"prompt,omitempty"`
	SimpleChoice          []SimpleChoice30 `xml:"simpleChoice"`
}

type SimpleChoice30 struct {
	XMLName     xml.Name `xml:"simpleChoice"`
	Identifier  string   `xml:"identifier,attr"`
	Fixed       bool     `xml:"fixed,attr,omitempty"`
	ShowHide    string   `xml:"showHide,attr,omitempty"`
	Content     string   `xml:",innerxml"`
}

type TextEntryInteraction30 struct {
	XMLName            xml.Name `xml:"textEntryInteraction"`
	ResponseIdentifier string   `xml:"responseIdentifier,attr"`
	Base               int      `xml:"base,attr,omitempty"`
	StringIdentifier   string   `xml:"stringIdentifier,attr,omitempty"`
	ExpectedLength     int      `xml:"expectedLength,attr,omitempty"`
	PatternMask        string   `xml:"patternMask,attr,omitempty"`
	PlaceholderText    string   `xml:"placeholderText,attr,omitempty"`
}

type ExtendedTextInteraction30 struct {
	XMLName            xml.Name  `xml:"extendedTextInteraction"`
	ResponseIdentifier string    `xml:"responseIdentifier,attr"`
	Base               int       `xml:"base,attr,omitempty"`
	StringIdentifier   string    `xml:"stringIdentifier,attr,omitempty"`
	ExpectedLength     int       `xml:"expectedLength,attr,omitempty"`
	PatternMask        string    `xml:"patternMask,attr,omitempty"`
	PlaceholderText    string    `xml:"placeholderText,attr,omitempty"`
	MaxStrings         int       `xml:"maxStrings,attr,omitempty"`
	MinStrings         int       `xml:"minStrings,attr,omitempty"`
	ExpectedLines      int       `xml:"expectedLines,attr,omitempty"`
	Format             string    `xml:"format,attr,omitempty"`
	Prompt             *Prompt30 `xml:"prompt,omitempty"`
}

type Prompt30 struct {
	XMLName xml.Name `xml:"prompt"`
	Content string   `xml:",innerxml"`
}

// QTI 3.0 Declarations - more structured than 2.x

type ResponseDecl30 struct {
	XMLName         xml.Name           `xml:"responseDeclaration"`
	Identifier      string             `xml:"identifier,attr"`
	Cardinality     string             `xml:"cardinality,attr"`
	BaseType        string             `xml:"baseType,attr,omitempty"`
	CorrectResponse *CorrectResponse30 `xml:"correctResponse,omitempty"`
	Mapping         *Mapping30         `xml:"mapping,omitempty"`
	AreaMapping     *AreaMapping30     `xml:"areaMapping,omitempty"`
}

type CorrectResponse30 struct {
	XMLName xml.Name  `xml:"correctResponse"`
	Value   []Value30 `xml:"value"`
}

type Value30 struct {
	XMLName    xml.Name `xml:"value"`
	FieldIdentifier string `xml:"fieldIdentifier,attr,omitempty"`
	BaseType   string   `xml:"baseType,attr,omitempty"`
	Content    string   `xml:",chardata"`
}

type Mapping30 struct {
	XMLName      xml.Name      `xml:"mapping"`
	LowerBound   float64       `xml:"lowerBound,attr,omitempty"`
	UpperBound   float64       `xml:"upperBound,attr,omitempty"`
	DefaultValue float64       `xml:"defaultValue,attr"`
	MapEntry     []MapEntry30  `xml:"mapEntry"`
}

type MapEntry30 struct {
	XMLName     xml.Name `xml:"mapEntry"`
	MapKey      string   `xml:"mapKey,attr"`
	MappedValue float64  `xml:"mappedValue,attr"`
}

type AreaMapping30 struct {
	XMLName          xml.Name         `xml:"areaMapping"`
	LowerBound       float64          `xml:"lowerBound,attr,omitempty"`
	UpperBound       float64          `xml:"upperBound,attr,omitempty"`
	DefaultValue     float64          `xml:"defaultValue,attr"`
	AreaMapEntry     []AreaMapEntry30 `xml:"areaMapEntry"`
}

type AreaMapEntry30 struct {
	XMLName     xml.Name `xml:"areaMapEntry"`
	Shape       string   `xml:"shape,attr"`
	Coords      string   `xml:"coords,attr"`
	MappedValue float64  `xml:"mappedValue,attr"`
}

type OutcomeDecl30 struct {
	XMLName         xml.Name         `xml:"outcomeDeclaration"`
	Identifier      string           `xml:"identifier,attr"`
	Cardinality     string           `xml:"cardinality,attr"`
	BaseType        string           `xml:"baseType,attr,omitempty"`
	View            []string         `xml:"view,attr,omitempty"`
	Interpretation  string           `xml:"interpretation,attr,omitempty"`
	LongInterpretation string        `xml:"longInterpretation,attr,omitempty"`
	NormalMaximum   float64          `xml:"normalMaximum,attr,omitempty"`
	NormalMinimum   float64          `xml:"normalMinimum,attr,omitempty"`
	MasteryValue    float64          `xml:"masteryValue,attr,omitempty"`
	DefaultValue    *DefaultValue30  `xml:"defaultValue,omitempty"`
}

type DefaultValue30 struct {
	XMLName xml.Name  `xml:"defaultValue"`
	Value   []Value30 `xml:"value"`
}

type TemplateDecl30 struct {
	XMLName       xml.Name        `xml:"templateDeclaration"`
	Identifier    string          `xml:"identifier,attr"`
	Cardinality   string          `xml:"cardinality,attr"`
	BaseType      string          `xml:"baseType,attr,omitempty"`
	ParamVariable bool            `xml:"paramVariable,attr,omitempty"`
	MathVariable  bool            `xml:"mathVariable,attr,omitempty"`
	DefaultValue  *DefaultValue30 `xml:"defaultValue,omitempty"`
}

// QTI 3.0 Response Processing - completely different from resprocessing

type ResponseProcessing30 struct {
	XMLName             xml.Name              `xml:"responseProcessing"`
	Template            string                `xml:"template,attr,omitempty"`
	TemplateLocation    string                `xml:"templateLocation,attr,omitempty"`
	ResponseRules       []interface{}         `xml:",any"` // Can contain various rule types
}

// QTI 3.0 Modal Feedback

type ModalFeedback30 struct {
	XMLName    xml.Name `xml:"modalFeedback"`
	Identifier string   `xml:"identifier,attr"`
	OutcomeIdentifier string `xml:"outcomeIdentifier,attr"`
	ShowHide   string   `xml:"showHide,attr"`
	Title      string   `xml:"title,attr,omitempty"`
	Content    string   `xml:",innerxml"`
}