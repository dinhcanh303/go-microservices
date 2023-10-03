-- name: GetPost
SELECT * FROM posts WHERE uuid = $1;

-- name: CreatePost :one
INSERT INTO
    post.posts (
        uuid,
        status,
        title,
        content,
        user_id,
        group_id,
        created_at,
        updated_at,
        deleted_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: UpdatePost :one
UPDATE post.posts 
SET
    title = $2 ,
    content = $3,
    status = $4,
    deleted_at = $5,
WHERE uuid = $1 RETURNING *;

-- name: UpdatePostWithUnscoped
UPDATE posts 
SET 
    title = $2, 
    content = $3, 
    status = $4, 
WHERE uuid = $1 AND deleted_at NOT NULL RETURNING *;

-- name: DeletePost
DELETE FROM posts WHERE uuid = $1;

-- name: ListPosts
SELECT * FROM posts OFFSET $1 LIMIT $2;

-- name: CountPosts
SELECT COUNT(*) FROM posts;

