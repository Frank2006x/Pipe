CREATE TABLE repository_job_templates (
    id BIGSERIAL PRIMARY KEY,
    repository_id BIGINT NOT NULL
        REFERENCES repositories(id)
        ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    order_index INT NOT NULL DEFAULT 0,
    image TEXT NOT NULL DEFAULT '',
    working_directory TEXT NOT NULL DEFAULT '',
    commands TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_repo_job_templates_repo_order
ON repository_job_templates(repository_id, order_index);

ALTER TABLE jobs
    ADD COLUMN template_id BIGINT REFERENCES repository_job_templates(id) ON DELETE SET NULL;
