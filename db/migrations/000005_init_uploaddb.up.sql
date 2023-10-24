START TRANSACTION;
CREATE SCHEMA IF NOT EXISTS "upload";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    upload.attachments (
        id uuid NOT NULL DEFAULT (uuid_generate_v4()),
        attachable_type VARCHAR(255) DEFAULT NULL,
        attachable_id uuid DEFAULT NULL,
        user_id uuid NOT NULL,
        filename VARCHAR(255) NOT NULL,
        extension VARCHAR(255) NOT NULL,
        mime_type VARCHAR(255) DEFAULT NULL,
        folder VARCHAR(255) DEFAULT NULL,
        url_thumbnail VARCHAR(255) DEFAULT NULL,
        url VARCHAR(255) NOT NULL,
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now()),
        -- deleted_at timestamp with time zone NULL,
        CONSTRAINT pk_attachments PRIMARY KEY (id)
    );

CREATE INDEX ix_user_id ON upload.attachments (user_id);
COMMIT;
