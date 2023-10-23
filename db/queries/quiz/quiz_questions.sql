-- name: GetQuestion :one
SELECT * FROM quiz_questions
WHERE id = $1 LIMIT 1;

-- name: ListQuizQuestions :many
SELECT * FROM quiz_questions
WHERE quiz_id = $1
ORDER BY title;

-- name: CreateQuizQuestion :one
INSERT INTO quiz_questions (
    quiz_id, title, body
) VALUES (
             $1, $2, $3
         )
RETURNING *;

-- name: UpdateQuizQuestion :one
UPDATE quiz_questions
set title = $2,
    body = $3
WHERE id = $1
RETURNING *;

-- name: DeleteQuizQuestion :exec
DELETE FROM quiz_questions
WHERE id = $1;

-- name: DeleteQuizQuestions :exec
DELETE FROM quiz_questions
WHERE quiz_id = $1;