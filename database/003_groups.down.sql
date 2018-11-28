DROP TABLE "focus.groups";
DROP TABLE "focus.event_group";

ALTER TABLE focus.events ADD COLUMN "group_id" INT NOT NULL;
