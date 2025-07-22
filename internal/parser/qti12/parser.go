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
	var doc models.QTIDocument
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

	return &doc, nil
}

func isValidQTI12Version(version string) bool {
	switch version {
	case "1.2", "1.2.1", "1.2.0":
		return true
	default:
		return false
	}
}