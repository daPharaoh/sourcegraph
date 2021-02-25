package migration

import (
	"database/sql"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/stores/lsifstore"
	"github.com/sourcegraph/sourcegraph/internal/oobmigration"
)

type diagnosticsCountMigrator struct {
	serializer *lsifstore.Serializer
}

// NewDiagnosticsCountMigrator creates a new Migrator instance that reads the documents
// table and populates their num_diagnostics value based on their decoded payload. This
// will update rows with a schema_version of 1, and will set the row's schema version
// to 2 after processing.
func NewDiagnosticsCountMigrator(store *lsifstore.Store) oobmigration.Migrator {
	driver := &diagnosticsCountMigrator{
		serializer: lsifstore.NewSerializer(),
	}

	//
	// TODO - other options
	return newMigrator(store, driver)
}

// TODO - rename
// TODO - document
func (m *diagnosticsCountMigrator) SomethingUp(rows *sql.Rows) ([]updateSpec, error) {
	//
	// TODO - why are we passing rows here?
	//

	var (
		specs   []updateSpec
		dumpID  int
		path    string
		rawData []byte
	)

	for rows.Next() {
		if err := rows.Scan(&dumpID, &path, &rawData); err != nil {
			return nil, err
		}

		data, err := m.serializer.UnmarshalDocumentData(rawData)
		if err != nil {
			return nil, err
		}

		specs = append(specs, updateSpec{
			DumpID:      dumpID,
			Conditions:  map[string]interface{}{"path": path},
			Assignments: map[string]interface{}{"num_diagnostics": len(data.Diagnostics)},
		})
	}

	return specs, nil
}

// TODO - rename
// TODO - document
func (m *diagnosticsCountMigrator) SomethingDown(rows *sql.Rows) ([]updateSpec, error) {
	var (
		specs   []updateSpec
		dumpID  int
		path    string
		rawData []byte // TODO - unnecessary
	)

	for rows.Next() {
		if err := rows.Scan(&dumpID, &path, &rawData); err != nil {
			return nil, err
		}

		specs = append(specs, updateSpec{
			DumpID:      dumpID,
			Conditions:  map[string]interface{}{"path": path},
			Assignments: map[string]interface{}{"num_diagnostics": 0},
		})
	}

	return specs, nil
}
