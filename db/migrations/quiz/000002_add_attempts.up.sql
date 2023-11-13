CREATE TABLE "attempts" (
                           "id" bigserial PRIMARY KEY,
                           "quiz_id" bigint not null,
                            "user_id" bigint not null,
                           "status" smallint not null default(0),
                           "score" smallint not null default(0),
                           "UUID" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
                           "createdAt" timestamp not null DEFAULT (now())
);

CREATE TABLE "attempt_answers" (
                            "id" bigserial PRIMARY KEY,
                            "attempt_id" bigint not null,
                            "question_id" bigint not null,
                            "answer_id" bigint not null,
                            "UUID" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
                            "createdAt" timestamp not null DEFAULT (now())
);

CREATE INDEX ON "attempts" ("id");

CREATE INDEX ON "attempts" ("UUID");

CREATE INDEX ON "attempt_answers" ("id");

CREATE INDEX ON "attempt_answers" ("UUID");

ALTER TABLE "attempts" ADD FOREIGN KEY ("quiz_id") REFERENCES "quizzes" ("id") ON DELETE CASCADE;

ALTER TABLE "attempt_answers" ADD FOREIGN KEY ("attempt_id") REFERENCES "attempts" ("id") ON DELETE CASCADE;

ALTER TABLE "attempt_answers" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE CASCADE;

ALTER TABLE "attempt_answers" ADD FOREIGN KEY ("answer_id") REFERENCES "answers" ("id") ON DELETE CASCADE ;