CREATE TABLE "tasks_payments" (
    "id" BIGSERIAL PRIMARY KEY,
    "task_id" BIGINT NOT NULL REFERENCES tasks(id),
    "amount" NUMERIC NOT NULL,
    "payee" VARCHAR NOT NULL REFERENCES payees(payee_name),
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "created_by" VARCHAR NOT NULL,
    "updated_on" timestamptz NOT NULL DEFAULT (now()),
    "updated_by" VARCHAR NOT NULL
);
ALTER TABLE "tasks" ADD COLUMN "paid_amount" NUMERIC DEFAULT (0);