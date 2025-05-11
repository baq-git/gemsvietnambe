-- +goose Up
CREATE TABLE gems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    gem_category_id UUID NOT NULL REFERENCES gem_categories(id) ON DELETE RESTRICT,
    gem_name TEXT NOT NULL,
    description TEXT NOT NULL,
    instruction TEXT NOT NULL,
    coordinates DOUBLE PRECISION [] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose Down
DROP TABLE gems;