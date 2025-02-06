-- name: GetRefreshToken :one
SELECT * FROM refresh_token_jti
WHERE id = $1 LIMIT 1;

-- name: ListRefreshToken :many
SELECT * FROM refresh_token_jti
WHERE user_id = $1
ORDER BY id;

-- name: CreateRefreshToken :one
INSERT INTO refresh_token_jti (
  id, user_id, fingerprint, expires_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_token_jti
WHERE id = $1;

-- name: DeleteRefreshTokenByUserAndFingerprint :exec
DELETE FROM refresh_token_jti
WHERE user_id = $1 AND fingerprint = $2;
