CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE "quiz" (
                        "id" bigserial PRIMARY KEY,
                        "title" varchar(255) NOT NULL,
                        "description" text,
                        "attempts" integer DEFAULT 0,
                        "UUID" uuid UNIQUE DEFAULT (uuid_generate_v4()),
                        "createdAt" timestamp DEFAULT (now()),
                        "publishedAt" timestamp
);

CREATE TABLE "quiz_question" (
                                 "id" bigserial PRIMARY KEY,
                                 "title" varchar(255) NOT NULL,
                                 "body" text,
                                 "quiz_id" bigint,
                                 "UUID" uuid UNIQUE DEFAULT (uuid_generate_v4()),
                                 "createdAt" timestamp DEFAULT (now())
);

CREATE TABLE "answers" (
                           "id" bigserial PRIMARY KEY,
                           "title" varchar(255) NOT NULL,
                           "correct" boolean DEFAULT false,
                           "quiz_question_id" bigint,
                           "UUID" uuid UNIQUE DEFAULT (uuid_generate_v4()),
                           "createdAt" timestamp DEFAULT (now())
);

CREATE INDEX ON "quiz" ("id");

CREATE INDEX ON "quiz" ("UUID");

CREATE INDEX ON "quiz_question" ("id");

CREATE INDEX ON "quiz_question" ("UUID");

CREATE INDEX ON "answers" ("id");

CREATE INDEX ON "answers" ("UUID");

ALTER TABLE "quiz_question" ADD FOREIGN KEY ("quiz_id") REFERENCES "quiz" ("id");

ALTER TABLE "answers" ADD FOREIGN KEY ("quiz_question_id") REFERENCES "quiz_question" ("id");