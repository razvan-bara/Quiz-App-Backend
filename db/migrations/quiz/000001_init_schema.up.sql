CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE "quizzes" (
                        "id" bigserial PRIMARY KEY,
                        "title" varchar(255) NOT NULL,
                        "description" text,
                        "attempts" integer DEFAULT 0,
                        "UUID" uuid UNIQUE DEFAULT (uuid_generate_v4()),
                        "createdAt" timestamp DEFAULT (now()),
                        "publishedAt" timestamp
);

CREATE TABLE "quiz_questions" (
                                 "id" bigserial PRIMARY KEY,
                                 "title" varchar(255) NOT NULL,
                                 "body" text,
                                 "quiz_id" bigint NOT NULL,
                                 "UUID" uuid UNIQUE DEFAULT (uuid_generate_v4()),
                                 "createdAt" timestamp DEFAULT (now())
);

CREATE TABLE "answers" (
                           "id" bigserial PRIMARY KEY,
                           "title" varchar(255) NOT NULL,
                           "correct" boolean DEFAULT false,
                           "quiz_question_id" bigint NOT NULL,
                           "UUID" uuid UNIQUE DEFAULT (uuid_generate_v4()),
                           "createdAt" timestamp DEFAULT (now())
);

CREATE INDEX ON "quizzes" ("id");

CREATE INDEX ON "quizzes" ("UUID");

CREATE INDEX ON "quiz_questions" ("id");

CREATE INDEX ON "quiz_questions" ("UUID");

CREATE INDEX ON "answers" ("id");

CREATE INDEX ON "answers" ("UUID");

ALTER TABLE "quiz_questions" ADD FOREIGN KEY ("quiz_id") REFERENCES "quizzes" ("id");

ALTER TABLE "answers" ADD FOREIGN KEY ("quiz_question_id") REFERENCES "quiz_questions" ("id");
