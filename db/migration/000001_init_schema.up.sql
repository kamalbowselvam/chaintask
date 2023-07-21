CREATE TABLE "tasks" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "budget" FLOAT NOT NULL,
  "created_on" timestamptz NOT NULL DEFAULT (now()),
  "created_by" VARCHAR NOT NULL,
  "updated_on" timestamptz NOT NULL DEFAULT (now()),
  "updated_by" VARCHAR NOT NULL,
  "done" BOOLEAN NOT NULL DEFAULT (FALSE)
);