CREATE TABLE IF NOT EXISTS blogs (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    title        VARCHAR(255) NOT NULL,
    slug         VARCHAR(255) NOT NULL UNIQUE,
    content      TEXT         NOT NULL,
    excerpt      TEXT,
    status       VARCHAR(20)  NOT NULL DEFAULT 'draft',
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_blogs_slug ON blogs(slug);
CREATE INDEX idx_blogs_status_published ON blogs(status, published_at DESC);
