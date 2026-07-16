ALTER TABLE repositories
    RENAME COLUMN owner_id TO user_id;

ALTER TABLE repositories
    RENAME COLUMN github_url TO html_url;

ALTER TABLE repositories
    RENAME COLUMN webhook_active TO is_active;

ALTER TABLE repositories
    ALTER COLUMN name TYPE TEXT,
    ALTER COLUMN full_name TYPE TEXT,
    ALTER COLUMN default_branch TYPE TEXT,
    ALTER COLUMN html_url TYPE TEXT;

ALTER TABLE repositories
    ADD COLUMN owner TEXT NOT NULL DEFAULT '', 
    ADD COLUMN description TEXT,
    ADD COLUMN clone_url TEXT NOT NULL DEFAULT '',
    ADD COLUMN webhook_secret TEXT,
    ALTER COLUMN private DROP DEFAULT,
    ALTER COLUMN is_active SET DEFAULT FALSE;

ALTER TABLE repositories 
    DROP CONSTRAINT IF EXISTS repositories_github_repo_id_key;

ALTER TABLE repositories 
    ADD CONSTRAINT repositories_user_id_github_repo_id_key UNIQUE (user_id, github_repo_id);

DROP INDEX IF EXISTS idx_repositories_owner;
CREATE INDEX idx_repositories_user 
ON repositories(user_id);
CREATE INDEX idx_repositories_github_repo_id
ON repositories(github_repo_id);

CREATE INDEX idx_repositories_full_name
ON repositories(full_name);