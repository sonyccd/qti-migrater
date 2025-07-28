package qti12

import (
	"encoding/xml"
	"fmt"

	"github.com/qti-migrator/pkg/models"
)

type Parser12 struct{}

func New() *Parser12 {
	return &Parser12{}
}

func (p *Parser12) Version() string {
	return "1.2"
}

func (p *Parser12) Parse(content []byte) (*models.QTIDocument, error) {
	var doc models.QTIDocument12
	err := xml.Unmarshal(content, &doc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse QTI 1.2 document: %w", err)
	}

	if doc.Version == "" {
		doc.Version = "1.2"
	}

	if !isValidQTI12Version(doc.Version) {
		return nil, fmt.Errorf("invalid QTI version for 1.2 parser: %s", doc.Version)
	}

	// Convert to generic QTIDocument for backward compatibility
	genericDoc := &models.QTIDocument{
		XMLName:    doc.XMLName,
		Version:    doc.Version,
		Items:      convertItems12ToGeneric(doc.Items),
		Assessment: convertAssessment12ToGeneric(doc.Assessment),
		Metadata:   doc.Metadata,
	}

	return genericDoc, nil
}

func isValidQTI12Version(version string) bool {
	switch version {
	case "1.2", "1.2.1", "1.2.0":
		return true
	default:
		return false
	}
}

// Converter functions to transform QTI 1.2 specific types to generic types

func convertItems12ToGeneric(items []models.Item12) []models.Item {
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
			Feedback:     convertFeedback12ToGeneric(item.Feedback),
			RubricBlock:  item.RubricBlock,
		}
	}
	return genericItems
}

func convertAssessment12ToGeneric(assessment *models.Assessment12) *models.Assessment {
	if assessment == nil {
		return nil
	}
	return &models.Assessment{
		XMLName:     assessment.XMLName,
		Title:       assessment.Title,
		Ident:       assessment.Ident,
		Sections:    convertSections12ToGeneric(assessment.Sections),
		Metadata:    assessment.Metadata,
		Objectives:  assessment.Objectives,
		RubricBlock: assessment.RubricBlock,
	}
}

func convertSections12ToGeneric(sections []models.Section12) []models.Section {
	genericSections := make([]models.Section, len(sections))
	for i, section := range sections {
		genericSections[i] = models.Section{
			XMLName:  section.XMLName,
			Title:    section.Title,
			Ident:    section.Ident,
			Items:    convertItems12ToGeneric(section.Items),
			Metadata: section.Metadata,
		}
	}
	return genericSections
}

func convertFeedback12ToGeneric(feedbacks []models.Feedback12) []models.Feedback {
	genericFeedbacks := make([]models.Feedback, len(feedbacks))
	for i, feedback := range feedbacks {
		genericFeedbacks[i] = models.Feedback(feedback)
	}
	return genericFeedbacks
}