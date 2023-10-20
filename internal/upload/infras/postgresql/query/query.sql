-- name: Get:one
SELECT * FROM upload.attachments WHERE id = $1;

-- name: Create:one

INSERT INTO upload.attachments 
    (
        id,
        attachable_type,
        attachable_id,
        user_id, 
        filename, 
        extension,
        mime_type
        url,
        url_thumbnail
        created_at, 
        updated_at)
VALUES ($1, $2, $3, $4, $5, $6,) RETURNING * ;

-- name: Update:one

UPDATE upload.attachments SET
    attachable_type = $2,
    attachable_id = $3,





   