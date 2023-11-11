CREATE TABLE company (
    "id" BIGSERIAL PRIMARY KEY,
    "companyname" varchar NOT NULL UNIQUE,
    "address" varchar,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "created_by" VARCHAR NOT NULL,
    "updated_on" timestamptz NOT NULL DEFAULT (now()),
    "updated_by" VARCHAR NOT NULL
)

