CREATE TABLE "accounts" (
    "id" SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,

    "username" VARCHAR(64) NOT NULL UNIQUE,
    "email" VARCHAR(320) NOT NULL UNIQUE,
    "password" BYTEA NOT NULL
);
