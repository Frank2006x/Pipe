CREATE TABLE repositories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    github_url TEXT UNIQUE NOT NULL,
    default_branch VARCHAR(50) DEFAULT 'main',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE pipelines (
    id BIGSERIAL PRIMARY KEY,
    repository_id BIGINT NOT NULL REFERENCES repositories(id),
    commit_sha VARCHAR(255),
    branch VARCHAR(100),
    status VARCHAR(30) NOT NULL,
    trigger_type VARCHAR(30),
    created_at TIMESTAMP DEFAULT NOW(),
    started_at TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE TABLE jobs (
    id BIGSERIAL PRIMARY KEY,
    pipeline_id BIGINT NOT NULL REFERENCES pipelines(id),
    name VARCHAR(100) NOT NULL,
    status VARCHAR(30) NOT NULL,
    started_at TIMESTAMP,
    completed_at TIMESTAMP
);


CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,
    job_id BIGINT NOT NULL REFERENCES jobs(id),
    content TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);