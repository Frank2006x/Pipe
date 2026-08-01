ALTER TABLE jobs
    DROP COLUMN IF EXISTS template_id;

DROP TABLE IF EXISTS repository_job_templates;
