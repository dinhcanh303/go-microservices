-- name: Create :one
INSERT INTO 
    comment.comments (
        id,
		user_id,
		content,
		post_id,
		parent_comment_id,
		reply_to_id
    )
VALUES ($1, $2, $3, $4 ,$5 , $6) RETURNING *;

-- name: Get :one
SELECT * FROM comment.comments WHERE id = $1;

-- name: Update :one
UPDATE comment.comments 
SET
    content = COALESCE(sqlc.narg(content),content),
    reply_to_id = COALESCE(sqlc.narg(reply_to_id),reply_to_id)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: Delete :exec
DELETE FROM comment.comments WHERE id = $1 OR parent_comment_id = $1;

-- name: DeleteAllByPostID :exec
DELETE FROM comment.comments WHERE post_id = $1;

-- name: GetCommentsByPostID :many
SELECT * FROM comment.comments WHERE post_id = $1;

-- name: CountByPostID :one
SELECT count(*) FROM comment.comments WHERE post_id = $1;

-- name: CountByCommentID :one
SELECT count(*) FROM comment.comments WHERE parent_comment_id = $1;