DELETE FROM projects;
DELETE FROM tasks;
DELETE FROM users;
ALTER TABLE projects ADD COLUMN company_id BIGINT REFERENCES company(id) NOT NULL;
ALTER TABLE users ADD COLUMN company_id BIGINT REFERENCES company(id) NOT NULL;
ALTER TABLE tasks ADD COLUMN company_id BIGINT REFERENCES company(id) NOT NULL;

INSERT INTO company VALUES (0, 'default_company', 'somewhere in Toulouse', NOW()::timestamp, 'kselvamADMIN', NOW()::timestamp, 'kselvamADMIN');
INSERT INTO users VALUES ('kselvamADMIN', '$2a$10$6nlF3jgvvVTU.uCpA2MFreu7Z5./IT5S2Rgr12KlSKuVHvvd4DRum', 'ADMIN ADMIN', 'admin@chaintask.org', 'ADMIN', NOW()::timestamp,NOW()::timestamp, 0, 0) on conflict do nothing;

CREATE OR REPLACE FUNCTION put_company_id() RETURNS trigger AS $put_company_id$
  DECLARE
    company_id projects.company_id%TYPE;
  BEGIN
    SELECT company_id INTO company_id FROM projects where id=NEW.project_id;
    NEW.company_id = company_id;
    RETURN NEW;
  END;

$put_company_id$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER put_company_id BEFORE INSERT ON tasks
    FOR EACH ROW EXECUTE FUNCTION put_company_id();