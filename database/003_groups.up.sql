CREATE TABLE focus.groups (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(300) NOT NULL
);

CREATE TABLE focus.event_group (
  "id" SERIAL PRIMARY KEY,
  "event_id" BIGINT NOT NULL REFERENCES focus.events("id"),
  "group_id" BIGINT NOT NULL REFERENCES focus.groups("id")
);

ALTER TABLE focus.events DROP COLUMN "group_id";
