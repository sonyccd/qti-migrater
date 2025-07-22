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
	var doc models.QTIDocument
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

	return &doc, nil
}

func isValidQTI21Version(version string) bool {
	switch version {
	case "2.1", "2.1.0", "2.1.1", "2.2", "2.2.0", "2.2.1", "2.2.2", "2.2.3", "2.2.4":
		return true
	default:
		return false
	}
}