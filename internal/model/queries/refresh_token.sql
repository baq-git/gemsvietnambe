-- name: SaveRefreshToken :one
INSERT INTO refresh_tokens(
        id,
        refresh_token,
        user_id,
        revoked
    )
VALUES (
        gen_random_uuid(),
        $1,
        $2,
        FALSE
    )
RETURNING *;
-- name: RetrieveRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE user_id = $1;
-- name: RotationRefreshToken :exec
DELETE FROM refresh_tokens
WHERE id = $1;
-- name: DeleteRefreshTokenByUserID :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;