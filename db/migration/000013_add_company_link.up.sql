DELETE FROM projects;
DELETE FROM tasks;
DELETE FROM users;
ALTER TABLE projects ADD COLUMN company_id BIGINT REFERENCES company(id) NOT NULL;
ALTER TABLE users ADD COLUMN company_id BIGINT REFERENCES company(id) NOT NULL;
INSERT INTO company VALUES (0, 'default_company', 'somewhere in Toulouse', NOW()::timestamp, 'kselvamADMIN', NOW()::timestamp, 'kselvamADMIN');
INSERT INTO users VALUES ('kselvamADMIN', '$2a$10$6nlF3jgvvVTU.uCpA2MFreu7Z5./IT5S2Rgr12KlSKuVHvvd4DRum', 'ADMIN ADMIN', 'admin@chaintask.org', 'ADMIN', NOW()::timestamp,NOW()::timestamp, 0, 0) on conflict do nothing;