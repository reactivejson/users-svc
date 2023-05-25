CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS users;

CREATE TABLE users (
                       id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                       first_name varchar(64) NOT NULL,
                       last_name varchar(64) NOT NULL,
                       nickname varchar(32) NOT NULL,
                       password varchar(256),
                       email varchar(128) NOT NULL UNIQUE,
                       country varchar(64) NOT NULL,
                       created_at timestamp with time zone NOT NULL DEFAULT NOW(),
                       updated_at timestamp with time zone
);

CREATE INDEX created_at_idx on users(created_at);
CREATE INDEX email_idx on users(email);