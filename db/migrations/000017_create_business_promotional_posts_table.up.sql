CREATE TABLE IF NOT EXISTS business_promotional_posts (
    id serial PRIMARY KEY,
    title_latin text NOT NULL UNIQUE,
    description_latin text,
    title_cyrillic text NOT NULL UNIQUE,
    description_cyrillic text,
    photos bytea[],
    videos bytea[],
    cover_image bytea,
    expiration timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 day',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
