ALTER TABLE projects DROP COLUMN IF EXISTS company_id;
ALTER TABLE users DROP COLUMN IF EXISTS company_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS company_id;

DROP FUNCTION IF EXISTS "put_company_id()";
DROP TRIGGER IF EXISTS "put_company_id" ON "tasks";