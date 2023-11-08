START TRANSACTION;
CREATE SCHEMA IF NOT EXISTS "auth";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    auth.users (
        id uuid NOT NULL DEFAULT (uuid_generate_v4()),
        email VARCHAR(255) NOT NULL UNIQUE,
        first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255) NOT NULL,
        full_name VARCHAR(255) DEFAULT NULL,
        password VARCHAR(255) NOT NULL,
        roles VARCHAR(20) DEFAULT 'user', 
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now()),
        CONSTRAINT pk_users PRIMARY KEY (id)
    );
CREATE TABLE 
    auth.api_keys (
        id BIGSERIAL PRIMARY KEY,
        key VARCHAR(255) NOT NULL UNIQUE,
        status BOOLEAN NOT NULL,
        permissions JSON DEFAULT '[]',
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now())
);
CREATE TABLE 
    auth.keys (
        id BIGSERIAL PRIMARY KEY,
        user_id uuid NOT NULL,
        public_key VARCHAR(255) NOT NULL,
        private_key VARCHAR(255) NOT NULL,
        refresh_token VARCHAR(255) NOT NULL,
        refresh_tokens_used JSON DEFAULT '[]',
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now()),
        FOREIGN KEY (user_id) REFERENCES auth.users (id)
);
CREATE INDEX ix_auth_key_token ON auth.keys (user_id);
COMMIT;
