package models

import "encoding/xml"

// QTI 1.2 specific structures

type QTIDocument12 struct {
	XMLName    xml.Name    `xml:"questestinterop"`
	Version    string      `xml:"version,attr"`
	Items      []Item12    `xml:"item"`
	Assessment *Assessment12 `xml:"assessment,omitempty"`
	Metadata   *Metadata   `xml:"metadata,omitempty"`
}

type Assessment12 struct {
	XMLName     xml.Name     `xml:"assessment"`
	Title       string       `xml:"title,attr"`
	Ident       string       `xml:"ident,attr"`
	Sections    []Section12  `xml:"section"`
	Metadata    *Metadata    `xml:"metadata,omitempty"`
	Objectives  []Objective  `xml:"objectives>objective,omitempty"`
	RubricBlock *RubricBlock `xml:"rubricBlock,omitempty"`
}

type Section12 struct {
	XMLName  xml.Name  `xml:"section"`
	Title    string    `xml:"title,attr"`
	Ident    string    `xml:"ident,attr"`
	Items    []Item12  `xml:"item"`
	Metadata *Metadata `xml:"metadata,omitempty"`
}

type Item12 struct {
	XMLName        xml.Name        `xml:"item"`
	Title          string          `xml:"title,attr"`
	Ident          string          `xml:"ident,attr"`
	MaxAttempts    int             `xml:"maxattempts,attr,omitempty"`
	Metadata       *Metadata       `xml:"metadata,omitempty"`
	Presentation   *Presentation   `xml:"presentation,omitempty"`
	ResponseProc   *ResponseProc   `xml:"resprocessing,omitempty"`
	Feedback       []Feedback12    `xml:"itemfeedback,omitempty"`
	RubricBlock    *RubricBlock    `xml:"rubricBlock,omitempty"`
}

// QTI 1.2 Presentation structures

type Presentation struct {
	XMLName  xml.Name  `xml:"presentation"`
	Label    string    `xml:"label,attr,omitempty"`
	Material *Material `xml:"material,omitempty"`
	Response []Response12 `xml:"response_lid,omitempty"`
	Flow     []Flow12     `xml:"flow,omitempty"`
}

type Response12 struct {
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

type Flow12 struct {
	XMLName    xml.Name    `xml:"flow"`
	Class      string      `xml:"class,attr,omitempty"`
	Material   []Material  `xml:"material,omitempty"`
	Response   []Response12 `xml:"response_lid,omitempty"`
	Flow       []Flow12     `xml:"flow,omitempty"`
}

// QTI 1.2 Response Processing structures

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

// QTI 1.2 Feedback structures

type Feedback12 struct {
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