ALTER TABLE photo_gallery_photos RENAME COLUMN file_path TO file_name;
ALTER TABLE photo_gallery_photos ADD COLUMN file bytea NOT NULL;
