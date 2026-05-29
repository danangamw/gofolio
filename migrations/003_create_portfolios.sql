CREATE TABLE IF NOT EXISTS portfolios (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    title          VARCHAR(255) NOT NULL,
    description    TEXT,
    image_url      VARCHAR(500),
    tech_stack     TEXT[]       NOT NULL DEFAULT '{}',
    project_url    VARCHAR(500),
    repository_url VARCHAR(500),
    sort_order     INTEGER      NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
