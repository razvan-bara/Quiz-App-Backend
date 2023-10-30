CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE "users" (
                         "id" bigserial PRIMARY KEY,
                         "email" varchar(255) UNIQUE NOT NULL,
                         "firstName" varchar(255) NOT NULL,
                         "lastName" varchar(255) NOT NULL,
                         "password" text NOT NULL,
                         "UUID" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
                         "createdAt" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("id");

CREATE INDEX ON "users" ("UUID");
