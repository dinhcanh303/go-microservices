-- name: Get :one
SELECT * FROM post.posts WHERE id = $1;

-- name: Create :one
INSERT INTO
    post.posts (
        id,
        status,
        content,
        bg_content,
        user_id,
        group_id
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: Update :one
UPDATE post.posts 
SET
    content = COALESCE(sqlc.narg(content),content),
    bg_content = COALESCE(sqlc.narg(bg_content),bg_content),
    status = COALESCE(sqlc.narg(status),status)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: GetByGroupId :many
SELECT * FROM post.posts WHERE group_id = ANY($1::uuid[]) ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: GetByUserId :many
SELECT * FROM post.posts 
WHERE user_id = $1 AND group_id IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: GetByFeed :many
SELECT *
FROM post.posts
WHERE user_id = ANY($1::uuid[])
   OR group_id = ANY($2::uuid[])
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;
-- name: Delete :exec
DELETE FROM post.posts WHERE id = $1;
