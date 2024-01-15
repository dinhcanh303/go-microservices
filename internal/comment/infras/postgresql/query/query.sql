-- name: Create :one
INSERT INTO 
    comment.comments (
        id,
		user_id,
		content,
		post_id,
		parent_comment_id,
		reply_id,
        tag_ids
    )
VALUES ( $1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: Get :one
SELECT * FROM comment.comments WHERE id = $1;

-- name: Update :one
UPDATE comment.comments 
SET
    content = COALESCE(sqlc.narg(content),content),
    reply_id = COALESCE(sqlc.narg(reply_id),reply_id),
    tag_ids = COALESCE(sqlc.narg(tag_ids),tag_ids)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: Delete :exec
DELETE FROM comment.comments WHERE id = $1 OR parent_comment_id = $1;

-- name: DeleteAllByPostID :exec
DELETE FROM comment.comments WHERE post_id = $1;

-- name: GetCommentsByPostID2 :many
SELECT * 
FROM comment.comments 
WHERE post_id = $1 AND parent_comment_id is not null
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
-- name: GetCommentsByPostID :many
WITH ranked_comments AS (
        SELECT
            c.*,
            ROW_NUMBER() OVER (PARTITION BY parent_comment_id ORDER BY created_at) AS row_num
        FROM
            comment.comments c
        WHERE
            c.parent_comment_id IN (
                SELECT id
                FROM comment.comments
                WHERE post_id = $1 AND parent_comment_id IS NULL
                LIMIT $2 OFFSET $3
            )
)
SELECT tb.*
FROM
(
    SELECT 
    tb1.id,
	tb1.user_id,
	tb1.content,
	tb1.reply_id,
    tb1.tag_ids,
	tb1.post_id,
	tb1.parent_comment_id,
	tb1.created_at,
	tb1.updated_at
	FROM (
		SELECT tb2.* FROM comment.comments AS tb2
		WHERE tb2.post_id = $1 AND tb2.parent_comment_id IS NULL
		LIMIT $2 OFFSET $3
	)  AS tb1
UNION
	SELECT 
	child.id,
    child.user_id,
    child.content,
    child.reply_id,
    child.tag_ids,
    child.post_id,
    child.parent_comment_id,
    child.created_at,
    child.updated_at
	FROM (
        
		SELECT *
		FROM ranked_comments
		WHERE row_num = 1
	) AS child
) AS tb ORDER BY created_at;

-- name: GetCommentsByCommentID :many
SELECT *
FROM comment.comments
WHERE parent_comment_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountByPostID :one
SELECT count(*) FROM comment.comments WHERE post_id = $1;

-- name: CountByCommentID :one
SELECT count(*) FROM comment.comments WHERE parent_comment_id = $1;