-- name: GetAttempt :one
SELECT * FROM attempts
WHERE id = $1 LIMIT 1;

-- name: ListAttemptsForUser :many
SELECT sqlc.embed(attempts), sqlc.embed(quizzes)
FROM attempts
         INNER JOIN quizzes ON attempts.quiz_id = quizzes.id
WHERE user_id = $1
ORDER BY attempts."createdAt" DESC;

-- name: ListAttemptsForUserWhereStatus :many
SELECT sqlc.embed(attempts), sqlc.embed(quizzes)
FROM attempts
         INNER JOIN quizzes ON attempts.quiz_id = quizzes.id
WHERE user_id = $1 AND status = $2
ORDER BY attempts."createdAt" DESC;


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