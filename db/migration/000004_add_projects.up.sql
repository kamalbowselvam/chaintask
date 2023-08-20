CREATE TABLE "projects" (
  "id" BIGSERIAL PRIMARY KEY,
  "pojectname" VARCHAR NOT NULL,
  "created_on" timestamptz NOT NULL DEFAULT (now()),
  "created_by" VARCHAR NOT NULL,
  "location" Point NOT NULL,
  "address" VARCHAR NOT NULL,
  "responsible" VARCHAR NOT NULL,
  "client" VARCHAR NOT NULL
);

ALTER TABLE "projects" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("username");
ALTER TABLE "projects" ADD FOREIGN KEY ("responsible") REFERENCES "users" ("username");

ALTER TABLE "tasks" ADD "project_id" INT NOT NULL REFERENCES "projects"("id"); 

CREATE OR REPLACE FUNCTION check_roles() RETURNS trigger AS $check_roles$
  DECLARE
    responsible_role users.role_id%TYPE;
    client_role      users.role_id%TYPE;
  BEGIN
    SELECT role_id INTO responsible_role FROM users where username=NEW.responsible;
    IF responsible_role != 2 THEN 
      RAISE EXCEPTION "Responsible has not the right role"
    END IF;
    SELECT role_id INTO client_role FROM users where username=NEW.client;
    IF client_role != 1 THEN
      RAISE EXCEPTION "Client has not the right role"
    END IF;
    RETURNS NEW;
  END;

$check_roles$ LANGUAGE plpgsql;

CREATE TRIGGER check_roles BEFORE INSERT OR UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION check_roles();