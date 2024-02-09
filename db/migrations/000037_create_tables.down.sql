BEGIN;

-- Delete article_comments table if it exists
DROP TABLE IF EXISTS article_comments;

-- Delete e_newspaper_comments table if it exists
DROP TABLE IF EXISTS e_newspaper_comments;

-- Delete news_post_comments table if it exists
DROP TABLE IF EXISTS news_post_comments;

-- Delete video_news_comments table if it exists
DROP TABLE IF EXISTS video_news_comments;

COMMIT;
