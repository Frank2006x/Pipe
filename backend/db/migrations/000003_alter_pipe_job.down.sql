-- ===========================
-- INDEXES
-- ===========================

DROP INDEX IF EXISTS idx_jobs_pipeline_order;
DROP INDEX IF EXISTS idx_pipelines_delivery;
DROP INDEX IF EXISTS idx_pipelines_commit_sha;

-- ===========================
-- LOGS
-- ===========================

ALTER TABLE logs
    DROP CONSTRAINT IF EXISTS logs_job_line_unique;

-- ===========================
-- JOBS
-- ===========================

ALTER TABLE jobs
    DROP CONSTRAINT IF EXISTS jobs_order_index_check;

ALTER TABLE jobs
    RENAME COLUMN finished_at TO completed_at;

ALTER TABLE jobs
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS order_index;

-- ===========================
-- PIPELINES
-- ===========================

ALTER TABLE pipelines
    RENAME COLUMN finished_at TO completed_at;

ALTER TABLE pipelines
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS trigger_username,
    DROP COLUMN IF EXISTS github_delivery_id;