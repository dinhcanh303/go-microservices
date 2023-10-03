START TRANSACTION;
CREATE SCHEMA IF NOT EXISTS "post";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    post.posts (
        id uuid NOT NULL DEFAULT (uuid_generate_v4()),
        user_id uuid NOT NULL,
        group_id uuid NULL,
        title text NOT NULL,
        content text NOT NULL,
        status int NOT NULL,
        with
            time zone NOT NULL,
            created_at timestamp
        with
            time zone NOT NULL DEFAULT (now()),
            updated_at timestamp
        with
            time zone NULL,
            CONSTRAINT pk_posts PRIMARY KEY (id)
    );

CREATE INDEX ix_user_id ON post.posts (user_id);
CREATE INDEX ix_group_id ON post.posts (group_id);
COMMIT;