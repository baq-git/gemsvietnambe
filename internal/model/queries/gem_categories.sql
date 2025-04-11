-- name: GetAllGemCategories :many
SELECT *
FROM gem_categories;
-- name: GetGemCategory :one
SELECT *
FROM gem_categories
WHERE id = $1;