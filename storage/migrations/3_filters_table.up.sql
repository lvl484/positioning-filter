ALTER TABLE filters ADD COLUMN id SERIAL PRIMARY KEY;
ALTER TABLE filters ALTER COLUMN user_id TYPE UUID;
ALTER TABLE filters DROP CONSTRAINT filters_name_key;
ALTER TABLE filters ADD CONSTRAINT unique_user_filtername UNIQUE(user_id, name);