package qti21

import (
	"encoding/xml"
	"fmt"

	"github.com/qti-migrator/pkg/models"
)

type Parser21 struct{}

func New() *Parser21 {
	return &Parser21{}
}

func (p *Parser21) Version() string {
	return "2.1"
}

func (p *Parser21) Parse(content []byte) (*models.QTIDocument, error) {
	var doc models.QTIDocument21
	err := xml.Unmarshal(content, &doc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse QTI 2.1 document: %w", err)
	}

	if doc.Version == "" {
		doc.Version = "2.1"
	}

	if !isValidQTI21Version(doc.Version) {
		return nil, fmt.Errorf("invalid QTI version for 2.1 parser: %s", doc.Version)
	}

	// Convert to generic QTIDocument for backward compatibility
	genericDoc := &models.QTIDocument{
		XMLName:    doc.XMLName,
		Version:    doc.Version,
		Items:      convertItems21ToGeneric(doc.Items),
		Assessment: convertAssessment21ToGeneric(doc.Assessment),
		Metadata:   doc.Metadata,
	}

	return genericDoc, nil
}

func isValidQTI21Version(version string) bool {
	switch version {
	case "2.1", "2.1.0", "2.1.1", "2.2", "2.2.0", "2.2.1", "2.2.2", "2.2.3", "2.2.4":
		return true
	default:
		return false
	}
}

// Converter functions to transform QTI 2.1 specific types to generic types

func convertItems21ToGeneric(items []models.Item21) []models.Item {
	genericItems := make([]models.Item, len(items))
	for i, item := range items {
		genericItems[i] = models.Item{
			XMLName:      item.XMLName,
			Title:        item.Title,
			Ident:        item.Ident,
			MaxAttempts:  item.MaxAttempts,
			Metadata:     item.Metadata,
			Presentation: item.Presentation,
			ResponseProc: item.ResponseProc,
			ItemBody:     item.ItemBody,
			ResponseDecl: convertResponseDecl21ToGeneric(item.ResponseDecl),
			OutcomeDecl:  convertOutcomeDecl21ToGeneric(item.OutcomeDecl),
			TemplateDecl: convertTemplateDecl21ToGeneric(item.TemplateDecl),
			Feedback:     convertFeedback21ToGeneric(item.Feedback),
			RubricBlock:  item.RubricBlock,
		}
	}
	return genericItems
}

func convertAssessment21ToGeneric(assessment *models.Assessment21) *models.Assessment {
	if assessment == nil {
		return nil
	}
	return &models.Assessment{
		XMLName:     assessment.XMLName,
		Title:       assessment.Title,
		Ident:       assessment.Ident,
		Sections:    convertSections21ToGeneric(assessment.Sections),
		Metadata:    assessment.Metadata,
		Objectives:  assessment.Objectives,
		RubricBlock: assessment.RubricBlock,
	}
}

func convertSections21ToGeneric(sections []models.Section21) []models.Section {
	genericSections := make([]models.Section, len(sections))
	for i, section := range sections {
		genericSections[i] = models.Section{
			XMLName:  section.XMLName,
			Title:    section.Title,
			Ident:    section.Ident,
			Items:    convertItems21ToGeneric(section.Items),
			Metadata: section.Metadata,
		}
	}
	return genericSections
}

func convertResponseDecl21ToGeneric(decls []models.ResponseDecl21) []models.ResponseDecl {
	genericDecls := make([]models.ResponseDecl, len(decls))
	for i, decl := range decls {
		genericDecls[i] = models.ResponseDecl{
			XMLName:         decl.XMLName,
			Identifier:      decl.Identifier,
			Cardinality:     decl.Cardinality,
			BaseType:        decl.BaseType,
			CorrectResponse: (*models.CorrectResponse)(decl.CorrectResponse),
			Mapping:         (*models.Mapping)(decl.Mapping),
		}
	}
	return genericDecls
}

func convertOutcomeDecl21ToGeneric(decls []models.OutcomeDecl21) []models.OutcomeDecl {
	genericDecls := make([]models.OutcomeDecl, len(decls))
	for i, decl := range decls {
		genericDecls[i] = models.OutcomeDecl{
			XMLName:      decl.XMLName,
			Identifier:   decl.Identifier,
			Cardinality:  decl.Cardinality,
			BaseType:     decl.BaseType,
			DefaultValue: (*models.DefaultValue)(decl.DefaultValue),
		}
	}
	return genericDecls
}

func convertTemplateDecl21ToGeneric(decls []models.TemplateDecl21) []models.TemplateDecl {
	genericDecls := make([]models.TemplateDecl, len(decls))
	for i, decl := range decls {
		genericDecls[i] = models.TemplateDecl{
			XMLName:       decl.XMLName,
			Identifier:    decl.Identifier,
			Cardinality:   decl.Cardinality,
			BaseType:      decl.BaseType,
			ParamVariable: decl.ParamVariable,
			DefaultValue:  (*models.DefaultValue)(decl.DefaultValue),
		}
	}
	return genericDecls
}

func convertFeedback21ToGeneric(feedbacks []models.Feedback21) []models.Feedback {
	genericFeedbacks := make([]models.Feedback, len(feedbacks))
	for i, feedback := range feedbacks {
		genericFeedbacks[i] = models.Feedback(feedback)
	}
	return genericFeedbacks
}