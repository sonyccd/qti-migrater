package models

import (
	"encoding/xml"
	"time"
)

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

type Item struct {
	XMLName        xml.Name        `xml:"item"`
	Title          string          `xml:"title,attr"`
	Ident          string          `xml:"ident,attr"`
	MaxAttempts    int             `xml:"maxattempts,attr,omitempty"`
	Metadata       *Metadata       `xml:"metadata,omitempty"`
	Presentation   *Presentation   `xml:"presentation,omitempty"`
	ResponseProc   *ResponseProc   `xml:"resprocessing,omitempty"`
	ItemBody       *ItemBody       `xml:"itemBody,omitempty"`
	ResponseDecl   []ResponseDecl  `xml:"responseDeclaration,omitempty"`
	OutcomeDecl    []OutcomeDecl   `xml:"outcomeDeclaration,omitempty"`
	TemplateDecl   []TemplateDecl  `xml:"templateDeclaration,omitempty"`
	Feedback       []Feedback      `xml:"itemfeedback,omitempty"`
	RubricBlock    *RubricBlock    `xml:"rubricBlock,omitempty"`
}

type Metadata struct {
	XMLName     xml.Name    `xml:"metadata"`
	Schema      string      `xml:"schema,omitempty"`
	SchemaVer   string      `xml:"schemaversion,omitempty"`
	LOM         *LOM        `xml:"lom,omitempty"`
	QTIMetadata *QTIMetadata `xml:"qtimetadata,omitempty"`
}

type LOM struct {
	XMLName     xml.Name     `xml:"lom"`
	General     *LOMGeneral  `xml:"general,omitempty"`
	Lifecycle   *LOMLifecycle `xml:"lifecycle,omitempty"`
	Technical   *LOMTechnical `xml:"technical,omitempty"`
	Educational *LOMEducational `xml:"educational,omitempty"`
	Rights      *LOMRights   `xml:"rights,omitempty"`
}

type LOMGeneral struct {
	Identifier  []LOMIdentifier `xml:"identifier,omitempty"`
	Title       *LOMString      `xml:"title,omitempty"`
	Language    []string        `xml:"language,omitempty"`
	Description []LOMString     `xml:"description,omitempty"`
	Keyword     []LOMString     `xml:"keyword,omitempty"`
}

type LOMIdentifier struct {
	Catalog string `xml:"catalog,omitempty"`
	Entry   string `xml:"entry,omitempty"`
}

type LOMString struct {
	Language string `xml:"language,attr,omitempty"`
	Value    string `xml:",chardata"`
}

type LOMLifecycle struct {
	Version     *LOMString      `xml:"version,omitempty"`
	Status      *LOMVocabulary  `xml:"status,omitempty"`
	Contribute  []LOMContribute `xml:"contribute,omitempty"`
}

type LOMVocabulary struct {
	Source string `xml:"source,omitempty"`
	Value  string `xml:"value,omitempty"`
}

type LOMContribute struct {
	Role   *LOMVocabulary `xml:"role,omitempty"`
	Entity []string       `xml:"entity,omitempty"`
	Date   *LOMDateTime   `xml:"date,omitempty"`
}

type LOMDateTime struct {
	DateTime    time.Time  `xml:"datetime,omitempty"`
	Description *LOMString `xml:"description,omitempty"`
}

type LOMTechnical struct {
	Format               []string              `xml:"format,omitempty"`
	Size                 string                `xml:"size,omitempty"`
	Location             []string              `xml:"location,omitempty"`
	Requirement          []LOMRequirement      `xml:"requirement,omitempty"`
	InstallationRemarks  *LOMString            `xml:"installationremarks,omitempty"`
	OtherPlatformReq     *LOMString            `xml:"otherplatformrequirements,omitempty"`
	Duration             *LOMDuration          `xml:"duration,omitempty"`
}

type LOMRequirement struct {
	OrComposite []LOMOrComposite `xml:"orcomposite,omitempty"`
}

type LOMOrComposite struct {
	Type         *LOMVocabulary `xml:"type,omitempty"`
	Name         *LOMVocabulary `xml:"name,omitempty"`
	MinVersion   string         `xml:"minimumversion,omitempty"`
	MaxVersion   string         `xml:"maximumversion,omitempty"`
}

type LOMDuration struct {
	Duration    string     `xml:"duration,omitempty"`
	Description *LOMString `xml:"description,omitempty"`
}

type LOMEducational struct {
	InteractivityType    *LOMVocabulary  `xml:"interactivitytype,omitempty"`
	LearningResourceType []LOMVocabulary `xml:"learningresourcetype,omitempty"`
	InteractivityLevel   *LOMVocabulary  `xml:"interactivitylevel,omitempty"`
	SemanticDensity      *LOMVocabulary  `xml:"semanticdensity,omitempty"`
	IntendedEndUserRole  []LOMVocabulary `xml:"intendedenduserrole,omitempty"`
	Context              []LOMVocabulary `xml:"context,omitempty"`
	TypicalAgeRange      []LOMString     `xml:"typicalagerange,omitempty"`
	Difficulty           *LOMVocabulary  `xml:"difficulty,omitempty"`
	TypicalLearningTime  *LOMDuration    `xml:"typicallearningtime,omitempty"`
	Description          []LOMString     `xml:"description,omitempty"`
	Language             []string        `xml:"language,omitempty"`
}

type LOMRights struct {
	Cost                *LOMVocabulary `xml:"cost,omitempty"`
	CopyrightAndOther   *LOMVocabulary `xml:"copyrightandotherrestrictions,omitempty"`
	Description         *LOMString     `xml:"description,omitempty"`
}

type QTIMetadata struct {
	XMLName            xml.Name    `xml:"qtimetadata"`
	TimeDependent      bool        `xml:"timedependent,omitempty"`
	Composite          bool        `xml:"composite,omitempty"`
	InteractionType    string      `xml:"interactiontype,omitempty"`
	FeedbackType       string      `xml:"feedbacktype,omitempty"`
	SolutionAvailable  bool        `xml:"solutionavailable,omitempty"`
	Scoringmode        string      `xml:"scoringmode,omitempty"`
	ToolName           string      `xml:"toolname,omitempty"`
	ToolVersion        string      `xml:"toolversion,omitempty"`
	ToolVendor         string      `xml:"toolvendor,omitempty"`
}

type Presentation struct {
	XMLName  xml.Name  `xml:"presentation"`
	Label    string    `xml:"label,attr,omitempty"`
	Material *Material `xml:"material,omitempty"`
	Response []Response `xml:"response_lid,omitempty"`
	Flow     []Flow     `xml:"flow,omitempty"`
}

type Material struct {
	XMLName  xml.Name  `xml:"material"`
	Label    string    `xml:"label,attr,omitempty"`
	MatText  []MatText `xml:"mattext,omitempty"`
	MatImage []MatImage `xml:"matimage,omitempty"`
	MatAudio []MatAudio `xml:"mataudio,omitempty"`
	MatVideo []MatVideo `xml:"matvideo,omitempty"`
}

type MatText struct {
	XMLName     xml.Name `xml:"mattext"`
	TextType    string   `xml:"texttype,attr,omitempty"`
	Charset     string   `xml:"charset,attr,omitempty"`
	XML         string   `xml:"xml:space,attr,omitempty"`
	Content     string   `xml:",chardata"`
}

type MatImage struct {
	XMLName    xml.Name `xml:"matimage"`
	ImageType  string   `xml:"imagetype,attr,omitempty"`
	URI        string   `xml:"uri,attr"`
	Width      int      `xml:"width,attr,omitempty"`
	Height     int      `xml:"height,attr,omitempty"`
}

type MatAudio struct {
	XMLName    xml.Name `xml:"mataudio"`
	AudioType  string   `xml:"audiotype,attr,omitempty"`
	URI        string   `xml:"uri,attr"`
}

type MatVideo struct {
	XMLName    xml.Name `xml:"matvideo"`
	VideoType  string   `xml:"videotype,attr,omitempty"`
	URI        string   `xml:"uri,attr"`
	Width      int      `xml:"width,attr,omitempty"`
	Height     int      `xml:"height,attr,omitempty"`
}

type Response struct {
	XMLName      xml.Name      `xml:"response_lid"`
	Ident        string        `xml:"ident,attr"`
	RCardinality string        `xml:"rcardinality,attr,omitempty"`
	RTiming      string        `xml:"rtiming,attr,omitempty"`
	RenderChoice *RenderChoice `xml:"render_choice,omitempty"`
	RenderFib    *RenderFib    `xml:"render_fib,omitempty"`
}

type RenderChoice struct {
	XMLName    xml.Name    `xml:"render_choice"`
	Shuffle    string      `xml:"shuffle,attr,omitempty"`
	MInNumber  int         `xml:"minnumber,attr,omitempty"`
	MaxNumber  int         `xml:"maxnumber,attr,omitempty"`
	ResponseLabel []ResponseLabel `xml:"response_label"`
}

type ResponseLabel struct {
	XMLName    xml.Name   `xml:"response_label"`
	Ident      string     `xml:"ident,attr"`
	RArea      string     `xml:"rarea,attr,omitempty"`
	RRange     string     `xml:"rrange,attr,omitempty"`
	Material   *Material  `xml:"material,omitempty"`
}

type RenderFib struct {
	XMLName    xml.Name `xml:"render_fib"`
	Encoding   string   `xml:"encoding,attr,omitempty"`
	FibType    string   `xml:"fibtype,attr,omitempty"`
	Rows       int      `xml:"rows,attr,omitempty"`
	MaxChars   int      `xml:"maxchars,attr,omitempty"`
	Prompt     string   `xml:"prompt,attr,omitempty"`
	Columns    int      `xml:"columns,attr,omitempty"`
}

type Flow struct {
	XMLName    xml.Name    `xml:"flow"`
	Class      string      `xml:"class,attr,omitempty"`
	Material   []Material  `xml:"material,omitempty"`
	Response   []Response  `xml:"response_lid,omitempty"`
	Flow       []Flow      `xml:"flow,omitempty"`
}

type ResponseProc struct {
	XMLName    xml.Name    `xml:"resprocessing"`
	ScoreModel string      `xml:"scoremodel,attr,omitempty"`
	Outcomes   *Outcomes   `xml:"outcomes,omitempty"`
	ResCondition []ResCondition `xml:"respcondition"`
}

type Outcomes struct {
	XMLName     xml.Name     `xml:"outcomes"`
	DecVar      []DecVar     `xml:"decvar"`
}

type DecVar struct {
	XMLName     xml.Name `xml:"decvar"`
	VarName     string   `xml:"varname,attr"`
	VarType     string   `xml:"vartype,attr,omitempty"`
	DefaultVal  string   `xml:"defaultval,attr,omitempty"`
	MinValue    string   `xml:"minvalue,attr,omitempty"`
	MaxValue    string   `xml:"maxvalue,attr,omitempty"`
}

type ResCondition struct {
	XMLName     xml.Name     `xml:"respcondition"`
	Title       string       `xml:"title,attr,omitempty"`
	Continue    string       `xml:"continue,attr,omitempty"`
	ConditionVar *ConditionVar `xml:"conditionvar"`
	SetVar      []SetVar     `xml:"setvar,omitempty"`
	DisplayFeedback []DisplayFeedback `xml:"displayfeedback,omitempty"`
}

type ConditionVar struct {
	XMLName     xml.Name     `xml:"conditionvar"`
	Not         *Not         `xml:"not,omitempty"`
	And         *And         `xml:"and,omitempty"`
	Or          *Or          `xml:"or,omitempty"`
	VarEqual    []VarEqual   `xml:"varequal,omitempty"`
	VarLT       []VarLT      `xml:"varlt,omitempty"`
	VarLTE      []VarLTE     `xml:"varlte,omitempty"`
	VarGT       []VarGT      `xml:"vargt,omitempty"`
	VarGTE      []VarGTE     `xml:"vargte,omitempty"`
	VarSubset   []VarSubset  `xml:"varsubset,omitempty"`
	VarInside   []VarInside  `xml:"varinside,omitempty"`
	VarSubstring []VarSubstring `xml:"varsubstring,omitempty"`
}

type Not struct {
	XMLName     xml.Name     `xml:"not"`
	VarEqual    []VarEqual   `xml:"varequal,omitempty"`
	And         *And         `xml:"and,omitempty"`
	Or          *Or          `xml:"or,omitempty"`
}

type And struct {
	XMLName     xml.Name     `xml:"and"`
	VarEqual    []VarEqual   `xml:"varequal,omitempty"`
	Not         *Not         `xml:"not,omitempty"`
	Or          *Or          `xml:"or,omitempty"`
}

type Or struct {
	XMLName     xml.Name     `xml:"or"`
	VarEqual    []VarEqual   `xml:"varequal,omitempty"`
	Not         *Not         `xml:"not,omitempty"`
	And         *And         `xml:"and,omitempty"`
}

type VarEqual struct {
	XMLName     xml.Name `xml:"varequal"`
	RespIdent   string   `xml:"respident,attr"`
	Case        string   `xml:"case,attr,omitempty"`
	Value       string   `xml:",chardata"`
}

type VarLT struct {
	XMLName     xml.Name `xml:"varlt"`
	RespIdent   string   `xml:"respident,attr"`
	Value       string   `xml:",chardata"`
}

type VarLTE struct {
	XMLName     xml.Name `xml:"varlte"`
	RespIdent   string   `xml:"respident,attr"`
	Value       string   `xml:",chardata"`
}

type VarGT struct {
	XMLName     xml.Name `xml:"vargt"`
	RespIdent   string   `xml:"respident,attr"`
	Value       string   `xml:",chardata"`
}

type VarGTE struct {
	XMLName     xml.Name `xml:"vargte"`
	RespIdent   string   `xml:"respident,attr"`
	Value       string   `xml:",chardata"`
}

type VarSubset struct {
	XMLName     xml.Name `xml:"varsubset"`
	RespIdent   string   `xml:"respident,attr"`
	SetMatch    string   `xml:"setmatch,attr,omitempty"`
	Value       string   `xml:",chardata"`
}

type VarInside struct {
	XMLName     xml.Name `xml:"varinside"`
	RespIdent   string   `xml:"respident,attr"`
	AreaMatch   string   `xml:"areamatch,attr,omitempty"`
	Value       string   `xml:",chardata"`
}

type VarSubstring struct {
	XMLName     xml.Name `xml:"varsubstring"`
	RespIdent   string   `xml:"respident,attr"`
	Case        string   `xml:"case,attr,omitempty"`
	Value       string   `xml:",chardata"`
}

type SetVar struct {
	XMLName     xml.Name `xml:"setvar"`
	Action      string   `xml:"action,attr"`
	VarName     string   `xml:"varname,attr,omitempty"`
	Value       string   `xml:",chardata"`
}

type DisplayFeedback struct {
	XMLName     xml.Name `xml:"displayfeedback"`
	FeedbackType string  `xml:"feedbacktype,attr,omitempty"`
	LinkRefId   string   `xml:"linkrefid,attr"`
}

type ItemBody struct {
	XMLName     xml.Name     `xml:"itemBody"`
	P           []P          `xml:"p,omitempty"`
	Div         []Div        `xml:"div,omitempty"`
	ChoiceInteraction []ChoiceInteraction `xml:"choiceInteraction,omitempty"`
	TextEntryInteraction []TextEntryInteraction `xml:"textEntryInteraction,omitempty"`
	ExtendedTextInteraction []ExtendedTextInteraction `xml:"extendedTextInteraction,omitempty"`
}

type P struct {
	XMLName xml.Name `xml:"p"`
	Content string   `xml:",innerxml"`
}

type Div struct {
	XMLName xml.Name `xml:"div"`
	Class   string   `xml:"class,attr,omitempty"`
	Content string   `xml:",innerxml"`
}

type ChoiceInteraction struct {
	XMLName         xml.Name        `xml:"choiceInteraction"`
	ResponseIdent   string          `xml:"responseIdentifier,attr"`
	Shuffle         bool            `xml:"shuffle,attr,omitempty"`
	MaxChoices      int             `xml:"maxChoices,attr,omitempty"`
	MinChoices      int             `xml:"minChoices,attr,omitempty"`
	Prompt          *Prompt         `xml:"prompt,omitempty"`
	SimpleChoice    []SimpleChoice  `xml:"simpleChoice"`
}

type SimpleChoice struct {
	XMLName     xml.Name `xml:"simpleChoice"`
	Identifier  string   `xml:"identifier,attr"`
	Fixed       bool     `xml:"fixed,attr,omitempty"`
	Content     string   `xml:",innerxml"`
}

type Prompt struct {
	XMLName xml.Name `xml:"prompt"`
	Content string   `xml:",innerxml"`
}

type TextEntryInteraction struct {
	XMLName         xml.Name `xml:"textEntryInteraction"`
	ResponseIdent   string   `xml:"responseIdentifier,attr"`
	ExpectedLength  int      `xml:"expectedLength,attr,omitempty"`
	PatternMask     string   `xml:"patternMask,attr,omitempty"`
	PlaceholderText string   `xml:"placeholderText,attr,omitempty"`
}

type ExtendedTextInteraction struct {
	XMLName         xml.Name `xml:"extendedTextInteraction"`
	ResponseIdent   string   `xml:"responseIdentifier,attr"`
	MinStrings      int      `xml:"minStrings,attr,omitempty"`
	MaxStrings      int      `xml:"maxStrings,attr,omitempty"`
	ExpectedLines   int      `xml:"expectedLines,attr,omitempty"`
	ExpectedLength  int      `xml:"expectedLength,attr,omitempty"`
	Prompt          *Prompt  `xml:"prompt,omitempty"`
}

type ResponseDecl struct {
	XMLName       xml.Name      `xml:"responseDeclaration"`
	Identifier    string        `xml:"identifier,attr"`
	Cardinality   string        `xml:"cardinality,attr"`
	BaseType      string        `xml:"baseType,attr,omitempty"`
	CorrectResponse *CorrectResponse `xml:"correctResponse,omitempty"`
	Mapping       *Mapping      `xml:"mapping,omitempty"`
}

type CorrectResponse struct {
	XMLName xml.Name `xml:"correctResponse"`
	Value   []string `xml:"value"`
}

type Mapping struct {
	XMLName        xml.Name    `xml:"mapping"`
	LowerBound     float64     `xml:"lowerBound,attr,omitempty"`
	UpperBound     float64     `xml:"upperBound,attr,omitempty"`
	DefaultValue   float64     `xml:"defaultValue,attr,omitempty"`
	MapEntry       []MapEntry  `xml:"mapEntry"`
}

type MapEntry struct {
	XMLName    xml.Name `xml:"mapEntry"`
	MapKey     string   `xml:"mapKey,attr"`
	MappedValue float64 `xml:"mappedValue,attr"`
}

type OutcomeDecl struct {
	XMLName         xml.Name        `xml:"outcomeDeclaration"`
	Identifier      string          `xml:"identifier,attr"`
	Cardinality     string          `xml:"cardinality,attr"`
	BaseType        string          `xml:"baseType,attr,omitempty"`
	DefaultValue    *DefaultValue   `xml:"defaultValue,omitempty"`
}

type DefaultValue struct {
	XMLName xml.Name `xml:"defaultValue"`
	Value   string   `xml:"value"`
}

type TemplateDecl struct {
	XMLName      xml.Name     `xml:"templateDeclaration"`
	Identifier   string       `xml:"identifier,attr"`
	Cardinality  string       `xml:"cardinality,attr"`
	BaseType     string       `xml:"baseType,attr,omitempty"`
	ParamVariable bool        `xml:"paramVariable,attr,omitempty"`
	DefaultValue *DefaultValue `xml:"defaultValue,omitempty"`
}

type Feedback struct {
	XMLName        xml.Name       `xml:"itemfeedback"`
	Ident          string         `xml:"ident,attr"`
	Title          string         `xml:"title,attr,omitempty"`
	FlowMat        []FlowMat      `xml:"flow_mat,omitempty"`
	Material       *Material      `xml:"material,omitempty"`
}

type FlowMat struct {
	XMLName  xml.Name  `xml:"flow_mat"`
	Material *Material `xml:"material,omitempty"`
}

type Objective struct {
	XMLName  xml.Name  `xml:"objective"`
	Title    string    `xml:"title,attr,omitempty"`
	Material *Material `xml:"material,omitempty"`
}

type RubricBlock struct {
	XMLName xml.Name  `xml:"rubricBlock"`
	Use     string    `xml:"use,attr,omitempty"`
	View    string    `xml:"view,attr,omitempty"`
	Content string    `xml:",innerxml"`
}