CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE "quizzes" (
                        "id" bigserial PRIMARY KEY,
                        "title" varchar(255) NOT NULL,
                        "description" text,
                        "attempts" integer DEFAULT 0,
                        "UUID" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
                        "createdAt" timestamp not null DEFAULT (now()),
                        "publishedAt" timestamp
);

CREATE TABLE "questions" (
                                 "id" bigserial PRIMARY KEY,
                                 "title" varchar(255) NOT NULL,
                                 "body" text,
                                 "quiz_id" bigint NOT NULL,
                                 "UUID" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
                                 "createdAt" timestamp not null DEFAULT (now())
);

CREATE TABLE "answers" (
                           "id" bigserial PRIMARY KEY,
                           "title" varchar(255) NOT NULL,
                           "correct" boolean NOT NULL DEFAULT false,
                           "question_id" bigint NOT NULL,
                           "UUID" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
                           "createdAt" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "quizzes" ("id");

CREATE INDEX ON "quizzes" ("UUID");

CREATE INDEX ON "questions" ("id");

CREATE INDEX ON "questions" ("UUID");

CREATE INDEX ON "answers" ("id");

CREATE INDEX ON "answers" ("UUID");

ALTER TABLE "questions" ADD FOREIGN KEY ("quiz_id") REFERENCES "quizzes" ("id") ON DELETE CASCADE ;

ALTER TABLE "answers" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE CASCADE ;
