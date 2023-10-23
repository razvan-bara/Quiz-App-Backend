-- name: GetAnswer :one
SELECT * FROM answers
WHERE id = $1 LIMIT 1;

-- name: ListQuestionAnswers :many
SELECT * FROM answers
WHERE quiz_question_id = $1
ORDER BY title;

-- name: CreateQuestionAnswer :one
INSERT INTO answers (
    quiz_question_id, title
) VALUES (
             $1, $2
         )
RETURNING *;

-- name: UpdateQuestionAnswer :one
UPDATE answers
set title = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAnswer :exec
DELETE FROM answers
WHERE id = $1;

-- name: DeleteQuestionAnswers :exec
DELETE FROM answers
WHERE quiz_question_id = $1;