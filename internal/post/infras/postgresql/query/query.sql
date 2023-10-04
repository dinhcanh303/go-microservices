-- name: GetPost :one
SELECT * FROM post.posts WHERE id = $1;

-- name: CreatePost :one
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

-- name: UpdatePost :one
UPDATE post.posts 
SET
    title = $2 ,
    content = $3,
    status = $4,
    deleted_at = $5
WHERE id = $1 RETURNING *;

-- name: UpdatePostWithUnscoped :one
UPDATE post.posts 
SET 
    title = $2, 
    content = $3, 
    status = $4
WHERE id = $1 AND deleted_at IS NOT NULL RETURNING *;

-- name: DeletePost :exec
DELETE FROM post.posts WHERE id = $1;

-- name: ListPosts :many
SELECT * FROM post.posts OFFSET $1 LIMIT $2;

-- name: CountPosts :one
SELECT COUNT(*) FROM post.posts;

