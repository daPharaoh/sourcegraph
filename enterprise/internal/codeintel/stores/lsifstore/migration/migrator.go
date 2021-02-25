package migration

import (
	"context"
	"database/sql"

	"github.com/keegancsmith/sqlf"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/stores/lsifstore"
	"github.com/sourcegraph/sourcegraph/internal/database/basestore"
	"github.com/sourcegraph/sourcegraph/internal/oobmigration"
)

type Migrator struct {
	store  *lsifstore.Store
	driver migrationDriver
}

type migrationDriver interface {
	SomethingUp(rows *sql.Rows) ([]updateSpec, error)
	SomethingDown(rows *sql.Rows) ([]updateSpec, error)
}

type updateSpec struct {
	DumpID      int
	Conditions  map[string]interface{}
	Assignments map[string]interface{}
}

// TODO
func newMigrator(store *lsifstore.Store, driver migrationDriver) oobmigration.Migrator {
	return &Migrator{
		store:  store,
		driver: driver,
	}
}

// TODO
// Progress returns the ratio of migrated records to total records. Any record with a
// schema version of two or greater is considered migrated.
func (m *Migrator) Progress(ctx context.Context) (float64, error) {
	// TODO - configure
	tableName := sqlf.Sprintf("lsif_data_documents")
	version := 2

	progress, _, err := basestore.ScanFirstFloat(m.store.Query(ctx, sqlf.Sprintf(migratorProgressQuery, tableName, version, tableName)))
	if err != nil {
		return 0, err
	}

	return progress, nil
}

const migratorProgressQuery = `
-- source: enterprise/internal/codeintel/stores/lsifstore/migration/migrator.go:Progress
SELECT CASE c2.count WHEN 0 THEN 1 ELSE cast(c1.count as float) / cast(c2.count as float) END FROM
	(SELECT COUNT(*) as count FROM %s_schema_versions WHERE min_schema_version >= %s) c1,
	(SELECT COUNT(*) as count FROM %s_schema_versions) c2
`

// TODO
func (m *Migrator) Up(ctx context.Context) (err error) {
	// TODO - configure
	tableName := sqlf.Sprintf("lsif_data_documents")
	version := 2
	batchSize := 1000
	fields := []string{"path", "data"}

	return m.run(ctx, tableName, fields, version, batchSize)
}

// TODO
func (m *Migrator) Down(ctx context.Context) error {
	// TODO - configure
	tableName := sqlf.Sprintf("lsif_data_documents")
	version := 2
	batchSize := 1000
	fields := []string{"path", "data"}

	// TODO - do it down though :/
	return m.run(ctx, tableName, fields, version, batchSize)
}

//
//
//

// TODO - document
// TODO - break apart
func (m *Migrator) run(ctx context.Context, tableName *sqlf.Query, fields []string, version, batchSize int) (err error) {
	fieldQueries := make([]*sqlf.Query, 0, len(fields))
	for _, field := range fields {
		fieldQueries = append(fieldQueries, sqlf.Sprintf(field))
	}

	tx, err := m.store.Transact(ctx)
	if err != nil {
		return err
	}
	defer func() { err = tx.Done(err) }()

	rows, err := tx.Query(ctx, sqlf.Sprintf(migratorSelectQuery, version, tableName, sqlf.Join(fieldQueries, ", "), tableName, version, batchSize))
	if err != nil {
		return err
	}
	defer func() { err = basestore.CloseRows(rows, err) }()

	updateQueries, err := m.driver.SomethingUp(rows)
	if err != nil {
		return err
	}

	for _, spec := range updateQueries {
		// TODO - rewrite this
		spec.Assignments["schema_version"] = version
		assignments := sqlf.Join(joiner(spec.Assignments), ", ")
		spec.Conditions["dump_id"] = spec.DumpID
		conditions := sqlf.Join(joiner(spec.Conditions), " AND ")

		if err := tx.Exec(ctx, sqlf.Sprintf(migratorUpdateQuery, tableName, assignments, conditions)); err != nil {
			return err
		}
	}

	idMap := map[int]struct{}{}
	for _, spec := range updateQueries {
		idMap[spec.DumpID] = struct{}{}
	}
	ids := make([]*sqlf.Query, 0, len(idMap))
	for key := range idMap {
		ids = append(ids, sqlf.Sprintf("%s", key))
	}

	if len(ids) == 0 {
		return nil
	}

	return m.store.Exec(ctx, sqlf.Sprintf(updateRangesQuery, sqlf.Join(ids, ", ")))
}

const migratorSelectQuery = `
-- source: enterprise/internal/codeintel/stores/lsifstore/migration/migrator.go:do
WITH candidates AS (SELECT dump_id FROM %s_schema_versions WHERE min_schema_version < %s)
SELECT dump_id, %s
FROM %s
WHERE dump_id IN (SELECT dump_id FROM candidates) AND schema_version < %s
LIMIT %s
FOR UPDATE SKIP LOCKED
`

const migratorUpdateQuery = `
-- source: enterprise/internal/codeintel/stores/lsifstore/migration/migrator.go:do
UPDATE %s SET %s WHERE %s
`

const updateRangesQuery = `
-- source: enterprise/internal/codeintel/stores/lsifstore/migration/migrator.go:do
INSERT INTO
	lsif_data_documents_schema_versions
FROM
	dump_id,
	MIN(schema_version) as min_schema_version,
	MAX(schema_version) as max_schema_version
FROM
	lsif_data_documents
WHERE
	dump_id IN (%s)
ON CONFLICT (dump_id) DO UPDATE SET
	min_schema_version = EXCLUDED.min_schema_version,
	max_schema_version = EXCLUDED.max_schema_version
`

// TODO - rename
// TODO - document
func joiner(m map[string]interface{}) []*sqlf.Query {
	queries := make([]*sqlf.Query, 0, len(m))
	for k, v := range m {
		queries = append(queries, sqlf.Sprintf("%s = %s", sqlf.Sprintf(k), v))
	}

	return queries
}
