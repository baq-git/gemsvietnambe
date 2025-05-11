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
-- name: ListingGemsComplex :many
SELECT g.id AS gem_id,
    g.gem_name,
    g.instruction,
    g.description,
    g.created_at,
    g.updated_at,
    g.gem_category_id,
    ts_rank(g.search_gems_vector, query) AS rank,
    gc.category_name,
    gc.description AS category_description,
    gc.slug AS category_slug,
    g.search_gems_vector,
    query
FROM gems g
    CROSS JOIN plainto_tsquery('simple', $1) query
    INNER JOIN gem_categories gc ON g.gem_category_id = gc.id
WHERE (
        (
            $1 = ''
            OR g.search_gems_vector @@ query
        )
        AND (
            gc.slug = $2
            OR $2 = ''
            OR $2 = 'all'
        )
    )
ORDER BY rank DESC
LIMIT $3 OFFSET $4;