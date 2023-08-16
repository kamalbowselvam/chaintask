CREATE TABLE "roles" (
    id SERIAL PRIMARY KEY NOT NULL,
    userRole TEXT
);

INSERT INTO roles (userRole) VALUES 
    ('USER'),
    ('WORKS_MANAGER'),
    ('ADMIN');

ALTER TABLE "users" ADD "role_id" int REFERENCES roles(id);