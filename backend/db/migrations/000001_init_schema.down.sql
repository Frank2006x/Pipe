DROP INDEX IF EXISTS idx_logs_job;
DROP INDEX IF EXISTS idx_jobs_status;
DROP INDEX IF EXISTS idx_jobs_pipeline;
DROP INDEX IF EXISTS idx_pipelines_status;
DROP INDEX IF EXISTS idx_pipelines_repository;
DROP INDEX IF EXISTS idx_repositories_github_repo_id;
DROP INDEX IF EXISTS idx_repositories_owner;
DROP INDEX IF EXISTS idx_users_github_id;

DROP TABLE IF EXISTS logs;

DROP TABLE IF EXISTS jobs;

DROP TABLE IF EXISTS pipelines;

DROP TABLE IF EXISTS repositories;

DROP TABLE IF EXISTS github_tokens;

DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS github_event;

DROP TYPE IF EXISTS job_status;

DROP TYPE IF EXISTS pipeline_status;