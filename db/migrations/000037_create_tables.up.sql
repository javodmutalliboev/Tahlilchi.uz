BEGIN;

-- Create article_comments table
CREATE TABLE IF NOT EXISTS article_comments (
    id SERIAL PRIMARY KEY,
    article INTEGER NOT NULL REFERENCES articles (id),
    text TEXT NOT NULL,
    contact TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create e_newspaper_comments table
CREATE TABLE IF NOT EXISTS e_newspaper_comments (
    id SERIAL PRIMARY KEY,
    e_newspaper INTEGER NOT NULL REFERENCES e_newspapers (id),
    text TEXT NOT NULL,
    contact TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create news_post_comments table
CREATE TABLE IF NOT EXISTS news_post_comments (
    id SERIAL PRIMARY KEY,
    news_post INTEGER NOT NULL REFERENCES news_posts (id),
    text TEXT NOT NULL,
    contact TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create video_news_comments table
CREATE TABLE IF NOT EXISTS video_news_comments (
    id SERIAL PRIMARY KEY,
    video_news INTEGER NOT NULL REFERENCES video_news (id),
    text TEXT NOT NULL,
    contact TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved BOOLEAN NOT NULL DEFAULT FALSE
);

COMMIT;
