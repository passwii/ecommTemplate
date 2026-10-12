-- +goose Up
CREATE TABLE storefront_metadata (
    id SMALLINT PRIMARY KEY CHECK (id = 1),
    template_id TEXT NOT NULL CHECK (template_id <> ''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO storefront_metadata (id, template_id)
VALUES (1, 'gas-grill');

-- +goose Down
DROP TABLE storefront_metadata;
