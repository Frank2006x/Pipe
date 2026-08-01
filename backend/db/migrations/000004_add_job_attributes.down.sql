ALTER TABLE jobs
    DROP COLUMN IF EXISTS image,
    DROP COLUMN IF EXISTS working_directory,
    DROP COLUMN IF EXISTS commands,
    DROP COLUMN IF EXISTS logs,
    DROP COLUMN IF EXISTS exit_code;
