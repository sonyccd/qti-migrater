package qti12to21

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/qti-migrator/pkg/models"
)

type Migrator12to21 struct{}

func New() *Migrator12to21 {
	return &Migrator12to21{}
}

func (m *Migrator12to21) Migrate(doc interface{}) ([]byte, error) {
	qtiDoc, ok := doc.(*models.QTIDocument)
	if !ok {
		return nil, fmt.Errorf("invalid document type for QTI 1.2 to 2.1 migration")
	}

	migratedDoc := m.migrateDocument(qtiDoc)

	output, err := xml.MarshalIndent(migratedDoc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal migrated document: %w", err)
	}

	xmlHeader := []byte(xml.Header)
	return append(xmlHeader, output...), nil
}

func (m *Migrator12to21) migrateDocument(doc *models.QTIDocument) *models.QTIDocument {
	migratedDoc := &models.QTIDocument{
		XMLName: doc.XMLName,
		Version: "2.1",
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

func (m *Migrator12to21) migrateAssessment(assessment *models.Assessment) *models.Assessment {
	migratedAssessment := &models.Assessment{
		XMLName:     assessment.XMLName,
		Title:       assessment.Title,
		Ident:       assessment.Ident,
		Objectives:  assessment.Objectives,
		RubricBlock: assessment.RubricBlock,
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

func (m *Migrator12to21) migrateSection(section *models.Section) *models.Section {
	migratedSection := &models.Section{
		XMLName: section.XMLName,
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

func (m *Migrator12to21) migrateItem(item *models.Item) *models.Item {
	migratedItem := &models.Item{
		XMLName:     item.XMLName,
		Title:       item.Title,
		Ident:       item.Ident,
		MaxAttempts: item.MaxAttempts,
		RubricBlock: item.RubricBlock,
	}

	if item.Metadata != nil {
		migratedItem.Metadata = m.migrateMetadata(item.Metadata)
	}

	if item.Presentation != nil {
		migratedItem.ItemBody = m.convertPresentationToItemBody(item.Presentation)
		migratedItem.ResponseDecl = m.extractResponseDeclarations(item.Presentation, item.ResponseProc)
	}

	if item.ResponseProc != nil {
		migratedItem.OutcomeDecl = m.extractOutcomeDeclarations(item.ResponseProc)
	}

	for _, feedback := range item.Feedback {
		migratedItem.Feedback = append(migratedItem.Feedback, m.migrateFeedback(&feedback))
	}

	return migratedItem
}

func (m *Migrator12to21) migrateMetadata(metadata *models.Metadata) *models.Metadata {
	return &models.Metadata{
		XMLName:     metadata.XMLName,
		Schema:      metadata.Schema,
		SchemaVer:   "2.1",
		LOM:         metadata.LOM,
		QTIMetadata: metadata.QTIMetadata,
	}
}

func (m *Migrator12to21) convertPresentationToItemBody(presentation *models.Presentation) *models.ItemBody {
	itemBody := &models.ItemBody{
		XMLName: xml.Name{Local: "itemBody"},
	}

	if presentation.Material != nil {
		itemBody.P = m.convertMaterialToParagraphs(presentation.Material)
	}

	for _, response := range presentation.Response {
		if response.RenderChoice != nil {
			choiceInteraction := m.convertResponseToChoiceInteraction(&response)
			itemBody.ChoiceInteraction = append(itemBody.ChoiceInteraction, *choiceInteraction)
		} else if response.RenderFib != nil {
			if response.RenderFib.Rows > 1 {
				extTextInteraction := m.convertResponseToExtendedTextInteraction(&response)
				itemBody.ExtendedTextInteraction = append(itemBody.ExtendedTextInteraction, *extTextInteraction)
			} else {
				textEntryInteraction := m.convertResponseToTextEntryInteraction(&response)
				itemBody.TextEntryInteraction = append(itemBody.TextEntryInteraction, *textEntryInteraction)
			}
		}
	}

	for _, flow := range presentation.Flow {
		m.processFlow(&flow, itemBody)
	}

	return itemBody
}

func (m *Migrator12to21) processFlow(flow *models.Flow, itemBody *models.ItemBody) {
	for _, material := range flow.Material {
		paragraphs := m.convertMaterialToParagraphs(&material)
		itemBody.P = append(itemBody.P, paragraphs...)
	}

	for _, response := range flow.Response {
		if response.RenderChoice != nil {
			choiceInteraction := m.convertResponseToChoiceInteraction(&response)
			itemBody.ChoiceInteraction = append(itemBody.ChoiceInteraction, *choiceInteraction)
		} else if response.RenderFib != nil {
			if response.RenderFib.Rows > 1 {
				extTextInteraction := m.convertResponseToExtendedTextInteraction(&response)
				itemBody.ExtendedTextInteraction = append(itemBody.ExtendedTextInteraction, *extTextInteraction)
			} else {
				textEntryInteraction := m.convertResponseToTextEntryInteraction(&response)
				itemBody.TextEntryInteraction = append(itemBody.TextEntryInteraction, *textEntryInteraction)
			}
		}
	}

	for _, subFlow := range flow.Flow {
		m.processFlow(&subFlow, itemBody)
	}
}

func (m *Migrator12to21) convertMaterialToParagraphs(material *models.Material) []models.P {
	var paragraphs []models.P

	for _, matText := range material.MatText {
		content := matText.Content
		if matText.TextType == "text/html" {
			content = m.sanitizeHTMLContent(content)
		}
		paragraphs = append(paragraphs, models.P{
			XMLName: xml.Name{Local: "p"},
			Content: content,
		})
	}

	for _, matImage := range material.MatImage {
		imgTag := fmt.Sprintf(`<img src="%s"`, matImage.URI)
		if matImage.Width > 0 {
			imgTag += fmt.Sprintf(` width="%d"`, matImage.Width)
		}
		if matImage.Height > 0 {
			imgTag += fmt.Sprintf(` height="%d"`, matImage.Height)
		}
		imgTag += " />"
		
		paragraphs = append(paragraphs, models.P{
			XMLName: xml.Name{Local: "p"},
			Content: imgTag,
		})
	}

	return paragraphs
}

func (m *Migrator12to21) convertResponseToChoiceInteraction(response *models.Response) *models.ChoiceInteraction {
	choiceInteraction := &models.ChoiceInteraction{
		XMLName:       xml.Name{Local: "choiceInteraction"},
		ResponseIdent: response.Ident,
	}

	if response.RenderChoice != nil {
		choiceInteraction.Shuffle = response.RenderChoice.Shuffle == "yes"
		
		if response.RenderChoice.MaxNumber > 0 {
			choiceInteraction.MaxChoices = response.RenderChoice.MaxNumber
		}
		if response.RenderChoice.MInNumber > 0 {
			choiceInteraction.MinChoices = response.RenderChoice.MInNumber
		}

		for _, label := range response.RenderChoice.ResponseLabel {
			simpleChoice := models.SimpleChoice{
				XMLName:    xml.Name{Local: "simpleChoice"},
				Identifier: label.Ident,
			}

			if label.Material != nil {
				content := m.extractMaterialContent(label.Material)
				simpleChoice.Content = content
			}

			choiceInteraction.SimpleChoice = append(choiceInteraction.SimpleChoice, simpleChoice)
		}
	}

	return choiceInteraction
}

func (m *Migrator12to21) convertResponseToTextEntryInteraction(response *models.Response) *models.TextEntryInteraction {
	textEntry := &models.TextEntryInteraction{
		XMLName:       xml.Name{Local: "textEntryInteraction"},
		ResponseIdent: response.Ident,
	}

	if response.RenderFib != nil {
		if response.RenderFib.MaxChars > 0 {
			textEntry.ExpectedLength = response.RenderFib.MaxChars
		}
	}

	return textEntry
}

func (m *Migrator12to21) convertResponseToExtendedTextInteraction(response *models.Response) *models.ExtendedTextInteraction {
	extText := &models.ExtendedTextInteraction{
		XMLName:       xml.Name{Local: "extendedTextInteraction"},
		ResponseIdent: response.Ident,
	}

	if response.RenderFib != nil {
		if response.RenderFib.Rows > 0 {
			extText.ExpectedLines = response.RenderFib.Rows
		}
		if response.RenderFib.MaxChars > 0 {
			extText.ExpectedLength = response.RenderFib.MaxChars
		}
	}

	return extText
}

func (m *Migrator12to21) extractMaterialContent(material *models.Material) string {
	var content strings.Builder

	for _, matText := range material.MatText {
		text := matText.Content
		if matText.TextType == "text/html" {
			text = m.sanitizeHTMLContent(text)
		}
		content.WriteString(text)
	}

	for _, matImage := range material.MatImage {
		imgTag := fmt.Sprintf(`<img src="%s"`, matImage.URI)
		if matImage.Width > 0 {
			imgTag += fmt.Sprintf(` width="%d"`, matImage.Width)
		}
		if matImage.Height > 0 {
			imgTag += fmt.Sprintf(` height="%d"`, matImage.Height)
		}
		imgTag += " />"
		content.WriteString(imgTag)
	}

	return content.String()
}

func (m *Migrator12to21) extractResponseDeclarations(presentation *models.Presentation, responseProc *models.ResponseProc) []models.ResponseDecl {
	var responseDecls []models.ResponseDecl

	for _, response := range presentation.Response {
		responseDecl := models.ResponseDecl{
			XMLName:     xml.Name{Local: "responseDeclaration"},
			Identifier:  response.Ident,
			Cardinality: m.determineCardinality(&response),
			BaseType:    m.determineBaseType(&response),
		}

		if responseProc != nil {
			correctResponse := m.extractCorrectResponse(response.Ident, responseProc)
			if correctResponse != nil {
				responseDecl.CorrectResponse = correctResponse
			}
		}

		responseDecls = append(responseDecls, responseDecl)
	}

	return responseDecls
}

func (m *Migrator12to21) determineCardinality(response *models.Response) string {
	if response.RCardinality != "" {
		switch response.RCardinality {
		case "single":
			return "single"
		case "multiple":
			return "multiple"
		case "ordered":
			return "ordered"
		}
	}

	if response.RenderChoice != nil && response.RenderChoice.MaxNumber > 1 {
		return "multiple"
	}

	return "single"
}

func (m *Migrator12to21) determineBaseType(response *models.Response) string {
	if response.RenderChoice != nil {
		return "identifier"
	} else if response.RenderFib != nil {
		if response.RenderFib.FibType == "integer" {
			return "integer"
		} else if response.RenderFib.FibType == "decimal" {
			return "float"
		}
		return "string"
	}
	return "string"
}

func (m *Migrator12to21) extractCorrectResponse(responseIdent string, responseProc *models.ResponseProc) *models.CorrectResponse {
	var correctValues []string

	for _, condition := range responseProc.ResCondition {
		if condition.ConditionVar != nil {
			for _, varEqual := range condition.ConditionVar.VarEqual {
				if varEqual.RespIdent == responseIdent {
					isCorrect := false
					for _, setVar := range condition.SetVar {
						if setVar.Action == "set" && setVar.Value == "1" {
							isCorrect = true
							break
						}
					}
					if isCorrect {
						correctValues = append(correctValues, varEqual.Value)
					}
				}
			}
		}
	}

	if len(correctValues) > 0 {
		return &models.CorrectResponse{
			XMLName: xml.Name{Local: "correctResponse"},
			Value:   correctValues,
		}
	}

	return nil
}

func (m *Migrator12to21) extractOutcomeDeclarations(responseProc *models.ResponseProc) []models.OutcomeDecl {
	var outcomeDecls []models.OutcomeDecl

	if responseProc.Outcomes != nil {
		for _, decVar := range responseProc.Outcomes.DecVar {
			outcomeDecl := models.OutcomeDecl{
				XMLName:     xml.Name{Local: "outcomeDeclaration"},
				Identifier:  decVar.VarName,
				Cardinality: "single",
				BaseType:    m.convertVarType(decVar.VarType),
			}

			if decVar.DefaultVal != "" {
				outcomeDecl.DefaultValue = &models.DefaultValue{
					XMLName: xml.Name{Local: "defaultValue"},
					Value:   decVar.DefaultVal,
				}
			}

			outcomeDecls = append(outcomeDecls, outcomeDecl)
		}
	}

	if len(outcomeDecls) == 0 {
		outcomeDecls = append(outcomeDecls, models.OutcomeDecl{
			XMLName:     xml.Name{Local: "outcomeDeclaration"},
			Identifier:  "SCORE",
			Cardinality: "single",
			BaseType:    "float",
			DefaultValue: &models.DefaultValue{
				XMLName: xml.Name{Local: "defaultValue"},
				Value:   "0.0",
			},
		})
	}

	return outcomeDecls
}

func (m *Migrator12to21) convertVarType(varType string) string {
	switch varType {
	case "integer":
		return "integer"
	case "decimal", "scientific":
		return "float"
	case "boolean":
		return "boolean"
	default:
		return "float"
	}
}

func (m *Migrator12to21) migrateFeedback(feedback *models.Feedback) models.Feedback {
	migratedFeedback := models.Feedback{
		XMLName: feedback.XMLName,
		Ident:   feedback.Ident,
		Title:   feedback.Title,
	}

	if feedback.Material != nil {
		migratedFeedback.Material = feedback.Material
	}

	migratedFeedback.FlowMat = append(migratedFeedback.FlowMat, feedback.FlowMat...)

	return migratedFeedback
}

func (m *Migrator12to21) sanitizeHTMLContent(content string) string {
	content = strings.ReplaceAll(content, "<br>", "<br/>")
	content = strings.ReplaceAll(content, "<hr>", "<hr/>")
	content = strings.ReplaceAll(content, "<img ", "<img ")
	
	if !strings.Contains(content, "/>") && strings.Contains(content, "<img") {
		content = strings.ReplaceAll(content, ">", "/>")
	}

	return content
}