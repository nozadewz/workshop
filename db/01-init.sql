CREATE SEQUENCE IF NOT EXISTS account_id;
CREATE SEQUENCE IF NOT EXISTS cloud_pocket_id;

CREATE TABLE "accounts" (
    "id" int4 NOT NULL DEFAULT nextval('account_id'::regclass),
    "balance" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE TABLE "pockets" (
    "id" int4 NOT NULL DEFAULT nextval('cloud_pocket_id'::regclass),
    "name" TEXT NOT NULL,
    "category" TEXT,
    "currency" TEXT NOT NULL,
    "balance" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);