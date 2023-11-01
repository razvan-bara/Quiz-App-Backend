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