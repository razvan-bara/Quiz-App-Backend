-- name: GetQuiz :one
SELECT * FROM quizzes
WHERE id = $1 LIMIT 1;

-- name: ListQuizzes :many
SELECT * FROM quizzes
ORDER BY title;

-- name: CreateQuizzes :one
INSERT INTO quizzes (
    title, description
) VALUES (
             $1, $2
         )
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM quizzes
WHERE id = $1;