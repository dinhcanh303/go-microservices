START TRANSACTION;
CREATE SCHEMA IF NOT EXISTS "noti";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    noti.notifications (
        id SERIAL PRIMARY KEY,
        actor_id uuid [] DEFAULT NULL,
        receiver_id uuid [] DEFAULT NULL,
        read_id uuid [] DEFAULT NULL,
        object_type VARCHAR(255) NOT NULL,
        object_id uuid NOT NULL,
        entity_type_id integer DEFAULT NULL,
        entity VARCHAR(255) DEFAULT NULL,
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now())
    );
-- CREATE TABLE 
--     noti.entity_types (
--         id SERIAL PRIMARY KEY,
--         name VARCHAR(255) NOT NULL,
--         description TEXT DEFAULT NULL
--     );
CREATE INDEX ix_receiver_id ON noti.notifications (receiver_id);
COMMIT;
