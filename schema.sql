DROP DATABASE IF EXISTS socialnetwork CASCADE;
CREATE DATABASE IF NOT EXISTS socialnetwork;
SET DATABASE = socialnetwork;

CREATE TABLE IF NOT EXISTS users(
    id SERIAL NOT NULL PRIMARY KEY,
    email STRING NOT NULL UNIQUE,
    username STRING NOT NULL UNIQUE,
    avatar_url STRING,
    followers_count INT NOT NULL CHECK (followers_count>=0) DEFAULT 0,
    following_count INT NOT NULL CHECK (following_count>=0) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO users(id, email, username) VALUES
    (1,'john@example.DEV', 'john_doe'), 
    (2,'jane@example.DEV', 'jane_doe');