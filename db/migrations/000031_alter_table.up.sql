BEGIN;

ALTER TABLE e_newspapers RENAME COLUMN edited_at TO updated_at;

COMMIT;
