-- BEGIN TRANSACTION;
BEGIN TRANSACTION;

-- delete large objects referenced by appeals table picture oid and video oid columns
-- do it.
DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN SELECT picture, video FROM appeals WHERE picture IS NOT NULL OR video IS NOT NULL LOOP
        IF r.picture IS NOT NULL THEN
            PERFORM lo_unlink(r.picture);
        END IF;
        IF r.video IS NOT NULL THEN
            PERFORM lo_unlink(r.video);
        END IF;
    END LOOP;
END;
$$;

-- drop picture and video columns from appeals table
ALTER TABLE appeals DROP COLUMN picture;
ALTER TABLE appeals DROP COLUMN video;

-- add picture bytea and video bytea columns to appeals table
ALTER TABLE appeals ADD COLUMN picture BYTEA;
ALTER TABLE appeals ADD COLUMN video BYTEA;

-- end transaction
COMMIT;