-- name: GetAnswer :one
SELECT * FROM answers
WHERE id = $1 LIMIT 1;

-- name: GetAnswerByIdAndQuestionId :one
SELECT * FROM answers
WHERE id = $1 AND  question_id = $2 LIMIT 1;

-- name: ListAnswersForQuestion :many
SELECT * FROM answers
WHERE question_id = $1
ORDER BY title;

-- name: CreateAnswer :one
INSERT INTO answers (
    question_id, title, correct
) VALUES (
             $1, $2, $3
         )
RETURNING *;

-- name: UpdateAnswer :one
UPDATE answers
set title = $2, correct = $3
WHERE id = $1
RETURNING *;

-- name: DeleteAnswer :exec
DELETE FROM answers
WHERE id = $1;

-- name: DeleteAnswersForQuestion :exec
DELETE FROM answers
WHERE question_id = $1;