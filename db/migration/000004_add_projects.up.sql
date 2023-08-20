CREATE TABLE "projects" (
  "id" BIGSERIAL PRIMARY KEY,
  "projectname" VARCHAR NOT NULL,
  "created_on" timestamptz NOT NULL DEFAULT (now()),
  "created_by" VARCHAR NOT NULL,
  "location" Point NOT NULL,
  "address" VARCHAR NOT NULL,
  "responsible" VARCHAR NOT NULL,
  "client" VARCHAR NOT NULL
);

ALTER TABLE "projects" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("username");
ALTER TABLE "projects" ADD FOREIGN KEY ("responsible") REFERENCES "users" ("username");

DELETE from "tasks";

ALTER TABLE "tasks" ADD "project_id" INT NOT NULL REFERENCES "projects"("id"); 

CREATE OR REPLACE FUNCTION check_roles() RETURNS trigger AS $check_roles$
  DECLARE
    responsible_role users.role_id%TYPE;
    client_role      users.role_id%TYPE;
  BEGIN
    SELECT role_id INTO responsible_role FROM users where username=NEW.responsible;
    IF responsible_role != 2 THEN 
      RAISE EXCEPTION 'Responsible %s has not the right role', NEW.responsible;
    END IF;
    SELECT role_id INTO client_role FROM users where username=NEW.client;
    IF client_role != 1 THEN
      RAISE EXCEPTION 'Client %s has not the right role', NEW.client;
    END IF;
    RETURN NEW;
  END;

$check_roles$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER check_roles BEFORE INSERT OR UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION check_roles();