ALTER TABLE repositories
    DROP CONSTRAINT IF EXISTS repositories_user_id_github_repo_id_key;

ALTER TABLE repositories
    ADD CONSTRAINT repositories_github_repo_id_key UNIQUE (github_repo_id);

ALTER TABLE repositories
    DROP COLUMN IF EXISTS owner,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS clone_url,
    DROP COLUMN IF EXISTS webhook_secret;

ALTER TABLE repositories
    ALTER COLUMN name TYPE VARCHAR(255),
    ALTER COLUMN full_name TYPE VARCHAR(255),
    ALTER COLUMN default_branch TYPE VARCHAR(100),
    ALTER COLUMN private SET DEFAULT FALSE,
    ALTER COLUMN is_active SET DEFAULT FALSE;

ALTER TABLE repositories
    RENAME COLUMN user_id TO owner_id;

ALTER TABLE repositories
    RENAME COLUMN html_url TO github_url;

ALTER TABLE repositories
    RENAME COLUMN is_active TO webhook_active;

DROP INDEX IF EXISTS idx_repositories_user;
DROP INDEX IF EXISTS idx_repositories_github_repo_id;
DROP INDEX IF EXISTS idx_repositories_full_name;
CREATE INDEX idx_repositories_owner ON repositories(owner_id);