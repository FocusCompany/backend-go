CREATE TABLE focus.filters (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(300) NOT NULL,
  "user_id" UUID NOT NULL
);