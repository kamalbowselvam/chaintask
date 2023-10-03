ALTER TABLE tasks DROP COLUMN if EXISTS version_id;
ALTER TABLE projects DROP COLUMN if EXISTS version_id;
ALTER TABLE users DROP COLUMN if EXISTS version_id;