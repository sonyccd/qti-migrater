package qti21to30

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/qti-migrator/pkg/models"
)

type Migrator21to30 struct{}

func New() *Migrator21to30 {
	return &Migrator21to30{}
}

// QTI 3.0 specific structs to handle XML naming properly
type QTI3ItemBody struct {
	XMLName                 xml.Name                      `xml:"qti-item-body"`
	P                       []QTI3P                       `xml:"p,omitempty"`
	Div                     []QTI3Div                     `xml:"div,omitempty"`
	ChoiceInteraction       []QTI3ChoiceInteraction       `xml:"qti-choice-interaction,omitempty"`
	TextEntryInteraction    []QTI3TextEntryInteraction    `xml:"qti-text-entry-interaction,omitempty"`
	ExtendedTextInteraction []QTI3ExtendedTextInteraction `xml:"qti-extended-text-interaction,omitempty"`
}

type QTI3P struct {
	XMLName xml.Name `xml:"p"`
	Content string   `xml:",innerxml"`
}

type QTI3Div struct {
	XMLName xml.Name `xml:"div"`
	Class   string   `xml:"data-qti-class,attr,omitempty"`
	Content string   `xml:",innerxml"`
}

type QTI3ChoiceInteraction struct {
	XMLName       xml.Name           `xml:"qti-choice-interaction"`
	ResponseIdent string             `xml:"response-identifier,attr"`
	Shuffle       bool               `xml:"shuffle,attr,omitempty"`
	MaxChoices    int                `xml:"max-choices,attr,omitempty"`
	MinChoices    int                `xml:"min-choices,attr,omitempty"`
	Prompt        *QTI3Prompt        `xml:"qti-prompt,omitempty"`
	SimpleChoice  []QTI3SimpleChoice `xml:"qti-simple-choice"`
}

type QTI3SimpleChoice struct {
	XMLName    xml.Name `xml:"qti-simple-choice"`
	Identifier string   `xml:"identifier,attr"`
	Fixed      bool     `xml:"fixed,attr,omitempty"`
	Content    string   `xml:",innerxml"`
}

type QTI3Prompt struct {
	XMLName xml.Name `xml:"qti-prompt"`
	Content string   `xml:",innerxml"`
}

type QTI3TextEntryInteraction struct {
	XMLName         xml.Name `xml:"qti-text-entry-interaction"`
	ResponseIdent   string   `xml:"response-identifier,attr"`
	ExpectedLength  int      `xml:"expected-length,attr,omitempty"`
	PatternMask     string   `xml:"pattern-mask,attr,omitempty"`
	PlaceholderText string   `xml:"placeholder-text,attr,omitempty"`
}

type QTI3ExtendedTextInteraction struct {
	XMLName        xml.Name    `xml:"qti-extended-text-interaction"`
	ResponseIdent  string      `xml:"response-identifier,attr"`
	MinStrings     int         `xml:"min-strings,attr,omitempty"`
	MaxStrings     int         `xml:"max-strings,attr,omitempty"`
	ExpectedLines  int         `xml:"expected-lines,attr,omitempty"`
	ExpectedLength int         `xml:"expected-length,attr,omitempty"`
	Prompt         *QTI3Prompt `xml:"qti-prompt,omitempty"`
}

type QTI3ResponseDecl struct {
	XMLName         xml.Name              `xml:"qti-response-declaration"`
	Identifier      string                `xml:"identifier,attr"`
	Cardinality     string                `xml:"cardinality,attr"`
	BaseType        string                `xml:"base-type,attr,omitempty"`
	CorrectResponse *QTI3CorrectResponse  `xml:"qti-correct-response,omitempty"`
	Mapping         *QTI3Mapping          `xml:"qti-mapping,omitempty"`
}

type QTI3CorrectResponse struct {
	XMLName xml.Name  `xml:"qti-correct-response"`
	Value   []QTI3Value `xml:"qti-value"`
}

type QTI3Value struct {
	XMLName xml.Name `xml:"qti-value"`
	Content string   `xml:",chardata"`
}

type QTI3Mapping struct {
	XMLName      xml.Name       `xml:"qti-mapping"`
	LowerBound   float64        `xml:"lower-bound,attr,omitempty"`
	UpperBound   float64        `xml:"upper-bound,attr,omitempty"`
	DefaultValue float64        `xml:"default-value,attr,omitempty"`
	MapEntry     []QTI3MapEntry `xml:"qti-map-entry"`
}

type QTI3MapEntry struct {
	XMLName     xml.Name `xml:"qti-map-entry"`
	MapKey      string   `xml:"map-key,attr"`
	MappedValue float64  `xml:"mapped-value,attr"`
}

type QTI3OutcomeDecl struct {
	XMLName      xml.Name          `xml:"qti-outcome-declaration"`
	Identifier   string            `xml:"identifier,attr"`
	Cardinality  string            `xml:"cardinality,attr"`
	BaseType     string            `xml:"base-type,attr,omitempty"`
	DefaultValue *QTI3DefaultValue `xml:"qti-default-value,omitempty"`
}

type QTI3DefaultValue struct {
	XMLName xml.Name  `xml:"qti-default-value"`
	Value   []QTI3Value `xml:"qti-value"`
}

type QTI3Feedback struct {
	XMLName xml.Name `xml:"qti-modal-feedback"`
	Ident   string   `xml:"identifier,attr"`
	Title   string   `xml:"title,attr,omitempty"`
	Content string   `xml:",innerxml"`
}

// QTI3Item represents a QTI 3.0 assessment item
type QTI3Item struct {
	XMLName         xml.Name                      `xml:"http://www.imsglobal.org/xsd/imsqtiasi_v3p0 qti-assessment-item"`
	Identifier      string                        `xml:"identifier,attr"`
	Title           string                        `xml:"title,attr"`
	Adaptive        string                        `xml:"adaptive,attr,omitempty"`
	TimeDependent   string                        `xml:"time-dependent,attr,omitempty"`
	ResponseDecl    []QTI3ResponseDecl            `xml:"qti-response-declaration,omitempty"`
	OutcomeDecl     []QTI3OutcomeDecl             `xml:"qti-outcome-declaration,omitempty"`
	ItemBody        *QTI3ItemBody                 `xml:"qti-item-body,omitempty"`
	Feedback        []QTI3Feedback                `xml:"qti-modal-feedback,omitempty"`
}

func (m *Migrator21to30) Migrate(doc interface{}) ([]byte, error) {
	qtiDoc, ok := doc.(*models.QTIDocument)
	if !ok {
		return nil, fmt.Errorf("invalid document type for QTI 2.1 to 3.0 migration")
	}

	// For simplicity, we'll handle single item documents with the proper QTI 3.0 structure
	if len(qtiDoc.Items) == 1 && qtiDoc.Assessment == nil {
		return m.migrateSingleItem(&qtiDoc.Items[0])
	}

	// For documents with assessments or multiple items, create a proper QTI 3.0 structure
	// But for now, just migrate using the document structure and apply QTI 3.0 element names
	migratedDoc := m.migrateDocument(qtiDoc)

	output, err := xml.MarshalIndent(migratedDoc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal migrated document: %w", err)
	}

	xmlHeader := []byte(xml.Header)
	return append(xmlHeader, output...), nil
}

func (m *Migrator21to30) migrateSingleItem(item *models.Item) ([]byte, error) {
	qti3Item := QTI3Item{
		Identifier:    item.Ident,
		Title:         item.Title,
		Adaptive:      "false",
		TimeDependent: "false",
	}

	// Migrate response declarations
	for _, decl := range item.ResponseDecl {
		qti3Item.ResponseDecl = append(qti3Item.ResponseDecl, m.migrateResponseDeclarationToQTI3(&decl))
	}

	// Migrate outcome declarations
	for _, decl := range item.OutcomeDecl {
		qti3Item.OutcomeDecl = append(qti3Item.OutcomeDecl, m.migrateOutcomeDeclarationToQTI3(&decl))
	}

	// Migrate item body
	if item.ItemBody != nil {
		qti3Item.ItemBody = m.migrateItemBodyToQTI3(item.ItemBody)
	}

	// Migrate feedback
	for _, feedback := range item.Feedback {
		qti3Item.Feedback = append(qti3Item.Feedback, m.migrateFeedbackToQTI3(&feedback))
	}

	output, err := xml.MarshalIndent(qti3Item, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal migrated item: %w", err)
	}

	xmlHeader := []byte(xml.Header)
	return append(xmlHeader, output...), nil
}

func (m *Migrator21to30) migrateDocument(doc *models.QTIDocument) *models.QTIDocument {
	// For QTI 3.0, the root element changes based on content
	rootName := "questestinterop"
	if doc.Assessment != nil {
		// Assessment document becomes qti-assessment-test
		rootName = "qti-assessment-test"
	}
	
	migratedDoc := &models.QTIDocument{
		XMLName: xml.Name{
			Space: "http://www.imsglobal.org/xsd/imsqtiasi_v3p0",
			Local: rootName,
		},
		Version: "3.0",
	}

	for _, item := range doc.Items {
		migratedItem := m.migrateItem(&item)
		migratedDoc.Items = append(migratedDoc.Items, *migratedItem)
	}

	if doc.Assessment != nil {
		migratedDoc.Assessment = m.migrateAssessment(doc.Assessment)
	}

	if doc.Metadata != nil {
		migratedDoc.Metadata = m.migrateMetadata(doc.Metadata)
	}

	return migratedDoc
}

func (m *Migrator21to30) migrateAssessment(assessment *models.Assessment) *models.Assessment {
	migratedAssessment := &models.Assessment{
		XMLName:     xml.Name{Local: "qti-assessment-test"},
		Title:       assessment.Title,
		Ident:       assessment.Ident,
		Objectives:  assessment.Objectives,
		RubricBlock: m.migrateRubricBlock(assessment.RubricBlock),
	}

	if assessment.Metadata != nil {
		migratedAssessment.Metadata = m.migrateMetadata(assessment.Metadata)
	}

	for _, section := range assessment.Sections {
		migratedSection := m.migrateSection(&section)
		migratedAssessment.Sections = append(migratedAssessment.Sections, *migratedSection)
	}

	return migratedAssessment
}

func (m *Migrator21to30) migrateSection(section *models.Section) *models.Section {
	migratedSection := &models.Section{
		XMLName: xml.Name{Local: "qti-assessment-section"},
		Title:   section.Title,
		Ident:   section.Ident,
	}

	if section.Metadata != nil {
		migratedSection.Metadata = m.migrateMetadata(section.Metadata)
	}

	for _, item := range section.Items {
		migratedItem := m.migrateItem(&item)
		migratedSection.Items = append(migratedSection.Items, *migratedItem)
	}

	return migratedSection
}

func (m *Migrator21to30) migrateItem(item *models.Item) *models.Item {
	migratedItem := &models.Item{
		XMLName:     xml.Name{Local: "qti-assessment-item"},
		Title:       item.Title,
		Ident:       item.Ident,
		MaxAttempts: item.MaxAttempts,
		RubricBlock: m.migrateRubricBlock(item.RubricBlock),
	}

	if item.Metadata != nil {
		migratedItem.Metadata = m.migrateMetadata(item.Metadata)
	}

	if item.ItemBody != nil {
		migratedItem.ItemBody = m.migrateItemBody(item.ItemBody)
	}

	for _, decl := range item.ResponseDecl {
		migratedItem.ResponseDecl = append(migratedItem.ResponseDecl, m.migrateResponseDeclaration(&decl))
	}

	for _, decl := range item.OutcomeDecl {
		migratedItem.OutcomeDecl = append(migratedItem.OutcomeDecl, m.migrateOutcomeDeclaration(&decl))
	}

	for _, decl := range item.TemplateDecl {
		migratedItem.TemplateDecl = append(migratedItem.TemplateDecl, m.migrateTemplateDeclaration(&decl))
	}

	for _, feedback := range item.Feedback {
		migratedItem.Feedback = append(migratedItem.Feedback, m.migrateFeedback(&feedback))
	}

	return migratedItem
}

func (m *Migrator21to30) migrateMetadata(metadata *models.Metadata) *models.Metadata {
	return &models.Metadata{
		XMLName:     xml.Name{Local: "qti-metadata"},
		Schema:      metadata.Schema,
		SchemaVer:   "3.0",
		LOM:         metadata.LOM,
		QTIMetadata: m.migrateQTIMetadata(metadata.QTIMetadata),
	}
}

func (m *Migrator21to30) migrateQTIMetadata(qtiMetadata *models.QTIMetadata) *models.QTIMetadata {
	if qtiMetadata == nil {
		return nil
	}

	return &models.QTIMetadata{
		XMLName:            xml.Name{Local: "qti-metadata-container"},
		TimeDependent:      qtiMetadata.TimeDependent,
		Composite:          qtiMetadata.Composite,
		InteractionType:    m.migrateInteractionType(qtiMetadata.InteractionType),
		FeedbackType:       qtiMetadata.FeedbackType,
		SolutionAvailable:  qtiMetadata.SolutionAvailable,
		Scoringmode:        qtiMetadata.Scoringmode,
		ToolName:           qtiMetadata.ToolName,
		ToolVersion:        qtiMetadata.ToolVersion,
		ToolVendor:         qtiMetadata.ToolVendor,
	}
}

func (m *Migrator21to30) migrateInteractionType(interactionType string) string {
	switch interactionType {
	case "choiceInteraction":
		return "qti-choice-interaction"
	case "textEntryInteraction":
		return "qti-text-entry-interaction"
	case "extendedTextInteraction":
		return "qti-extended-text-interaction"
	case "matchInteraction":
		return "qti-match-interaction"
	case "associateInteraction":
		return "qti-associate-interaction"
	case "orderInteraction":
		return "qti-order-interaction"
	case "hotspotInteraction":
		return "qti-hotspot-interaction"
	case "selectPointInteraction":
		return "qti-select-point-interaction"
	case "graphicAssociateInteraction":
		return "qti-graphic-associate-interaction"
	case "graphicOrderInteraction":
		return "qti-graphic-order-interaction"
	case "graphicGapMatchInteraction":
		return "qti-graphic-gap-match-interaction"
	case "positionObjectInteraction":
		return "qti-position-object-interaction"
	case "sliderInteraction":
		return "qti-slider-interaction"
	case "drawingInteraction":
		return "qti-drawing-interaction"
	case "gapMatchInteraction":
		return "qti-gap-match-interaction"
	case "inlineChoiceInteraction":
		return "qti-inline-choice-interaction"
	case "hottextInteraction":
		return "qti-hottext-interaction"
	case "uploadInteraction":
		return "qti-upload-interaction"
	default:
		return interactionType
	}
}

func (m *Migrator21to30) migrateItemBody(itemBody *models.ItemBody) *models.ItemBody {
	migratedItemBody := &models.ItemBody{
		XMLName: xml.Name{Local: "qti-item-body"},
	}

	for _, p := range itemBody.P {
		migratedItemBody.P = append(migratedItemBody.P, models.P{
			XMLName: xml.Name{Local: "p"},
			Content: m.updateHTMLContent(p.Content),
		})
	}

	for _, div := range itemBody.Div {
		migratedItemBody.Div = append(migratedItemBody.Div, models.Div{
			XMLName: xml.Name{Local: "div"},
			Class:   div.Class,
			Content: m.updateHTMLContent(div.Content),
		})
	}

	for _, interaction := range itemBody.ChoiceInteraction {
		migratedItemBody.ChoiceInteraction = append(migratedItemBody.ChoiceInteraction, m.migrateChoiceInteraction(&interaction))
	}

	for _, interaction := range itemBody.TextEntryInteraction {
		migratedItemBody.TextEntryInteraction = append(migratedItemBody.TextEntryInteraction, m.migrateTextEntryInteraction(&interaction))
	}

	for _, interaction := range itemBody.ExtendedTextInteraction {
		migratedItemBody.ExtendedTextInteraction = append(migratedItemBody.ExtendedTextInteraction, m.migrateExtendedTextInteraction(&interaction))
	}

	return migratedItemBody
}

func (m *Migrator21to30) migrateChoiceInteraction(interaction *models.ChoiceInteraction) models.ChoiceInteraction {
	migratedInteraction := models.ChoiceInteraction{
		XMLName:       xml.Name{Local: "qti-choice-interaction"},
		ResponseIdent: interaction.ResponseIdent,
		Shuffle:       interaction.Shuffle,
		MaxChoices:    interaction.MaxChoices,
		MinChoices:    interaction.MinChoices,
	}

	if interaction.Prompt != nil {
		migratedInteraction.Prompt = &models.Prompt{
			XMLName: xml.Name{Local: "qti-prompt"},
			Content: m.updateHTMLContent(interaction.Prompt.Content),
		}
	}

	for _, choice := range interaction.SimpleChoice {
		migratedInteraction.SimpleChoice = append(migratedInteraction.SimpleChoice, models.SimpleChoice{
			XMLName:    xml.Name{Local: "qti-simple-choice"},
			Identifier: choice.Identifier,
			Fixed:      choice.Fixed,
			Content:    m.updateHTMLContent(choice.Content),
		})
	}

	return migratedInteraction
}

func (m *Migrator21to30) migrateTextEntryInteraction(interaction *models.TextEntryInteraction) models.TextEntryInteraction {
	return models.TextEntryInteraction{
		XMLName:         xml.Name{Local: "qti-text-entry-interaction"},
		ResponseIdent:   interaction.ResponseIdent,
		ExpectedLength:  interaction.ExpectedLength,
		PatternMask:     interaction.PatternMask,
		PlaceholderText: interaction.PlaceholderText,
	}
}

func (m *Migrator21to30) migrateExtendedTextInteraction(interaction *models.ExtendedTextInteraction) models.ExtendedTextInteraction {
	migratedInteraction := models.ExtendedTextInteraction{
		XMLName:        xml.Name{Local: "qti-extended-text-interaction"},
		ResponseIdent:  interaction.ResponseIdent,
		MinStrings:     interaction.MinStrings,
		MaxStrings:     interaction.MaxStrings,
		ExpectedLines:  interaction.ExpectedLines,
		ExpectedLength: interaction.ExpectedLength,
	}

	if interaction.Prompt != nil {
		migratedInteraction.Prompt = &models.Prompt{
			XMLName: xml.Name{Local: "qti-prompt"},
			Content: m.updateHTMLContent(interaction.Prompt.Content),
		}
	}

	return migratedInteraction
}

func (m *Migrator21to30) migrateResponseDeclaration(decl *models.ResponseDecl) models.ResponseDecl {
	migratedDecl := models.ResponseDecl{
		XMLName:     xml.Name{Local: "qti-response-declaration"},
		Identifier:  decl.Identifier,
		Cardinality: decl.Cardinality,
		BaseType:    m.migrateBaseType(decl.BaseType),
	}

	if decl.CorrectResponse != nil {
		migratedDecl.CorrectResponse = &models.CorrectResponse{
			XMLName: xml.Name{Local: "qti-correct-response"},
			Value:   decl.CorrectResponse.Value,
		}
	}

	if decl.Mapping != nil {
		migratedDecl.Mapping = m.migrateMapping(decl.Mapping)
	}

	return migratedDecl
}

func (m *Migrator21to30) migrateBaseType(baseType string) string {
	switch baseType {
	case "string":
		return "string"
	case "integer":
		return "integer"
	case "float":
		return "float"
	case "boolean":
		return "boolean"
	case "identifier":
		return "identifier"
	case "point":
		return "point"
	case "pair":
		return "directedPair"
	case "duration":
		return "duration"
	case "file":
		return "uri"
	default:
		return baseType
	}
}

func (m *Migrator21to30) migrateMapping(mapping *models.Mapping) *models.Mapping {
	migratedMapping := &models.Mapping{
		XMLName:      xml.Name{Local: "qti-mapping"},
		LowerBound:   mapping.LowerBound,
		UpperBound:   mapping.UpperBound,
		DefaultValue: mapping.DefaultValue,
	}

	for _, entry := range mapping.MapEntry {
		migratedMapping.MapEntry = append(migratedMapping.MapEntry, models.MapEntry{
			XMLName:     xml.Name{Local: "qti-map-entry"},
			MapKey:      entry.MapKey,
			MappedValue: entry.MappedValue,
		})
	}

	return migratedMapping
}

func (m *Migrator21to30) migrateOutcomeDeclaration(decl *models.OutcomeDecl) models.OutcomeDecl {
	migratedDecl := models.OutcomeDecl{
		XMLName:     xml.Name{Local: "qti-outcome-declaration"},
		Identifier:  decl.Identifier,
		Cardinality: decl.Cardinality,
		BaseType:    m.migrateBaseType(decl.BaseType),
	}

	if decl.DefaultValue != nil {
		migratedDecl.DefaultValue = &models.DefaultValue{
			XMLName: xml.Name{Local: "qti-default-value"},
			Value:   decl.DefaultValue.Value,
		}
	}

	return migratedDecl
}

func (m *Migrator21to30) migrateTemplateDeclaration(decl *models.TemplateDecl) models.TemplateDecl {
	migratedDecl := models.TemplateDecl{
		XMLName:       xml.Name{Local: "qti-template-declaration"},
		Identifier:    decl.Identifier,
		Cardinality:   decl.Cardinality,
		BaseType:      m.migrateBaseType(decl.BaseType),
		ParamVariable: decl.ParamVariable,
	}

	if decl.DefaultValue != nil {
		migratedDecl.DefaultValue = &models.DefaultValue{
			XMLName: xml.Name{Local: "qti-default-value"},
			Value:   decl.DefaultValue.Value,
		}
	}

	return migratedDecl
}

func (m *Migrator21to30) migrateFeedback(feedback *models.Feedback) models.Feedback {
	migratedFeedback := models.Feedback{
		XMLName: xml.Name{Local: "qti-modal-feedback"},
		Ident:   feedback.Ident,
		Title:   feedback.Title,
	}

	if feedback.Material != nil {
		migratedFeedback.Material = m.migrateMaterial(feedback.Material)
	}

	for _, flowMat := range feedback.FlowMat {
		if flowMat.Material != nil {
			flowMat.Material = m.migrateMaterial(flowMat.Material)
		}
		migratedFeedback.FlowMat = append(migratedFeedback.FlowMat, flowMat)
	}

	return migratedFeedback
}

func (m *Migrator21to30) migrateMaterial(material *models.Material) *models.Material {
	if material == nil {
		return nil
	}

	migratedMaterial := &models.Material{
		XMLName: material.XMLName,
		Label:   material.Label,
	}

	for _, matText := range material.MatText {
		migratedMaterial.MatText = append(migratedMaterial.MatText, models.MatText{
			XMLName:  matText.XMLName,
			TextType: matText.TextType,
			Charset:  matText.Charset,
			XML:      matText.XML,
			Content:  m.updateHTMLContent(matText.Content),
		})
	}

	migratedMaterial.MatImage = material.MatImage
	migratedMaterial.MatAudio = material.MatAudio
	migratedMaterial.MatVideo = material.MatVideo

	return migratedMaterial
}

func (m *Migrator21to30) migrateRubricBlock(rubricBlock *models.RubricBlock) *models.RubricBlock {
	if rubricBlock == nil {
		return nil
	}

	return &models.RubricBlock{
		XMLName: xml.Name{Local: "qti-rubric-block"},
		Use:     rubricBlock.Use,
		View:    m.migrateView(rubricBlock.View),
		Content: m.updateHTMLContent(rubricBlock.Content),
	}
}

func (m *Migrator21to30) migrateView(view string) string {
	switch view {
	case "author":
		return "author"
	case "candidate":
		return "candidate"
	case "proctor":
		return "proctor"
	case "scorer":
		return "scorer"
	case "testConstructor":
		return "test-constructor"
	case "tutor":
		return "tutor"
	default:
		return view
	}
}

func (m *Migrator21to30) updateHTMLContent(content string) string {
	content = strings.ReplaceAll(content, "class=", "data-qti-class=")
	
	content = strings.ReplaceAll(content, "<object", "<qti-object")
	content = strings.ReplaceAll(content, "</object>", "</qti-object>")
	
	return content
}

// Migration functions for QTI3 types
func (m *Migrator21to30) migrateItemBodyToQTI3(itemBody *models.ItemBody) *QTI3ItemBody {
	qti3ItemBody := &QTI3ItemBody{}

	for _, p := range itemBody.P {
		qti3ItemBody.P = append(qti3ItemBody.P, QTI3P{
			Content: m.updateHTMLContent(p.Content),
		})
	}

	for _, div := range itemBody.Div {
		qti3ItemBody.Div = append(qti3ItemBody.Div, QTI3Div{
			Class:   div.Class,
			Content: m.updateHTMLContent(div.Content),
		})
	}

	for _, interaction := range itemBody.ChoiceInteraction {
		qti3ItemBody.ChoiceInteraction = append(qti3ItemBody.ChoiceInteraction, m.migrateChoiceInteractionToQTI3(&interaction))
	}

	for _, interaction := range itemBody.TextEntryInteraction {
		qti3ItemBody.TextEntryInteraction = append(qti3ItemBody.TextEntryInteraction, m.migrateTextEntryInteractionToQTI3(&interaction))
	}

	for _, interaction := range itemBody.ExtendedTextInteraction {
		qti3ItemBody.ExtendedTextInteraction = append(qti3ItemBody.ExtendedTextInteraction, m.migrateExtendedTextInteractionToQTI3(&interaction))
	}

	return qti3ItemBody
}

func (m *Migrator21to30) migrateChoiceInteractionToQTI3(interaction *models.ChoiceInteraction) QTI3ChoiceInteraction {
	qti3Interaction := QTI3ChoiceInteraction{
		ResponseIdent: interaction.ResponseIdent,
		Shuffle:       interaction.Shuffle,
		MaxChoices:    interaction.MaxChoices,
		MinChoices:    interaction.MinChoices,
	}

	if interaction.Prompt != nil {
		qti3Interaction.Prompt = &QTI3Prompt{
			Content: m.updateHTMLContent(interaction.Prompt.Content),
		}
	}

	for _, choice := range interaction.SimpleChoice {
		qti3Interaction.SimpleChoice = append(qti3Interaction.SimpleChoice, QTI3SimpleChoice{
			Identifier: choice.Identifier,
			Fixed:      choice.Fixed,
			Content:    m.updateHTMLContent(choice.Content),
		})
	}

	return qti3Interaction
}

func (m *Migrator21to30) migrateTextEntryInteractionToQTI3(interaction *models.TextEntryInteraction) QTI3TextEntryInteraction {
	return QTI3TextEntryInteraction{
		ResponseIdent:   interaction.ResponseIdent,
		ExpectedLength:  interaction.ExpectedLength,
		PatternMask:     interaction.PatternMask,
		PlaceholderText: interaction.PlaceholderText,
	}
}

func (m *Migrator21to30) migrateExtendedTextInteractionToQTI3(interaction *models.ExtendedTextInteraction) QTI3ExtendedTextInteraction {
	qti3Interaction := QTI3ExtendedTextInteraction{
		ResponseIdent:  interaction.ResponseIdent,
		MinStrings:     interaction.MinStrings,
		MaxStrings:     interaction.MaxStrings,
		ExpectedLines:  interaction.ExpectedLines,
		ExpectedLength: interaction.ExpectedLength,
	}

	if interaction.Prompt != nil {
		qti3Interaction.Prompt = &QTI3Prompt{
			Content: m.updateHTMLContent(interaction.Prompt.Content),
		}
	}

	return qti3Interaction
}

func (m *Migrator21to30) migrateResponseDeclarationToQTI3(decl *models.ResponseDecl) QTI3ResponseDecl {
	qti3Decl := QTI3ResponseDecl{
		Identifier:  decl.Identifier,
		Cardinality: decl.Cardinality,
		BaseType:    m.migrateBaseType(decl.BaseType),
	}

	if decl.CorrectResponse != nil {
		qti3Decl.CorrectResponse = &QTI3CorrectResponse{}
		for _, value := range decl.CorrectResponse.Value {
			qti3Decl.CorrectResponse.Value = append(qti3Decl.CorrectResponse.Value, QTI3Value{
				Content: value,
			})
		}
	}

	if decl.Mapping != nil {
		qti3Decl.Mapping = &QTI3Mapping{
			LowerBound:   decl.Mapping.LowerBound,
			UpperBound:   decl.Mapping.UpperBound,
			DefaultValue: decl.Mapping.DefaultValue,
		}
		for _, entry := range decl.Mapping.MapEntry {
			qti3Decl.Mapping.MapEntry = append(qti3Decl.Mapping.MapEntry, QTI3MapEntry{
				MapKey:      entry.MapKey,
				MappedValue: entry.MappedValue,
			})
		}
	}

	return qti3Decl
}

func (m *Migrator21to30) migrateOutcomeDeclarationToQTI3(decl *models.OutcomeDecl) QTI3OutcomeDecl {
	qti3Decl := QTI3OutcomeDecl{
		Identifier:  decl.Identifier,
		Cardinality: decl.Cardinality,
		BaseType:    m.migrateBaseType(decl.BaseType),
	}

	if decl.DefaultValue != nil {
		qti3Decl.DefaultValue = &QTI3DefaultValue{}
		qti3Decl.DefaultValue.Value = append(qti3Decl.DefaultValue.Value, QTI3Value{
			Content: decl.DefaultValue.Value,
		})
	}

	return qti3Decl
}

func (m *Migrator21to30) migrateFeedbackToQTI3(feedback *models.Feedback) QTI3Feedback {
	qti3Feedback := QTI3Feedback{
		Ident: feedback.Ident,
		Title: feedback.Title,
	}

	// Convert material/flowmat content to simple content
	var content strings.Builder
	if feedback.Material != nil {
		for _, matText := range feedback.Material.MatText {
			content.WriteString(m.updateHTMLContent(matText.Content))
		}
	}
	for _, flowMat := range feedback.FlowMat {
		if flowMat.Material != nil {
			for _, matText := range flowMat.Material.MatText {
				content.WriteString(m.updateHTMLContent(matText.Content))
			}
		}
	}
	qti3Feedback.Content = content.String()

	return qti3Feedback
}