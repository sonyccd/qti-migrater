# QTI Migrator

A command-line tool for migrating QTI (Question and Test Interoperability) files between different versions.

## Features

- **Version Support**: Supports migration from QTI 1.2 to QTI 2.1 and QTI 2.1 to QTI 3.0
- **Preprocessing Analysis**: Analyze files before migration to identify potential issues
- **Detailed Reports**: Configurable verbosity levels for migration reports
- **Pipe Support**: Can be used in scripts with stdin/stdout support
- **Error Handling**: Clear error messages for migration issues
- **Modular Architecture**: Easy to extend for new QTI versions
- **Version-Specific Models**: Clean separation of QTI version structures for better maintainability

## Installation

```bash
go get github.com/qti-migrator
```

Or build from source:

```bash
git clone https://github.com/qti-migrator
cd qti-migrator
go build -o qti-migrator cmd/qti-migrator/main.go
```

## Usage

### Basic Migration

```bash
# Migrate a file from QTI 1.2 to 2.1
qti-migrator migrate -f 1.2 -t 2.1 -i input.xml -o output.xml

# Migrate a file from QTI 2.1 to 3.0
qti-migrator migrate -f 2.1 -t 3.0 -i input.xml -o output.xml

# Use stdin/stdout for scripting
cat input.xml | qti-migrator migrate -f 1.2 -t 2.1 > output.xml
```

### Preview Mode

Preview the migration without making changes:

```bash
qti-migrator migrate -f 1.2 -t 2.1 -i input.xml --preview
```

### Verbosity Levels

Control the amount of detail in reports:

```bash
# Minimal output (0)
qti-migrator migrate -f 1.2 -t 2.1 -i input.xml -o output.xml -v 0

# Normal output (1) - default
qti-migrator migrate -f 1.2 -t 2.1 -i input.xml -o output.xml -v 1

# Detailed output (2)
qti-migrator migrate -f 1.2 -t 2.1 -i input.xml -o output.xml -v 2

# Debug output (3)
qti-migrator migrate -f 1.2 -t 2.1 -i input.xml -o output.xml -v 3
```

### Batch Processing

Use with shell scripts for batch processing:

```bash
# Process all QTI 1.2 files in a directory
for file in *.xml; do
    qti-migrator migrate -f 1.2 -t 2.1 -i "$file" -o "migrated_$file"
done

# Process all QTI 2.1 files to QTI 3.0
for file in *_v21.xml; do
    qti-migrator migrate -f 2.1 -t 3.0 -i "$file" -o "${file/_v21/_v30}"
done

# Using find and xargs
find . -name "*.xml" -print0 | xargs -0 -I {} qti-migrator migrate -f 1.2 -t 2.1 -i {} -o migrated_{}
```

## Migration Report

The tool generates detailed migration reports that include:

- **Summary**: Overall migration status and statistics
- **Errors**: Migration blockers that must be resolved
- **Warnings**: Potential issues that may need attention
- **Migration Details**: Specific changes that will be made (verbosity 2+)

Example report:
```
================================================================================
                          QTI Migration Analysis Report
================================================================================
Migration Path: QTI 1.2 → QTI 2.1
================================================================================

SUMMARY
-------
Status: READY
Total Items: 10
Compatible Items: 8
Items Requiring Attention: 2
Errors: 0
Warnings: 3

WARNINGS
--------
1. [Item: q001] Score model not specified in QTI 1.2
   → Default score model 'SumOfScores' will be applied

2. [Item: q003] Interaction type 'multiple_choice' may need adjustment for QTI 2.1
   → Review interaction type mapping for QTI 2.1 compliance
```

## Supported Migrations

### QTI 1.2 to 2.1

- Converts presentation elements to itemBody
- Transforms response processing
- Updates attribute values (e.g., yes/no to true/false)
- Generates response and outcome declarations
- Validates and converts HTML content to XHTML

### QTI 2.1 to 3.0

- Updates XML namespaces to QTI 3.0 specification
- Converts element names to QTI 3.0 conventions (e.g., `itemBody` → `qti-item-body`)
- Transforms interaction types to new naming scheme (e.g., `choiceInteraction` → `qti-choice-interaction`)
- Updates base types and attributes for QTI 3.0 compliance
- Converts HTML class attributes to data-qti-class
- Transforms object elements to qti-object elements
- Migrates metadata structures to QTI 3.0 format

## Architecture

The tool is built with a modular architecture:

- **Parser**: Handles parsing of different QTI versions with version-specific parsers
  - `qti12.Parser12` for QTI 1.2 documents
  - `qti21.Parser21` for QTI 2.1/2.2 documents  
  - `qti30.Parser30` for QTI 3.0 documents
- **Models**: Version-specific data structures
  - `models/common.go` - Shared structures (LOM metadata, materials)
  - `models/qti12.go` - QTI 1.2 specific structures
  - `models/qti21.go` - QTI 2.1/2.2 specific structures
  - `models/qti30.go` - QTI 3.0 specific structures
  - `models/qti.go` - Generic structures for backward compatibility
- **Preprocessor**: Analyzes documents for migration compatibility
- **Migrator**: Performs the actual migration transformations
- **Reporter**: Generates human-readable reports
- **Error Handler**: Provides detailed error information

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues.

## License

MIT License