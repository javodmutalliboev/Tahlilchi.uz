BEGIN;

-- FILEPATH: /home/javod/Desktop/Projects/Tahlilchi.uz/applications/server/db/migrations/000033_alter_table.up.sql
ALTER TABLE news_posts
ADD COLUMN completed BOOLEAN NOT NULL DEFAULT FALSE;

-- FILEPATH: /home/javod/Desktop/Projects/Tahlilchi.uz/applications/server/db/migrations/000033_alter_table.up.sql
ALTER TABLE articles
ADD COLUMN completed BOOLEAN NOT NULL DEFAULT FALSE;

-- FILEPATH: /home/javod/Desktop/Projects/Tahlilchi.uz/applications/server/db/migrations/000033_alter_table.up.sql
ALTER TABLE business_promotional_posts
ADD COLUMN completed BOOLEAN NOT NULL DEFAULT FALSE;

-- FILEPATH: /home/javod/Desktop/Projects/Tahlilchi.uz/applications/server/db/migrations/000033_alter_table.up.sql
ALTER TABLE e_newspapers
ADD COLUMN completed BOOLEAN NOT NULL DEFAULT FALSE;

COMMIT;
