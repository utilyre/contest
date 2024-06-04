-- name: GetAccount :one
SELECT * FROM "accounts" WHERE "email" = $1 LIMIT 1;
