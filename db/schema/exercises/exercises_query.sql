-- name: GetExercise :one
SELECT * FROM exercise
WHERE id = $1 LIMIT 1;

-- name: ListDefaultExercises :many
SELECT * FROM exercise
WHERE user_id = NULL
ORDER BY id;

-- name: ListExercisesForUser :many
SELECT * FROM exercise
WHERE user_id = $1
ORDER BY id;