BEGIN;

ALTER TABLE articles
    ALTER COLUMN videos TYPE text[] USING videos::text[],
    ADD COLUMN category integer,
    ADD COLUMN related integer;

ALTER TABLE articles
    RENAME COLUMN edited_at TO updated_at;

ALTER TABLE articles
    ADD CONSTRAINT fk_category
    FOREIGN KEY (category)
    REFERENCES article_category(id);

ALTER TABLE articles
    ADD CONSTRAINT fk_related
    FOREIGN KEY (related)
    REFERENCES articles(id);

COMMIT;
