-- name: CreateUser :one
INSERT INTO users (
        id,
        email,
        username,
        password_hash,
        created_at,
        updated_at
    )
VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW())
RETURNING *;
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;
-- name: CheckUserById :one
SELECT *
FROM users
WHERE id = $1;