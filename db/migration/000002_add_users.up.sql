CREATE TABLE "users" (
  "username" varchar UNIQUE PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),  
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "tasks" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("username");
ALTER TABLE "tasks" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("username");