BEGIN;

-- Add new columns with the new data type
ALTER TABLE appeals ADD COLUMN picture_new oid;
ALTER TABLE appeals ADD COLUMN video_new oid;

-- Migrate the data
-- This step depends on how you want to convert the bytea data to oid.
-- You might need to write a function to perform this conversion.

-- Drop the old columns
ALTER TABLE appeals DROP COLUMN picture;
ALTER TABLE appeals DROP COLUMN video;

-- Rename the new columns to the old names
ALTER TABLE appeals RENAME COLUMN picture_new TO picture;
ALTER TABLE appeals RENAME COLUMN video_new TO video;

COMMIT;
