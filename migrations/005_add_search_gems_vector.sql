-- +goose Up
ALTER TABLE gems
ADD search_gems_vector tsvector GENERATED ALWAYS AS (
        setweight(
            to_tsvector('simple', coalesce(gem_name, '')),
            'A'
        ) || setweight(
            to_tsvector('simple', coalesce(instruction, '')),
            'B'
        ) || setweight(
            to_tsvector('simple', coalesce(description, '')),
            'C'
        )
    ) STORED;
CREATE INDEX IF NOT EXISTS gems_search_idx ON gems USING GIN(search_gems_vector);
-- +goose Down
DROP INDEX IF EXISTS gems_search_idx;
DROP COLUMN search_gems_vector;