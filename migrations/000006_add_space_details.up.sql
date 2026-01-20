-- Add description and floor_plan columns to areas and rooms
-- Add warning column to desks

ALTER TABLE areas ADD COLUMN description TEXT NOT NULL DEFAULT '';
ALTER TABLE areas ADD COLUMN floor_plan TEXT NOT NULL DEFAULT '';

ALTER TABLE rooms ADD COLUMN description TEXT NOT NULL DEFAULT '';
ALTER TABLE rooms ADD COLUMN floor_plan TEXT NOT NULL DEFAULT '';

ALTER TABLE desks ADD COLUMN warning TEXT NOT NULL DEFAULT '';
