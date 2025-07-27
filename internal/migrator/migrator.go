package migrator

import (
	"fmt"

	"github.com/qti-migrator/internal/migrator/qti12to21"
	"github.com/qti-migrator/internal/migrator/qti21to30"
	"github.com/qti-migrator/internal/parser"
)

type Migrator interface {
	Migrate(doc interface{}) ([]byte, error)
}

type MigratorService struct{}

func New() *MigratorService {
	return &MigratorService{}
}

func (m *MigratorService) Migrate(content []byte, fromVersion, toVersion string) ([]byte, error) {
	sourceParser, err := parser.GetParser(fromVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get parser for version %s: %w", fromVersion, err)
	}

	doc, err := sourceParser.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source document: %w", err)
	}

	var migrator Migrator
	if fromVersion == "1.2" && toVersion == "2.1" {
		migrator = qti12to21.New()
	} else if fromVersion == "2.1" && toVersion == "3.0" {
		migrator = qti21to30.New()
	} else {
		return nil, fmt.Errorf("unsupported migration path: %s to %s", fromVersion, toVersion)
	}

	result, err := migrator.Migrate(doc)
	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return result, nil
}