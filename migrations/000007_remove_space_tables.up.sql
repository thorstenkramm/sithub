-- Remove space tables (areas, rooms, desks) - spaces now loaded from YAML only
-- Remove FK constraint on bookings.desk_id

-- SQLite doesn't support DROP FOREIGN KEY, so we recreate the bookings table without it.
-- First, backup existing bookings data
CREATE TABLE bookings_backup AS SELECT * FROM bookings;

-- Drop the old bookings table
DROP TABLE bookings;

-- Recreate bookings table without FK constraint on desk_id
CREATE TABLE bookings (
  id TEXT PRIMARY KEY,
  desk_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  user_name TEXT NOT NULL DEFAULT '',
  booked_by_user_id TEXT NOT NULL DEFAULT '',
  booked_by_user_name TEXT NOT NULL DEFAULT '',
  booking_date TEXT NOT NULL,
  is_guest INTEGER NOT NULL DEFAULT 0,
  guest_email TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  UNIQUE (desk_id, booking_date)
);

-- Restore bookings data
INSERT INTO bookings SELECT * FROM bookings_backup;

-- Drop the backup table
DROP TABLE bookings_backup;

-- Now drop the space tables (order matters due to FKs)
DROP TABLE IF EXISTS desks;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS areas;
