CREATE TABLE jobs (
    id TEXT PRIMARY KEY,
    operation TEXT NOT NULL,
    status TEXT NOT NULL,
    progress SMALLINT NOT NULL DEFAULT 0,
    attempt INTEGER NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ,
    result_key TEXT,
    failure_code TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT jobs_id_not_empty
        CHECK (id <> ''),

    CONSTRAINT jobs_operation_valid
        CHECK (
            operation IN (
                'image_grayscale',
                'image_convert',
                'image_compress',
                'image_resize',
                'image_crop',
                'image_rotate',
                'image_flip',
                'image_thumbnail',
                'image_strip_metadata',
                'image_adjust',
                'image_blur',
                'image_sharpen',
                'image_pixelate',
                'image_padding'
            )
        ),

    CONSTRAINT job_status_valid
        CHECK (
            status IN (
                'queued',
                'processing',
                'cancelling',
                'completed',
                'failed',
                'cancelled',
                'expired'
            )
        ),

    CONSTRAINT jobs_progress_range
        CHECK (progress BETWEEN 0 AND 100),

    CONSTRAINT jobs_attempt_nonnegative
        CHECK (attempt >= 0),
    
    CONSTRAINT job_time_order
        CHECK (updated_at >= created_at),
    
    CONSTRAINT job_outcome_exclusive
        CHECK (result_key IS NULL OR failure_code IS NULL)
);



CREATE INDEX jobs_status_created_at_idx
    ON jobs(status, created_at);


CREATE INDEX jobs_expires_at_idx
    ON jobs(expires_at)
    WHERE expires_at IS NOT NULL;
