START TRANSACTION;
CREATE SCHEMA IF NOT EXISTS "comment";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    comment.comments (
        id uuid NOT NULL DEFAULT (uuid_generate_v4()),
        user_id uuid NOT NULL,
        content text NOT NULL,
        reply_id uuid DEFAULT NULL,
        tag_ids uuid [] DEFAULT NULL,
        post_id uuid NOT NULL,
        parent_comment_id uuid DEFAULT NULL,
        edited boolean DEFAULT NULL,
        created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
        updated_at timestamp with time zone NOT NULL DEFAULT (now()),
        -- deleted_at timestamp with time zone NULL,
        CONSTRAINT pk_comments PRIMARY KEY (id)
    );

CREATE INDEX ix_user_id ON comment.comments (user_id);
CREATE INDEX ix_post_id ON comment.comments (post_id);
CREATE INDEX ix_parent_comment_id ON comment.comments (parent_comment_id);
COMMIT;