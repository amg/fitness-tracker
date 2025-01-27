-- name: GetRefreshToken :one
SELECT * FROM refresh_token
WHERE id = $1 LIMIT 1;

-- name: ListRefreshToken :many
SELECT * FROM refresh_token
WHERE user_id = $1
ORDER BY id;

-- name: CreateRefreshToken :one
INSERT INTO refresh_token (
  id, user_id, fingerprint
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_token
WHERE id = $1;

-- name: DeleteRefreshTokenByUserAndFingerprint :exec
DELETE FROM refresh_token
WHERE user_id = $1 AND fingerprint = $2;
