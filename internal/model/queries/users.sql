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