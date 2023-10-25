-- name: Get :one
SELECT * FROM upload.attachments WHERE id = $1;

-- name: Create :one
INSERT INTO upload.attachments 
    (
        id,
        user_id, 
        filename, 
        extension,
        mime_type,
        version_id,
        url,
        url_thumbnail)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING * ;

-- name: Update :one
UPDATE upload.attachments 
SET
    attachable_type = $2,
    attachable_id = $3
WHERE id = $1 RETURNING *;

-- name: Delete :exec
DELETE FROM upload.attachments WHERE id = $1;


   