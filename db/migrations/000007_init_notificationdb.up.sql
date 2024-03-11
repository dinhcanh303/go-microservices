START TRANSACTION;
CREATE SCHEMA IF NOT EXISTS "noti";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    noti.notifications (
        id SERIAL PRIMARY KEY,
        actor_id uuid [] DEFAULT NULL,
        sender_id uuid [] DEFAULT NULL,
        data json DEFAULT NULL,
        type VARCHAR(255) DEFAULT NULL,
        object_type VARCHAR(255) DEFAULT NULL,
        object_id uuid DEFAULT NULL,
        read_at timestamp DEFAULT NULL,
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now())
    );

CREATE INDEX ix_sender_id ON noti.notifications (sender_id);
COMMIT;
