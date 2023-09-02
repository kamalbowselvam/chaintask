ALTER TABLE IF EXISTS "tasks" DROP CONSTRAINT IF EXISTS "tasks_project_id_fkey";
ALTER TABLE IF EXISTS "tasks" DROP COLUMN IF EXISTS "project_id";


DROP FUNCTION IF EXISTS "check_roles()";
DROP TRIGGER IF EXISTS "check_roles" ON "projects";
DROP TABLE if EXISTS "projects";