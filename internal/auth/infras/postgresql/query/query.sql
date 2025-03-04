-- name: CreateUser :one 
INSERT INTO auth.users
    (
        id,
        email,
        password,
        first_name,
        last_name,
        full_name,
        nick_name
    )
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetUser :one
SELECT * FROM auth.users WHERE id = $1 AND resigned = FALSE;

-- name: GetUserByEmail :one
SELECT * FROM auth.users WHERE email = $1 AND resigned = FALSE;

-- name: GetUsers :many
SELECT * FROM auth.users 
WHERE resigned = FALSE AND full_name LIKE COALESCE('%'||$1||'%','%%')
LIMIT $2
OFFSET $3;


-- name: GetUsersInviteGroup :many
SELECT * FROM auth.users 
WHERE 
    resigned = FALSE AND id != ANY($1::uuid[])
LIMIT $2
OFFSET $3;

-- name: GetUsersBirthDayByCurrentMonth :many
SELECT * FROM auth.users 
WHERE 
    resigned = FALSE AND EXTRACT(MONTH FROM date_of_birth) = EXTRACT(MONTH FROM CURRENT_DATE);

-- name: GetUsersBirthDayByCurrentDay :many
SELECT * FROM auth.users 
WHERE resigned = FALSE AND EXTRACT(MONTH FROM date_of_birth) = EXTRACT(MONTH FROM CURRENT_DATE)
AND EXTRACT(DAY FROM date_of_birth) = EXTRACT(DAY FROM CURRENT_DATE);

-- name: UpdateUser :one
UPDATE auth.users 
SET
    avatar_url = COALESCE(sqlc.narg(avatar_url),avatar_url),
    profile_url = COALESCE(sqlc.narg(profile_url),profile_url),
    gender = COALESCE(sqlc.narg(gender),gender),
    phone = COALESCE(sqlc.narg(phone),phone),
    address = COALESCE(sqlc.narg(address),address),
    date_of_birth = COALESCE(sqlc.narg(date_of_birth),date_of_birth),
    position = COALESCE(sqlc.narg(position),position),
    settings = COALESCE(sqlc.narg(settings),settings)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: FindKeyByUserID :one
SELECT * FROM auth.keys WHERE user_id = $1;

-- name: FindKeyByRefreshToken :one
SELECT * FROM auth.keys WHERE refresh_token = $1;

-- name: FindKeyByRefreshTokenUsed :one
SELECT * FROM auth.keys WHERE refresh_tokens_used = $1;

-- name: DeleteKeyByID :exec
DELETE FROM auth.keys WHERE id = $1;

-- name: DeleteKeyByUserID :exec
DELETE FROM auth.keys WHERE user_id = $1;

-- name: CreateKey :one 
INSERT INTO auth.keys
    (
        user_id,
        public_key,
        private_key
    )
VALUES ($1,$2,$3) RETURNING *;

-- name: UpdateKeyByUserID :one
UPDATE auth.keys 
SET
    public_key = COALESCE(sqlc.narg(public_key),public_key),
    private_key = COALESCE(sqlc.narg(private_key),private_key),
    refresh_token = COALESCE(sqlc.narg(refresh_token),refresh_token),
    refresh_tokens_used = COALESCE(sqlc.narg(refresh_tokens_used),refresh_tokens_used)
WHERE user_id = sqlc.arg(user_id) RETURNING *;
