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
        nick_name VARCHAR(255) DEFAULT NULL UNIQUE,
        avatar_url text DEFAULT NULL,
        profile_url text DEFAULT NULL,
        role VARCHAR(20) DEFAULT 'user',
        resigned BOOLEAN DEFAULT FALSE, 
        gender BOOLEAN DEFAULT NULL, 
        phone VARCHAR(255) DEFAULT NULL,
        address VARCHAR(255) DEFAULT NULL,
        position VARCHAR(255) DEFAULT NULL,
        date_of_birth DATE DEFAULT NULL,
        password VARCHAR(255) NOT NULL,
        settings JSONB NOT NULL DEFAULT '{}'::JSONB,
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now()),
        CONSTRAINT pk_users PRIMARY KEY (id)
    );
CREATE TABLE 
    auth.api_keys (
        id BIGSERIAL PRIMARY KEY,
        key VARCHAR(255) NOT NULL UNIQUE,
        status BOOLEAN NOT NULL,
        permissions TEXT [],
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now())
);
CREATE TABLE 
    auth.keys (
        id BIGSERIAL PRIMARY KEY,
        user_id uuid NOT NULL,
        public_key VARCHAR(255) NOT NULL,
        private_key VARCHAR(255) NOT NULL,
        refresh_token TEXT DEFAULT NULL,
        refresh_tokens_used TEXT [],
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now()),
        FOREIGN KEY (user_id) REFERENCES auth.users (id) ON DELETE CASCADE
);
CREATE INDEX ix_auth_key_token ON auth.keys (user_id);

-- Trigger event INSERT and UPDATE and DELETE
CREATE OR REPLACE FUNCTION notify_user_change()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        PERFORM pg_notify('user_change_event', 'INSERT');
    ELSIF TG_OP = 'UPDATE' THEN
        PERFORM pg_notify('user_change_event', 'UPDATE');
    ELSIF TG_OP = 'DELETE' THEN
        PERFORM pg_notify('user_change_event','DELETE');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_change_trigger
AFTER INSERT OR UPDATE OR DELETE ON auth.users
FOR EACH ROW
EXECUTE FUNCTION notify_user_change();

COMMIT;
