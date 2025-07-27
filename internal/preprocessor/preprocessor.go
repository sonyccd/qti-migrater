package preprocessor

import (
	"fmt"
	"strings"

	"github.com/qti-migrator/internal/parser"
	"github.com/qti-migrator/pkg/models"
)

type Preprocessor struct {
	verbosity int
}

type AnalysisReport struct {
	SourceVersion      string
	TargetVersion      string
	TotalItems         int
	CompatibleItems    int
	IncompatibleItems  int
	Warnings           []Warning
	Errors             []Error
	MigrationDetails   []MigrationDetail
}

type Warning struct {
	ItemID      string
	ElementPath string
	Message     string
	Suggestion  string
}

type Error struct {
	ItemID      string
	ElementPath string
	Message     string
	Fatal       bool
}

type MigrationDetail struct {
	ItemID      string
	ElementPath string
	OldValue    string
	NewValue    string
	Action      string
	Description string
}

func New(verbosity int) *Preprocessor {
	return &Preprocessor{
		verbosity: verbosity,
	}
}

func (p *Preprocessor) Analyze(content []byte, fromVersion, toVersion string) (*AnalysisReport, error) {
	sourceParser, err := parser.GetParser(fromVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get parser for version %s: %w", fromVersion, err)
	}

	doc, err := sourceParser.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	report := &AnalysisReport{
		SourceVersion: fromVersion,
		TargetVersion: toVersion,
		TotalItems:    len(doc.Items),
	}

	if fromVersion == "1.2" && toVersion == "2.1" {
		p.analyzeQTI12to21(doc, report)
	} else if fromVersion == "2.1" && toVersion == "3.0" {
		p.analyzeQTI21to30(doc, report)
	} else {
		return nil, fmt.Errorf("unsupported migration path: %s to %s", fromVersion, toVersion)
	}

	report.CompatibleItems = report.TotalItems - report.IncompatibleItems

	return report, nil
}

func (p *Preprocessor) analyzeQTI12to21(doc *models.QTIDocument, report *AnalysisReport) {
	for _, item := range doc.Items {
		p.analyzeItem12to21(&item, report)
	}

	if doc.Assessment != nil {
		for _, section := range doc.Assessment.Sections {
			for _, item := range section.Items {
				p.analyzeItem12to21(&item, report)
			}
		}
	}
}

func (p *Preprocessor) analyzeItem12to21(item *models.Item, report *AnalysisReport) {
	if item.Presentation != nil {
		for _, response := range item.Presentation.Response {
			if response.RenderChoice != nil && response.RenderChoice.Shuffle == "yes" {
				report.MigrationDetails = append(report.MigrationDetails, MigrationDetail{
					ItemID:      item.Ident,
					ElementPath: fmt.Sprintf("item[@ident='%s']/presentation/response_lid[@ident='%s']/render_choice", item.Ident, response.Ident),
					OldValue:    `shuffle="yes"`,
					NewValue:    `shuffle="true"`,
					Action:      "transform",
					Description: "Convert shuffle attribute from yes/no to true/false",
				})
			}
		}

		if item.Presentation.Material != nil {
			p.analyzeMaterial12to21(item.Ident, "presentation/material", item.Presentation.Material, report)
		}
	}

	if item.ResponseProc != nil {
		if item.ResponseProc.ScoreModel == "" {
			report.Warnings = append(report.Warnings, Warning{
				ItemID:      item.Ident,
				ElementPath: fmt.Sprintf("item[@ident='%s']/resprocessing", item.Ident),
				Message:     "Score model not specified in QTI 1.2",
				Suggestion:  "Default score model 'SumOfScores' will be applied",
			})
		}

		for _, condition := range item.ResponseProc.ResCondition {
			if condition.Continue == "yes" || condition.Continue == "no" {
				report.MigrationDetails = append(report.MigrationDetails, MigrationDetail{
					ItemID:      item.Ident,
					ElementPath: fmt.Sprintf("item[@ident='%s']/resprocessing/respcondition", item.Ident),
					OldValue:    fmt.Sprintf(`continue="%s"`, condition.Continue),
					NewValue:    fmt.Sprintf(`continue="%t"`, condition.Continue == "yes"),
					Action:      "transform",
					Description: "Convert continue attribute from yes/no to true/false",
				})
			}
		}
	}

	if item.Metadata != nil && item.Metadata.QTIMetadata != nil {
		meta := item.Metadata.QTIMetadata
		if meta.InteractionType != "" {
			if !isValidQTI21InteractionType(meta.InteractionType) {
				report.Warnings = append(report.Warnings, Warning{
					ItemID:      item.Ident,
					ElementPath: fmt.Sprintf("item[@ident='%s']/metadata/qtimetadata/interactiontype", item.Ident),
					Message:     fmt.Sprintf("Interaction type '%s' may need adjustment for QTI 2.1", meta.InteractionType),
					Suggestion:  "Review interaction type mapping for QTI 2.1 compliance",
				})
			}
		}
	}
}

func (p *Preprocessor) analyzeMaterial12to21(itemID, path string, material *models.Material, report *AnalysisReport) {
	for i, matText := range material.MatText {
		if matText.TextType == "text/html" && p.verbosity >= 2 {
			report.MigrationDetails = append(report.MigrationDetails, MigrationDetail{
				ItemID:      itemID,
				ElementPath: fmt.Sprintf("%s/mattext[%d]", path, i+1),
				OldValue:    "text/html content",
				NewValue:    "XHTML content (validated)",
				Action:      "validate",
				Description: "HTML content will be validated and converted to XHTML if necessary",
			})
		}
	}

	for i, matImage := range material.MatImage {
		if matImage.ImageType == "" && p.verbosity >= 2 {
			report.Warnings = append(report.Warnings, Warning{
				ItemID:      itemID,
				ElementPath: fmt.Sprintf("%s/matimage[%d]", path, i+1),
				Message:     "Image type not specified",
				Suggestion:  "Image type will be inferred from file extension or set to 'image/jpeg' as default",
			})
		}
	}
}

func (p *Preprocessor) analyzeQTI21to30(doc *models.QTIDocument, report *AnalysisReport) {
	for _, item := range doc.Items {
		p.analyzeItem21to30(&item, report)
	}

	if doc.Assessment != nil {
		for _, section := range doc.Assessment.Sections {
			for _, item := range section.Items {
				p.analyzeItem21to30(&item, report)
			}
		}
		
		// Assessment structure will change in QTI 3.0
		report.MigrationDetails = append(report.MigrationDetails, MigrationDetail{
			ItemID:      doc.Assessment.Ident,
			ElementPath: "assessment",
			OldValue:    "assessment",
			NewValue:    "qti-assessment-test",
			Action:      "rename",
			Description: "Assessment element renamed to qti-assessment-test in QTI 3.0",
		})
	}
}

func (p *Preprocessor) analyzeItem21to30(item *models.Item, report *AnalysisReport) {
	// Item element naming changes
	report.MigrationDetails = append(report.MigrationDetails, MigrationDetail{
		ItemID:      item.Ident,
		ElementPath: fmt.Sprintf("item[@ident='%s']", item.Ident),
		OldValue:    "item",
		NewValue:    "qti-assessment-item",
		Action:      "rename",
		Description: "Item element renamed to qti-assessment-item in QTI 3.0",
	})

	// ItemBody changes
	if item.ItemBody != nil {
		report.MigrationDetails = append(report.MigrationDetails, MigrationDetail{
			ItemID:      item.Ident,
			ElementPath: fmt.Sprintf("item[@ident='%s']/itemBody", item.Ident),
			OldValue:    "itemBody",
			NewValue:    "qti-item-body",
			Action:      "rename",
			Description: "ItemBody element renamed to qti-item-body in QTI 3.0",
		})

		// Analyze interactions
		for _, interaction := range item.ItemBody.ChoiceInteraction {
			report.MigrationDetails = append(report.MigrationDetails, MigrationDetail{
				ItemID:      item.Ident,
				ElementPath: fmt.Sprintf("item[@ident='%s']/itemBody/choiceInteraction[@responseIdentifier='%s']", item.Ident, interaction.ResponseIdent),
				OldValue:    "choiceInteraction",
				NewValue:    "qti-choice-interaction",
				Action:      "rename",
				Description: "Interaction elements prefixed with 'qti-' in QTI 3.0",
			})

			report.MigrationDetails = append(report.MigrationDetails, MigrationDetail{
				ItemID:      item.Ident,
				ElementPath: fmt.Sprintf("item[@ident='%s']/itemBody/choiceInteraction/simpleChoice", item.Ident),
				OldValue:    "simpleChoice",
				NewValue:    "qti-simple-choice",
				Action:      "rename",
				Description: "SimpleChoice elements renamed to qti-simple-choice in QTI 3.0",
			})
		}

		// Check for HTML content that needs updating
		for _, para := range item.ItemBody.P {
			if strings.Contains(para.Content, "class=") && p.verbosity >= 2 {
				report.Warnings = append(report.Warnings, Warning{
					ItemID:      item.Ident,
					ElementPath: fmt.Sprintf("item[@ident='%s']/itemBody/p", item.Ident),
					Message:     "HTML class attributes found in content",
					Suggestion:  "Class attributes will be converted to data-qti-class in QTI 3.0",
				})
			}
		}
	}

	// Response declaration changes
	for _, decl := range item.ResponseDecl {
		report.MigrationDetails = append(report.MigrationDetails, MigrationDetail{
			ItemID:      item.Ident,
			ElementPath: fmt.Sprintf("item[@ident='%s']/responseDeclaration[@identifier='%s']", item.Ident, decl.Identifier),
			OldValue:    "responseDeclaration",
			NewValue:    "qti-response-declaration",
			Action:      "rename",
			Description: "Response declaration renamed in QTI 3.0",
		})

		// Base type changes
		if decl.BaseType == "pair" {
			report.MigrationDetails = append(report.MigrationDetails, MigrationDetail{
				ItemID:      item.Ident,
				ElementPath: fmt.Sprintf("item[@ident='%s']/responseDeclaration[@identifier='%s']/@baseType", item.Ident, decl.Identifier),
				OldValue:    "pair",
				NewValue:    "directedPair",
				Action:      "transform",
				Description: "BaseType 'pair' renamed to 'directedPair' in QTI 3.0",
			})
		}
	}

	// Outcome declaration changes
	for _, decl := range item.OutcomeDecl {
		report.MigrationDetails = append(report.MigrationDetails, MigrationDetail{
			ItemID:      item.Ident,
			ElementPath: fmt.Sprintf("item[@ident='%s']/outcomeDeclaration[@identifier='%s']", item.Ident, decl.Identifier),
			OldValue:    "outcomeDeclaration",
			NewValue:    "qti-outcome-declaration",
			Action:      "rename",
			Description: "Outcome declaration renamed in QTI 3.0",
		})
	}

	// Metadata changes
	if item.Metadata != nil && item.Metadata.QTIMetadata != nil {
		report.MigrationDetails = append(report.MigrationDetails, MigrationDetail{
			ItemID:      item.Ident,
			ElementPath: fmt.Sprintf("item[@ident='%s']/metadata/qtimetadata", item.Ident),
			OldValue:    "qtimetadata",
			NewValue:    "qti-metadata-container",
			Action:      "rename",
			Description: "QTI metadata container renamed in QTI 3.0",
		})
	}
}

func isValidQTI21InteractionType(interactionType string) bool {
	validTypes := map[string]bool{
		"choiceInteraction":        true,
		"orderInteraction":         true,
		"associateInteraction":     true,
		"matchInteraction":         true,
		"gapMatchInteraction":      true,
		"inlineChoiceInteraction":  true,
		"textEntryInteraction":     true,
		"extendedTextInteraction":  true,
		"hotspotInteraction":       true,
		"selectPointInteraction":   true,
		"graphicOrderInteraction":  true,
		"graphicAssociateInteraction": true,
		"positionObjectInteraction": true,
		"sliderInteraction":        true,
		"drawingInteraction":       true,
		"uploadInteraction":        true,
		"customInteraction":        true,
	}
	return validTypes[interactionType]
}

func (r *AnalysisReport) HasErrors() bool {
	for _, err := range r.Errors {
		if err.Fatal {
			return true
		}
	}
	return false
}