package qti30

import (
	"encoding/xml"
	"fmt"

	"github.com/qti-migrator/pkg/models"
)

type Parser30 struct{}

func New() *Parser30 {
	return &Parser30{}
}

func (p *Parser30) Version() string {
	return "3.0"
}

func (p *Parser30) Parse(content []byte) (*models.QTIDocument, error) {
	var doc models.QTIDocument30
	err := xml.Unmarshal(content, &doc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse QTI 3.0 document: %w", err)
	}

	if doc.Version == "" {
		doc.Version = "3.0"
	}

	if !isValidQTI30Version(doc.Version) {
		return nil, fmt.Errorf("invalid QTI version for 3.0 parser: %s", doc.Version)
	}

	// QTI 3.0 has a completely different structure, so we need to transform it
	// to the generic format for backward compatibility
	genericDoc := &models.QTIDocument{
		XMLName:  xml.Name{Local: "questestinterop"},
		Version:  doc.Version,
		Items:    []models.Item{convertItem30ToGeneric(doc)},
		Metadata: doc.Metadata,
	}

	return genericDoc, nil
}

func isValidQTI30Version(version string) bool {
	switch version {
	case "3.0", "3.0.0":
		return true
	default:
		return false
	}
}

// Converter function to transform a QTI 3.0 document (which is an item) to generic Item
func convertItem30ToGeneric(doc models.QTIDocument30) models.Item {
	return models.Item{
		XMLName:      xml.Name{Local: "item"},
		Title:        doc.Title,
		Ident:        doc.Identifier,
		Metadata:     doc.Metadata,
		ItemBody:     convertItemBody30ToGeneric(doc.ItemBody),
		ResponseDecl: convertResponseDecl30ToGeneric(doc.ResponseDeclarations),
		OutcomeDecl:  convertOutcomeDecl30ToGeneric(doc.OutcomeDeclarations),
		TemplateDecl: convertTemplateDecl30ToGeneric(doc.TemplateDeclarations),
		// Note: QTI 3.0 uses modalFeedback instead of itemfeedback
		Feedback: convertModalFeedback30ToGeneric(doc.ModalFeedback),
	}
}

func convertItemBody30ToGeneric(body *models.ItemBody30) *models.ItemBody {
	if body == nil {
		return nil
	}
	// QTI 3.0 ItemBody has a more flexible structure
	// For now, we'll create a simple generic ItemBody
	return &models.ItemBody{
		XMLName: body.XMLName,
		// The content needs special handling as it's an interface{} slice in QTI 3.0
		// This would require more sophisticated conversion based on actual content types
	}
}

func convertResponseDecl30ToGeneric(decls []models.ResponseDecl30) []models.ResponseDecl {
	genericDecls := make([]models.ResponseDecl, len(decls))
	for i, decl := range decls {
		genericDecls[i] = models.ResponseDecl{
			XMLName:      decl.XMLName,
			Identifier:   decl.Identifier,
			Cardinality:  decl.Cardinality,
			BaseType:     decl.BaseType,
			CorrectResponse: convertCorrectResponse30ToGeneric(decl.CorrectResponse),
			Mapping:      convertMapping30ToGeneric(decl.Mapping),
		}
	}
	return genericDecls
}

func convertCorrectResponse30ToGeneric(cr *models.CorrectResponse30) *models.CorrectResponse {
	if cr == nil {
		return nil
	}
	values := make([]string, len(cr.Value))
	for i, v := range cr.Value {
		values[i] = v.Content
	}
	return &models.CorrectResponse{
		XMLName: cr.XMLName,
		Value:   values,
	}
}

func convertMapping30ToGeneric(m *models.Mapping30) *models.Mapping {
	if m == nil {
		return nil
	}
	entries := make([]models.MapEntry, len(m.MapEntry))
	for i, e := range m.MapEntry {
		entries[i] = models.MapEntry(e)
	}
	return &models.Mapping{
		XMLName:      m.XMLName,
		LowerBound:   m.LowerBound,
		UpperBound:   m.UpperBound,
		DefaultValue: m.DefaultValue,
		MapEntry:     entries,
	}
}

func convertOutcomeDecl30ToGeneric(decls []models.OutcomeDecl30) []models.OutcomeDecl {
	genericDecls := make([]models.OutcomeDecl, len(decls))
	for i, decl := range decls {
		genericDecls[i] = models.OutcomeDecl{
			XMLName:      decl.XMLName,
			Identifier:   decl.Identifier,
			Cardinality:  decl.Cardinality,
			BaseType:     decl.BaseType,
			DefaultValue: convertDefaultValue30ToGeneric(decl.DefaultValue),
		}
	}
	return genericDecls
}

func convertDefaultValue30ToGeneric(dv *models.DefaultValue30) *models.DefaultValue {
	if dv == nil {
		return nil
	}
	// Take the first value for simplicity
	value := ""
	if len(dv.Value) > 0 {
		value = dv.Value[0].Content
	}
	return &models.DefaultValue{
		XMLName: dv.XMLName,
		Value:   value,
	}
}

func convertTemplateDecl30ToGeneric(decls []models.TemplateDecl30) []models.TemplateDecl {
	genericDecls := make([]models.TemplateDecl, len(decls))
	for i, decl := range decls {
		genericDecls[i] = models.TemplateDecl{
			XMLName:       decl.XMLName,
			Identifier:    decl.Identifier,
			Cardinality:   decl.Cardinality,
			BaseType:      decl.BaseType,
			ParamVariable: decl.ParamVariable,
			DefaultValue:  convertDefaultValue30ToGeneric(decl.DefaultValue),
		}
	}
	return genericDecls
}

func convertModalFeedback30ToGeneric(feedbacks []models.ModalFeedback30) []models.Feedback {
	genericFeedbacks := make([]models.Feedback, len(feedbacks))
	for i, feedback := range feedbacks {
		// Convert modalFeedback to itemfeedback format
		genericFeedbacks[i] = models.Feedback{
			XMLName: xml.Name{Local: "itemfeedback"},
			Ident:   feedback.Identifier,
			Title:   feedback.Title,
			Material: &models.Material{
				MatText: []models.MatText{
					{
						Content: feedback.Content,
					},
				},
			},
		}
	}
	return genericFeedbacks
}