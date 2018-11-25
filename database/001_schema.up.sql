CREATE SCHEMA IF NOT EXISTS focus;

CREATE TABLE focus.events (
  "id" SERIAL PRIMARY KEY,
  "user_id" VARCHAR(36) NOT NULL,
  "device_id" VARCHAR(36) NOT NULL,
  "group_id" VARCHAR(36) NOT NULL,
  "window_name" VARCHAR(100),
  "process_name" VARCHAR(400),
  "afk" BOOLEAN DEFAULT false,
  "time" TIMESTAMP NOT NULL
)