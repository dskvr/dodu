// Package export serializes a scan.Snapshot to stable, documented formats
// (JSON and CSV) for use by external tools.
//
// The JSON schema is versioned via SchemaVersion; bump on any incompatible
// change. CSV column order is stable per writer.
package export

// SchemaVersion is the JSON schema version. Bump on incompatible changes.
const SchemaVersion = "1.0"
