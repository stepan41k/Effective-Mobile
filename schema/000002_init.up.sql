CREATE TYPE gen AS ENUM ('male', 'female', 'other');

ALTER TABLE profiles
ALTER COLUMN gender TYPE gen USING gender::gen;