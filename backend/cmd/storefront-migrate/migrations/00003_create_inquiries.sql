-- +goose Up
CREATE TABLE inquiries (
    id TEXT PRIMARY KEY CHECK (id <> ''),
    request_id TEXT UNIQUE,
    contact_name TEXT NOT NULL CHECK (contact_name <> ''),
    contact_email TEXT NOT NULL CHECK (contact_email <> ''),
    contact_company TEXT NOT NULL DEFAULT '',
    contact_phone TEXT NOT NULL DEFAULT '',
    destination_country TEXT NOT NULL DEFAULT '',
    product_id TEXT REFERENCES products (id),
    product_slug TEXT NOT NULL DEFAULT '',
    product_model TEXT NOT NULL DEFAULT '',
    product_name TEXT NOT NULL DEFAULT '',
    quantity INTEGER CHECK (quantity IS NULL OR quantity > 0),
    customization_request TEXT NOT NULL DEFAULT '',
    privacy_notice_version TEXT NOT NULL CHECK (privacy_notice_version <> ''),
    privacy_consented_at TIMESTAMPTZ NOT NULL,
    source_path TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'in_progress', 'quoted', 'won', 'lost', 'on_hold')),
    owner TEXT NOT NULL DEFAULT '',
    last_contacted_at TIMESTAMPTZ,
    next_action TEXT NOT NULL DEFAULT '',
    next_action_due_at TIMESTAMPTZ,
    loss_reason TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX inquiries_status_created_idx ON inquiries (status, created_at DESC, id);
CREATE INDEX inquiries_product_idx ON inquiries (product_id);

CREATE TABLE inquiry_activities (
    id BIGSERIAL PRIMARY KEY,
    inquiry_id TEXT NOT NULL REFERENCES inquiries (id) ON DELETE CASCADE,
    activity_type TEXT NOT NULL CHECK (activity_type IN ('created', 'note', 'contact', 'status_change')),
    body TEXT NOT NULL DEFAULT '',
    actor TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX inquiry_activities_inquiry_idx ON inquiry_activities (inquiry_id, created_at, id);

-- +goose Down
DROP TABLE inquiry_activities;
DROP TABLE inquiries;
