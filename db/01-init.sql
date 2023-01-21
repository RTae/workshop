CREATE SEQUENCE IF NOT EXISTS account_id;

CREATE TABLE IF NOT EXISTS TBL_Accounts (
    "id" int4 NOT NULL DEFAULT nextval('account_id'::regclass),
    "name" varchar(255) NOT NULL,
    PRIMARY KEY ("id")
);

CREATE SEQUENCE IF NOT EXISTS pocket_id;

CREATE TABLE IF NOT EXISTS TBL_Pockets (
    "id" int4 NOT NULL DEFAULT nextval('pocket_id'::regclass),
    "accountId" int4 REFERENCES "tbl_accounts" ("id"),
    "name" varchar(255) NOT NULL,
    "amount" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE SEQUENCE IF NOT EXISTS transaction_id;

CREATE TABLE IF NOT EXISTS TBL_Transactions (
    "id" int4 NOT NULL DEFAULT nextval('transaction_id'::regclass),
    "accountId" int4 REFERENCES "tbl_accounts" ("id"),
    "fromPocketId" int4 REFERENCES "tbl_pockets" ("id"),
    "toPocketId" int4 REFERENCES "tbl_pockets" ("id"),
    "amount" float8 NOT NULL,
    "date" timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);
