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

-- -- name: GetWithUnscoped :one
-- SELECT * FROM post.posts 
-- WHERE id = $1 AND deleted_at IS NOT NULL;


-- name: GetByGroupId :many
SELECT * FROM post.posts 
WHERE group_id = $1;

-- name: GetByUserId :many
SELECT * FROM post.posts 
WHERE user_id = $1 AND group_id IS NULL;

-- name: Delete :exec
DELETE FROM post.posts WHERE id = $1;

-- name: List :many
SELECT * FROM post.posts OFFSET $1 LIMIT $2;

-- -- name: CountPost :one
-- SELECT COUNT(*) FROM post.posts WHERE user_id = $1 AND group_id IS NULL;

-- -- name: CountPostInGroup :one
-- SELECT COUNT(*) FROM post.posts WHERE group_id = $1;
