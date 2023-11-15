CREATE TABLE resources (
    "id" BIGSERIAL PRIMARY KEY,
    "resource_name" varchar NOT NULL UNIQUE,
    "current" FLOAT DEFAULT (0) NOT NULL ,
    "availed" FLOAT NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "created_by" VARCHAR NOT NULL,
    "updated_on" timestamptz NOT NULL DEFAULT (now()),
    "updated_by" VARCHAR NOT NULL
)

