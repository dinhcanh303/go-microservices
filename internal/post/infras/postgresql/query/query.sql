-- name: Get :one
SELECT * FROM post.posts WHERE id = $1;

-- name: Create :one
INSERT INTO
    post.posts (
        id,
        status,
        title,
        content,
        user_id,
        group_id,
        deleted_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: Update :one
UPDATE post.posts 
SET
    title = $2 ,
    content = $3,
    status = $4,
    deleted_at = $5
WHERE id = $1 RETURNING *;

-- name: GetWithUnscoped :one
SELECT * FROM post.posts 
WHERE id = $1 AND deleted_at IS NOT NULL;

-- name: Delete :exec
DELETE FROM post.posts WHERE id = $1;

-- name: List :many
SELECT * FROM post.posts OFFSET $1 LIMIT $2;

-- name: Count :one
SELECT COUNT(*) FROM post.posts;
