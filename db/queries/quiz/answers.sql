-- name: GetAnswer :one
SELECT * FROM answers
WHERE id = $1 LIMIT 1;

-- name: ListAnswers :many
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
set title = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAnswer :exec
DELETE FROM answers
WHERE id = $1;

-- name: DeleteAnswers :exec
DELETE FROM answers
WHERE question_id = $1;