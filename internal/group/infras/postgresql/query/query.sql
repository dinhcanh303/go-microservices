-- name: Get :one
SELECT * FROM "group".groups WHERE id = $1;

-- name: Create :one
INSERT INTO
    "group".groups (
        id,
        status,
        name,
        description,
        user_id,
        deleted_at
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: Update :one
UPDATE "group".groups 
SET
    name = $2 ,
    description = $3,
    status = $4,
    deleted_at = $5
WHERE id = $1 RETURNING *;

-- name: Delete :exec
DELETE FROM "group".groups WHERE id = $1;


