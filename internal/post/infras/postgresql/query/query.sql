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
        group_id
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: Update :one
UPDATE post.posts 
SET
    title = COALESCE(sqlc.narg(title),title),
    content = COALESCE(sqlc.narg(content),content),
    status = COALESCE(sqlc.narg(status),status)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: GetByGroupId :many
SELECT * FROM post.posts WHERE group_id = $1 LIMIT $2 OFFSET $3;

-- name: GetByUserId :many
SELECT * FROM post.posts 
WHERE user_id = $1 AND group_id IS NULL LIMIT $2 OFFSET $3;

-- name: GetByFeed :many
SELECT *
FROM post.posts
WHERE user_id IN (sqlc.slice(User_ids))
   OR group_id IN (sqlc.slice(Group_ids))
LIMIT $1 OFFSET $2;
-- name: Delete :exec
DELETE FROM post.posts WHERE id = $1;
