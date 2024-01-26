-- name: Get :one
SELECT * FROM upload.attachments WHERE id = $1;

-- name: GetByIds :many
SELECT * FROM upload.attachments WHERE id = ANY($1::uuid[]);

-- name: GetAttachmentsByType :many
SELECT * FROM upload.attachments WHERE attachable_type = $1 AND attachable_id = $2 ORDER BY updated_at;

-- name: GetLastAttachmentByType :one
SELECT * FROM upload.attachments WHERE attachable_type = $1 AND attachable_id = $2 ORDER BY updated_at DESC LIMIT 1;

-- name: GetAttachmentsByOptional :many
SELECT * FROM upload.attachments 
WHERE attachable_type = $1 
AND mime_type LIKE '%' || $2 ||'%'
AND (entity_upload_id = $3 OR $3 IS NULL)
ORDER BY created_at DESC;

-- name: GetAttachmentsByUserId :many
SELECT * FROM upload.attachments 
WHERE attachable_type = $1 
AND mime_type LIKE '%' || $2 ||'%'
AND user_id = $3
AND (entity_upload_id = sqlc.narg('entity_upload_id') OR entity_upload_id IS NULL)
ORDER BY created_at DESC;

-- name: Create :one
INSERT INTO upload.attachments 
    (
        id,
        user_id, 
        entity_upload_id, 
        filename, 
        extension,
        mime_type,
        folder,
        url,
        url_thumbnail)
VALUES ($1, $2, $3, $4, $5, $6, $7 ,$8 , $9) RETURNING * ;

-- name: Update :one
UPDATE upload.attachments 
SET
    attachable_type = $2,
    attachable_id = $3,
    entity_upload_id = $4
WHERE id = $1 RETURNING *;

-- name: UpdateByIds :many
UPDATE upload.attachments 
SET
    attachable_type = $2,
    attachable_id = $3,
    entity_upload_id = $4
WHERE id = ANY($1::uuid[]) RETURNING *;

-- name: Delete :exec
DELETE FROM upload.attachments WHERE id = $1;

-- name: DeleteByIds :exec
DELETE FROM upload.attachments WHERE id = ANY($1::uuid[]);
