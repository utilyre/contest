-- name: GetAccount :one
SELECT * FROM "accounts" WHERE "email" = $1 LIMIT 1;

-- name: CreateAccount :exec
INSERT INTO "accounts"
("created_at", "username", "email", "password")
VALUES (NOW(), $1, $2, $3)
RETURNING "id";

-- name: GetAccountEmailPassword :one
SELECT "email", "password"
FROM "accounts"
WHERE "username" = $1
LIMIT 1;
