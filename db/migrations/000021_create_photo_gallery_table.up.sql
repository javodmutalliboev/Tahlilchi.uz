-- Creating the 'photo_gallery' table
CREATE TABLE IF NOT EXISTS photo_gallery (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Creating the 'photo_gallery_photos' table
CREATE TABLE IF NOT EXISTS photo_gallery_photos (
    id SERIAL PRIMARY KEY,
    photo_gallery INTEGER NOT NULL REFERENCES photo_gallery(id),
    file_path TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
