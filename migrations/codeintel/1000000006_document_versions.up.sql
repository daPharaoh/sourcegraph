BEGIN;

--
-- TODO - put this in metadata instead?
--

CREATE TABLE lsif_data_documents_schema_versions (
    dump_id integer NOT NULL,
    min_schema_version integer,
    max_schema_version integer
);
ALTER TABLE lsif_data_documents_schema_versions ADD PRIMARY KEY (dump_id);

COMMIT;
