-- ===========================
-- PIPELINES
-- ===========================

ALTER TABLE pipelines
    ADD COLUMN github_delivery_id TEXT,
    ADD COLUMN trigger_username VARCHAR(255),
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE pipelines
    RENAME COLUMN completed_at TO finished_at;

-- ===========================
-- JOBS
-- ===========================

ALTER TABLE jobs
    ADD COLUMN order_index INT NOT NULL DEFAULT 0,
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE jobs
    RENAME COLUMN completed_at TO finished_at;

ALTER TABLE jobs
    ADD CONSTRAINT jobs_order_index_check
    CHECK (order_index >= 0);

-- ===========================
-- LOGS
-- ===========================

ALTER TABLE logs
    ADD CONSTRAINT logs_job_line_unique
    UNIQUE (job_id, line_number);

-- ===========================
-- INDEXES
-- ===========================

CREATE INDEX idx_pipelines_commit_sha
ON pipelines(commit_sha);

CREATE INDEX idx_pipelines_delivery
ON pipelines(github_delivery_id);

CREATE INDEX idx_jobs_pipeline_order
ON jobs(pipeline_id, order_index);