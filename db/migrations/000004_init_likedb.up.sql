START TRANSACTION;
CREATE SCHEMA IF NOT EXISTS "like";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    "like".likes (
        id uuid NOT NULL DEFAULT (uuid_generate_v4()),
        emoji VARCHAR(255) NOT NULL,
        likeable_type VARCHAR(255) NOT NULL,
        likeable_id uuid NOT NULL,
        user_id uuid NOT NULL,
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now()),
        -- deleted_at timestamp with time zone NULL,
        CONSTRAINT pk_likes PRIMARY KEY (id)
    );

CREATE INDEX ix_user_id ON "like".likes (user_id);
COMMIT;