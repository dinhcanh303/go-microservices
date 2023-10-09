-- name: Create :one
INSERT INTO 
    comment.comments (
        id,
		user_id,
		content,
		post_id,
		parent_comment_id,
		reply_to
    )
VALUES ($1, $2, $3, $4 ,$5 , $6) RETURNING *;

-- name: Get :one
SELECT * FROM comment.comments WHERE id = $1;

-- name: Update :one
UPDATE comment.comments 
SET
    content = $2 ,
    reply_to = $3
WHERE id = $1 RETURNING *;

-- name: Delete :exec
DELETE FROM comment.comments WHERE id = $1;

-- name: DeleteByCommentID :exec
DELETE FROM comment.comments WHERE parent_comment_id = $1;

-- name: ListByPostID :many
SELECT * FROM comment.comments WHERE post_id = $1;

-- name: CountByPostID :one
SELECT count(*) FROM comment.comments WHERE post_id = $1;

-- name: CountByCommentID :one
SELECT count(*) FROM comment.comments WHERE parent_comment_id = $1;