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
    name = COALESCE(sqlc.narg('name'),name),
    description = COALESCE(sqlc.narg('description'),description),
    status = COALESCE(sqlc.narg('status'),status)
WHERE id = sqlc.arg(id) RETURNING *;

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

-- name: CreateGroupMembers :many
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
    role = COALESCE(sqlc.narg(role),role)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: DeleteGroupMember :exec
DELETE FROM "group".group_members WHERE id = $1;

-- name: GetGroupMembers :many
SELECT * FROM "group".group_members WHERE group_id = $1;

-- name: CountGroupMembers :one
SELECT count(*) FROM "group".group_members WHERE group_id = $1;

-- name: DeleteGroupMembersByGroupId :exec
DELETE FROM "group".group_members WHERE group_id = $1;

-- name: GetGroupsByUserId :many
SELECT g.*
FROM "group".groups AS g
INNER JOIN "group".group_members AS gm ON g.id = gm.group_id
WHERE gm.user_id = $1 LIMIT $2 OFFSET $3;

-- name: GetGroupIdsByUserId :many
SELECT gm.group_id
FROM "group".group_members as gm
WHERE user_id = $1;

