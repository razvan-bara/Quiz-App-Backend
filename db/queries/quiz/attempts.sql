-- name: GetAttempt :one
SELECT * FROM attempts
WHERE id = $1 LIMIT 1;

-- name: CreateAttempt :one
INSERT INTO attempts (
    quiz_id, user_id
) VALUES (
             $1, $2
         )
RETURNING *;

-- name: CreateAttemptAnswer :one
INSERT INTO attempt_answers (
    attempt_id, question_id, answer_id
) VALUES (
             $1, $2, $3
         )
RETURNING *;

-- name: UpdateAttempt :one
UPDATE attempts
set score = $2, status = $3
WHERE id = $1
RETURNING *;