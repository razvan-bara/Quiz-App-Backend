-- name: GetQuiz :one
SELECT * FROM quizzes
WHERE id = $1 LIMIT 1;

-- name: ListQuizzes :many
SELECT * FROM quizzes
ORDER BY title;

-- name: CreateQuiz :one
INSERT INTO quizzes (
    title, description
) VALUES (
             $1, $2
         )
RETURNING *;

-- name: UpdateQuiz :one
UPDATE quizzes
set title = $2,
    description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteQuiz :exec
DELETE FROM quizzes
WHERE id = $1;