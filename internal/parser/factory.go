package parser

import (
	"fmt"
	"strings"

	"github.com/qti-migrator/internal/parser/qti12"
	"github.com/qti-migrator/internal/parser/qti21"
	"github.com/qti-migrator/internal/parser/qti30"
)

func GetParser(version string) (Parser, error) {
	version = strings.TrimSpace(version)
	
	switch {
	case strings.HasPrefix(version, "1.2"):
		return qti12.New(), nil
	case strings.HasPrefix(version, "2.1") || strings.HasPrefix(version, "2.2"):
		return qti21.New(), nil
	case strings.HasPrefix(version, "3.0"):
		return qti30.New(), nil
	default:
		return nil, fmt.Errorf("unsupported QTI version: %s", version)
	}
}