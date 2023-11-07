-- name: CreateUser :one 
INSERT INTO auth.users
    (
        id,
        email,
        password,
        fist_name,
        last_name,
        full_name
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetUser :one
SELECT * FROM auth.users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM auth.users WHERE email = $1;

-- name: FindKeyByUserID :one
SELECT * FROM auth.keys WHERE user_id = $1;

-- name: FindKeyByRefreshToken :one
SELECT * FROM auth.keys WHERE refresh_token = $1;

-- name: FindKeyByRefreshTokenUsed :one
SELECT * FROM auth.keys WHERE refresh_tokens_used = $1;

-- name: DeleteKeyByID :exec
DELETE FROM auth.keys WHERE id = $1;

-- name: DeleteKeyByUserID :exec
DELETE FROM auth.users WHERE user_id = $1;

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
