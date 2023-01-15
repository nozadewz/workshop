CREATE SEQUENCE IF NOT EXISTS account_id;

CREATE TABLE "accounts" (
    "id" INT NOT NULL DEFAULT nextval('account_id'::regclass),
    "balance" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE TABLE "pockets" (
    "id" int NOT NULL DEFAULT nextval('pockets_id'::regclass),
    "account_id" INT,
    "name" TEXT,
    "currency" TEXT,
    "balance" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
    CONSTRAINT fk_account_id FOREIGN KEY(account_id) REFERENCES accounts(id)
)
        