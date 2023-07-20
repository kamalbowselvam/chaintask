CREATE TABLE "tasks" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "budget" float,
  "created_at" timestamp,
  "created_by" varchar,
  "validated_on" timestamp,
  "validate_by" varchar,
  "updated_on" timestamp,
  "updated_by" varchar
);