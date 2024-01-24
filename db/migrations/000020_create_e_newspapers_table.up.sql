CREATE TABLE e_newspapers (
    id SERIAL PRIMARY KEY,
    title_latin TEXT,
    title_cyrillic TEXT,
    file_latin BYTEA,
    file_cyrillic BYTEA,
    cover_image BYTEA,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    archived BOOLEAN NOT NULL DEFAULT FALSE
);