package parser

import (
	"github.com/qti-migrator/pkg/models"
)

type Parser interface {
	Parse(content []byte) (*models.QTIDocument, error)
	Version() string
}