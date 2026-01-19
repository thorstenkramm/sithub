-- SQLite doesn't support DROP COLUMN directly, so we recreate the table
CREATE TABLE bookings_backup AS SELECT id, desk_id, user_id, user_name, booking_date, created_at, updated_at FROM bookings;
DROP TABLE bookings;
CREATE TABLE bookings (
  id TEXT PRIMARY KEY,
  desk_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  user_name TEXT NOT NULL DEFAULT '',
  booking_date TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  FOREIGN KEY (desk_id) REFERENCES desks(id) ON DELETE CASCADE,
  UNIQUE (desk_id, booking_date)
);
INSERT INTO bookings SELECT * FROM bookings_backup;
DROP TABLE bookings_backup;
