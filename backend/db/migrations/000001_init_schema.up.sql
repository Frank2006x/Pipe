-- ===========================
-- ENUMS
-- ===========================

CREATE TYPE pipeline_status AS ENUM (
    'pending',
    'running',
    'success',
    'failed',
    'cancelled'
);

CREATE TYPE job_status AS ENUM (
    'pending',
    'running',
    'success',
    'failed',
    'cancelled'
);

CREATE TYPE github_event AS ENUM (
    'push',
    'pull_request',
    'workflow_dispatch',
    'release',
    'tag'
);

-- ===========================
-- USERS
-- ===========================

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    github_id BIGINT NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    avatar_url TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===========================
-- GITHUB TOKENS
-- ===========================

CREATE TABLE github_tokens (
    user_id BIGINT PRIMARY KEY
        REFERENCES users(id)
        ON DELETE CASCADE,

    access_token TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===========================
-- REPOSITORIES
-- ===========================

CREATE TABLE repositories (
    id BIGSERIAL PRIMARY KEY,

    owner_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    github_repo_id BIGINT NOT NULL UNIQUE,

    name VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    github_url TEXT NOT NULL,

    default_branch VARCHAR(100) NOT NULL,

    private BOOLEAN NOT NULL DEFAULT FALSE,

    webhook_id BIGINT,
    webhook_active BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===========================
-- PIPELINES
-- ===========================

CREATE TABLE pipelines (
    id BIGSERIAL PRIMARY KEY,

    repository_id BIGINT NOT NULL
        REFERENCES repositories(id)
        ON DELETE CASCADE,

    commit_sha CHAR(40) NOT NULL,
    commit_message TEXT,

    branch VARCHAR(255) NOT NULL,

    event_type github_event NOT NULL,

    status pipeline_status NOT NULL DEFAULT 'pending',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

-- ===========================
-- JOBS
-- ===========================

CREATE TABLE jobs (
    id BIGSERIAL PRIMARY KEY,

    pipeline_id BIGINT NOT NULL
        REFERENCES pipelines(id)
        ON DELETE CASCADE,

    name VARCHAR(100) NOT NULL,

    status job_status NOT NULL DEFAULT 'pending',

    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

-- ===========================
-- LOGS
-- ===========================

CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,

    job_id BIGINT NOT NULL
        REFERENCES jobs(id)
        ON DELETE CASCADE,

    line_number INT NOT NULL,

    content TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ===========================
-- INDEXES
-- ===========================

CREATE INDEX idx_users_github_id
ON users(github_id);

CREATE INDEX idx_repositories_owner
ON repositories(owner_id);

CREATE INDEX idx_repositories_github_repo_id
ON repositories(github_repo_id);

CREATE INDEX idx_pipelines_repository
ON pipelines(repository_id);

CREATE INDEX idx_pipelines_status
ON pipelines(status);

CREATE INDEX idx_jobs_pipeline
ON jobs(pipeline_id);

CREATE INDEX idx_jobs_status
ON jobs(status);

CREATE INDEX idx_logs_job
ON logs(job_id);