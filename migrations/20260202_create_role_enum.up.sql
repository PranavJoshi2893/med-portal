CREATE TYPE role AS ENUM('super_admin','admin','user');

ALTER TABLE users ADD COLUMN role role NOT NULL DEFAULT 'user';