BEGIN;

ALTER TABLE photo_gallery RENAME COLUMN edited_at TO updated_at;

COMMIT;