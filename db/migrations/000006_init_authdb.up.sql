START TRANSACTION;
CREATE SCHEMA IF NOT EXISTS "auth";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    auth.user (
        id uuid NOT NULL DEFAULT (uuid_generate_v4()),
        username VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password uuid NOT NULL,
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now()),
        CONSTRAINT pk_users PRIMARY KEY (id)
    );
CREATE TABLE 
    auth.api_keys (
        id bigint unsigned not null auto_increment PRIMARY KEY,
        key VARCHAR(255) NOT NULL UNIQUE,
        status BOOLEAN NOT NULL,
        permissions JSON DEFAULT '[]',
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now()),
);
CREATE TABLE 
    auth.key_tokens (
        id bigint unsigned not null auto_increment PRIMARY KEY,
        user_id uuid NOT NULL,
        public_key VARCHAR(255) NOT NULL,
        private_key VARCHAR(255) NOT NULL,
        refresh_token VARCHAR(255) NOT NULL,
        refresh_tokens_used JSON DEFAULT '[]',
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now()),
        FOREIGN KEY (user_id) REFERENCES auth.users (id)
);
CREATE INDEX ix_auth_key_token ON auth.key_tokens (user_id);
COMMIT;
