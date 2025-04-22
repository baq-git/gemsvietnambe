-- name: GetAllGems :many
SELECT *
FROM gems;
-- name: CreateGem :one
INSERT INTO gems (
        id,
        gem_name,
        description,
        instruction,
        coordinates,
        gem_category_id,
        created_at,
        updated_at
    )
VALUES (
        gen_random_uuid(),
        $1,
        $2,
        $3,
        $4,
        $5,
        NOW(),
        NOW()
    )
RETURNING *;
-- name: DeleteGem :exec
DELETE FROM gems
WHERE id = $1;
-- name: UpdateGem :one
UPDATE gems
SET gem_name = $2,
    description = $3,
    instruction = $4,
    coordinates = $5,
    gem_category_id = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
-- name: GetGem :one
SELECT *
FROM gems
    INNER JOIN gem_categories ON gems.gem_category_id = gem_categories.id
WHERE gems.id = $1;