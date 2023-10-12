START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS "group";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "group".groups (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    status integer NOT NULL,
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
    created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp with time zone NOT NULL DEFAULT (now()),
    -- deleted_at timestamp with time zone NULL,
    CONSTRAINT pk_group_members PRIMARY KEY (id)
);

CREATE INDEX ix_group_user_id ON "group".groups (user_id);
CREATE INDEX ix_group_member_user_id ON "group".group_members (user_id);
CREATE INDEX ix_group_member_group_id ON "group".group_members (group_id);

COMMIT;




