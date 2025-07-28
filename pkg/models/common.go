package models

import (
	"encoding/xml"
	"time"
)

// Common metadata structures used across all QTI versions

type Metadata struct {
	XMLName     xml.Name    `xml:"metadata"`
	Schema      string      `xml:"schema,omitempty"`
	SchemaVer   string      `xml:"schemaversion,omitempty"`
	LOM         *LOM        `xml:"lom,omitempty"`
	QTIMetadata *QTIMetadata `xml:"qtimetadata,omitempty"`
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

// LOM (Learning Object Metadata) structures

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

// Common material structures used across versions

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

// Common structures that appear in multiple versions but might have slight variations

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