-- +goose Up
CREATE TABLE collections (
    id BIGSERIAL PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE CHECK (slug <> ''),
    name TEXT NOT NULL CHECK (name <> ''),
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published')),
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE products (
    id TEXT PRIMARY KEY CHECK (id <> ''),
    slug TEXT NOT NULL UNIQUE CHECK (slug <> ''),
    name TEXT NOT NULL CHECK (name <> ''),
    model TEXT NOT NULL DEFAULT '',
    sku TEXT NOT NULL DEFAULT '',
    collection_id BIGINT NOT NULL REFERENCES collections (id),
    short_description TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    price_amount NUMERIC(12, 2),
    price_currency TEXT NOT NULL DEFAULT '',
    compare_at_price NUMERIC(12, 2),
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published')),
    content_status TEXT NOT NULL DEFAULT 'sample' CHECK (content_status IN ('sample', 'verified')),
    featured BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    badge TEXT NOT NULL DEFAULT '',
    fuel TEXT NOT NULL DEFAULT '',
    burners INTEGER,
    cooking_area TEXT NOT NULL DEFAULT '',
    finish TEXT NOT NULL DEFAULT '',
    minimum_order_quantity TEXT NOT NULL DEFAULT '',
    lead_time TEXT NOT NULL DEFAULT '',
    applications JSONB NOT NULL DEFAULT '[]'::jsonb,
    images JSONB NOT NULL DEFAULT '[]'::jsonb,
    specifications JSONB NOT NULL DEFAULT '[]'::jsonb,
    resources JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT products_applications_is_array CHECK (jsonb_typeof(applications) = 'array'),
    CONSTRAINT products_images_is_array CHECK (jsonb_typeof(images) = 'array'),
    CONSTRAINT products_specifications_is_array CHECK (jsonb_typeof(specifications) = 'array'),
    CONSTRAINT products_resources_is_array CHECK (jsonb_typeof(resources) = 'array')
);

CREATE INDEX products_status_sort_order_idx ON products (status, sort_order, id);
CREATE INDEX products_collection_status_idx ON products (collection_id, status, sort_order);

CREATE TABLE company_content (
    id SMALLINT PRIMARY KEY CHECK (id = 1),
    about JSONB NOT NULL DEFAULT '{}'::jsonb,
    manufacturing JSONB NOT NULL DEFAULT '{}'::jsonb,
    credentials JSONB NOT NULL DEFAULT '[]'::jsonb,
    case_studies JSONB NOT NULL DEFAULT '[]'::jsonb,
    faqs JSONB NOT NULL DEFAULT '[]'::jsonb,
    content_status TEXT NOT NULL DEFAULT 'sample' CHECK (content_status IN ('sample', 'verified')),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE company_content;
DROP TABLE products;
DROP TABLE collections;
