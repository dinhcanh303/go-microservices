START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS "group";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "group".groups (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    status integer NOT NULL DEFAULT 1,
    profile_url text DEFAULT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT (now()),
    -- deleted_at timestamp with time zone NULL,
    CONSTRAINT pk_groups PRIMARY KEY (id)
);

CREATE TABLE "group".group_members (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    group_id uuid NOT NULL,
    user_id uuid NOT NULL,
    role integer NOT NULL,
    status integer NOT NULL DEFAULT 1,
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT (now()),
    -- deleted_at timestamp with time zone NULL,
    CONSTRAINT pk_group_members PRIMARY KEY (id),
    FOREIGN KEY (group_id) REFERENCES "group".groups (id) ON DELETE CASCADE
);

CREATE INDEX ix_group_user_id ON "group".groups (user_id);
CREATE INDEX ix_group_member_user_id ON "group".group_members (user_id);
CREATE INDEX ix_group_member_group_id ON "group".group_members (group_id);
CREATE INDEX ix_group_member ON "group".group_members (group_id,user_id);
-- Trigger event INSERT and UPDATE and DELETE
CREATE OR REPLACE FUNCTION notify_group_change()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        PERFORM pg_notify('group_change_event', 'INSERT');
    ELSIF TG_OP = 'UPDATE' THEN
        PERFORM pg_notify('group_change_event', 'UPDATE');
    ELSIF TG_OP = 'DELETE' THEN
        PERFORM pg_notify('group_change_event','DELETE');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER group_change_trigger
AFTER INSERT OR UPDATE OR DELETE ON "group".groups
FOR EACH ROW
EXECUTE FUNCTION notify_group_change();
COMMIT;




