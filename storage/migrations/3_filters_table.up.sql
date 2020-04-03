BEGIN;

ALTER TABLE filters ADD COLUMN id SERIAL PRIMARY KEY;
ALTER TABLE filters ALTER COLUMN user_id TYPE UUID;
ALTER TABLE filters ADD CONSTRAINT unique_user_filtername UNIQUE(user_id, name);

COMMIT;