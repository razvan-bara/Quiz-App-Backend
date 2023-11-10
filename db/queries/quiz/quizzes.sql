-- name: GetQuiz :one
SELECT * FROM quizzes
WHERE id = $1 LIMIT 1;

-- name: ListQuizzes :many
SELECT * FROM quizzes
ORDER BY "createdAt" DESC;

-- name: ListDraftQuizzes :many
SELECT * FROM quizzes
where "publishedAt" is null
ORDER BY "createdAt" DESC;

-- name: ListPublishedQuizzesByTitlePaginate :many
SELECT * FROM quizzes
where "publishedAt" is not null and position(lower($1) in lower(title))>0
ORDER BY "createdAt" DESC LIMIT 8 OFFSET $2 * 8;

-- name: ListPublishedQuizzes :many
SELECT * FROM quizzes
where "publishedAt" is not null
ORDER BY "createdAt" DESC;

-- name: CreateQuiz :one
INSERT INTO quizzes (
    title, description, "publishedAt"
) VALUES (
             $1, $2, $3
         )
RETURNING *;

-- name: UpdateQuiz :one
UPDATE quizzes
set title = $2,
    description = $3,
    "publishedAt" = $4
WHERE id = $1
RETURNING *;

-- name: DeleteQuiz :exec
DELETE FROM quizzes
WHERE id = $1;