CREATE TABLE payees (
    "id" BIGSERIAL PRIMARY KEY,
    "payee_name" varchar NOT NULL UNIQUE,
    "address" varchar NOT NULL ,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "created_by" VARCHAR NOT NULL,
    "updated_on" timestamptz NOT NULL DEFAULT (now()),
    "updated_by" VARCHAR NOT NULL
)

