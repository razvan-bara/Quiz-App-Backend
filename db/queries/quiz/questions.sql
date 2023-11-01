-- name: GetQuestion :one
SELECT * FROM questions
WHERE id = $1 LIMIT 1;

-- name: GetQuestionByIdAndQuizId :one
SELECT * FROM questions
WHERE id = $1 AND quiz_id = $2 LIMIT 1;

-- name: ListQuestions :many
SELECT * FROM questions
WHERE quiz_id = $1
ORDER BY title;

-- name: CreateQuestion :one
INSERT INTO questions (
    quiz_id, title, body
) VALUES (
             $1, $2, $3
         )
RETURNING *;

-- name: UpdateQuestion :one
UPDATE questions
set title = $2,
    body = $3
WHERE id = $1
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM questions
WHERE id = $1;

-- name: DeleteQuestions :exec
DELETE FROM questions
WHERE quiz_id = $1;