-- name: Get :one
SELECT * FROM "group".groups WHERE id = $1;

-- name: GetWithUnscoped :one
SELECT * FROM "group".groups WHERE id = $1 AND deleted_at IS NOT NULL;

-- name: Create :one
INSERT INTO
    "group".groups (
        id,
        status,
        name,
        description,
        user_id
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: Update :one
UPDATE "group".groups 
SET
    name = $2 ,
    description = $3,
    status = $4
WHERE id = $1 RETURNING *;

-- name: Delete :exec
DELETE FROM "group".groups WHERE id = $1;


-- name: CreateGroupMember :one
INSERT INTO "group".group_members
(
    id,
    group_id,
    user_id,
    role
)
VALUES ($1,$2,$3,$4) RETURNING *;

-- name: UpdateGroupMember :one
UPDATE "group".group_members
SET
    role = $2
WHERE id = $1 RETURNING *;

-- name: DeleteGroupMember :exec
DELETE FROM "group".group_members WHERE id = $1;

-- name: GetAllGroupMembers :many
SELECT * FROM "group".group_members WHERE group_id = $1;

-- name: CountGroupMembers :one
SELECT count(*) FROM "group".group_members WHERE group_id = $1;

-- name: DeleteAllGroupMembersByGroupId :exec
DELETE FROM "group".group_members WHERE group_id = $1;

