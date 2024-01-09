-- name: Create :one
INSERT INTO 
    "like".likes (
        id,
		emoji,
		likeable_type,
		likeable_id,
		user_id
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: Update :one
UPDATE "like".likes
SET
	emoji = COALESCE(sqlc.narg(emoji),emoji)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: Delete :exec
DELETE FROM "like".likes WHERE id = $1;

-- name: GetAllByType :many
SELECT * FROM "like".likes WHERE likeable_type = $1
AND likeable_id = $2;

-- name: GetLikesInfoByType :one
SELECT
    COUNT(*) FILTER (WHERE user_id = $1) AS your_like,
    COUNT(*) FILTER (WHERE user_id != $1) AS others_likes,
	COALESCE(MIN(CASE WHEN user_id = $1 THEN emoji END),'') AS your_liked_emoji,
	COALESCE(ARRAY_AGG(DISTINCT CASE WHEN user_id != $1 THEN emoji END), '{}') AS others_liked_emojis
FROM "like".likes
WHERE likeable_type = $2 AND likeable_id = $3
LIMIT 1;

