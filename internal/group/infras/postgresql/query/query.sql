-- name: Get
SELECT * FROM "group".groups WHERE uuid = $1;

-- name: Create :one
INSERT INTO
    "group".groups (
        uuid,
        status,
        title,
        content,
        user_id,
        group_id,
        created_at,
        updated_at,
        deleted_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: Update :one
UPDATE "group".groups 
SET
    title = $2 ,
    content = $3,
    status = $4,
    deleted_at = $5,
WHERE uuid = $1 RETURNING *;

-- name: Delete
DELETE FROM "group".groups WHERE uuid = $1;


