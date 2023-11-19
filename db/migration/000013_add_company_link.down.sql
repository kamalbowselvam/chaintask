ALTER TABLE projects DROP COLUMN IF EXISTS company_id;
ALTER TABLE users DROP COLUMN IF EXISTS company_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS company_id;

DROP FUNCTION IF EXISTS "put_company_id_for_tasks()";
DROP TRIGGER IF EXISTS "put_company_id_for_tasks" ON "tasks";

DROP FUNCTION IF EXISTS "put_company_id_for_projects()";
DROP TRIGGER IF EXISTS "put_company_id_for_projects" ON "projects";