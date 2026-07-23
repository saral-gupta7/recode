ALTER TABLE jobs
    ADD COLUMN input_key TEXT,
    ADD COLUMN original_filename TEXT,
    ADD COLUMN mime_type TEXT,
    ADD COLUMN size_bytes BIGINT,
    ADD COLUMN options JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN owner_token_hash BYTEA;

ALTER TABLE jobs
    ADD CONSTRAINT jobs_size_positive
        CHECK (size_bytes IS NULL OR size_bytes > 0),
    ADD CONSTRAINT jobs_input_key_not_empty
        CHECK (input_key IS NULL OR input_key <> ''),
    ADD CONSTRAINT jobs_owner_hash_not_empty
        CHECK (owner_token_hash IS NULL OR octet_length(owner_token_hash) > 0);
