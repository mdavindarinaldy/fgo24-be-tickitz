ALTER TABLE users
ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN phone_number VARCHAR(255) NOT NULL DEFAULT '' UNIQUE;

UPDATE users
SET name = profiles.name,
phone_number = profiles.phone_number
FROM profiles
WHERE profiles.user_id = users.id;

DROP TABLE profiles;

ALTER TABLE users
ALTER COLUMN name DROP DEFAULT,
ALTER COLUMN phone_number DROP DEFAULT;