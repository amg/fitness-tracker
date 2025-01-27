-- name: GetUserInfo :one
SELECT * FROM user_info
WHERE id = $1 LIMIT 1;

-- name: GetUserInfoByEmail :one
SELECT * FROM user_info
WHERE email = $1 LIMIT 1;

-- name: ListUserInfo :many
SELECT * FROM user_info
ORDER BY first_name;

-- name: CreateUserInfo :one
INSERT INTO user_info (
  id, email, first_name, last_name, picture_url
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: DeleteUserInfo :exec
DELETE FROM user_info
WHERE id = $1;
