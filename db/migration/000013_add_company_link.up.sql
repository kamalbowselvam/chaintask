DELETE FROM projects;
DELETE FROM tasks;
DELETE FROM users;
ALTER TABLE projects ADD COLUMN company_id BIGINT REFERENCES company(id) NOT NULL;
ALTER TABLE users ADD COLUMN company_id BIGINT REFERENCES company(id) NOT NULL;
ALTER TABLE tasks ADD COLUMN company_id BIGINT REFERENCES company(id) NOT NULL;

INSERT INTO company VALUES (0, 'default_company', 'somewhere in Toulouse', NOW()::timestamp, 'kselvamADMIN', NOW()::timestamp, 'kselvamADMIN');
INSERT INTO users VALUES ('kselvamADMIN', '$2a$10$6nlF3jgvvVTU.uCpA2MFreu7Z5./IT5S2Rgr12KlSKuVHvvd4DRum', 'ADMIN ADMIN', 'admin@chaintask.org', 'ADMIN', NOW()::timestamp,NOW()::timestamp, 0, 0) on conflict do nothing;

CREATE OR REPLACE FUNCTION put_company_id_for_tasks() RETURNS trigger AS $put_company_id_for_tasks$
  DECLARE
    company_id_var projects.company_id%TYPE;
  BEGIN
    SELECT company_id INTO company_id_var FROM projects where id=NEW.project_id;
    NEW.company_id = company_id_var;
    RETURN NEW;
  END;

$put_company_id_for_tasks$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER put_company_id_for_tasks BEFORE INSERT ON tasks
    FOR EACH ROW EXECUTE FUNCTION put_company_id_for_tasks();

CREATE OR REPLACE FUNCTION put_company_id_for_projects() RETURNS trigger AS $put_company_id_for_projects$
  DECLARE
    company_id_var projects.company_id%TYPE;
  BEGIN
    SELECT company_id INTO company_id_var FROM users where username=NEW.created_by;
    NEW.company_id = company_id_var;
    RETURN NEW;
  END;

$put_company_id_for_projects$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER put_company_id_for_projects BEFORE INSERT ON projects
    FOR EACH ROW EXECUTE FUNCTION put_company_id_for_projects();