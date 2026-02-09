-- Remove denormalized user_name and booked_by_user_name columns.
-- Add guest_name column for guest bookings (user_name was the only source of guest names).
-- SQLite requires table recreation to drop columns.
CREATE TABLE bookings_new (
  id TEXT PRIMARY KEY,
  item_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  booked_by_user_id TEXT NOT NULL DEFAULT '',
  booking_date TEXT NOT NULL,
  is_guest INTEGER NOT NULL DEFAULT 0,
  guest_name TEXT NOT NULL DEFAULT '',
  guest_email TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  UNIQUE (item_id, booking_date)
);

INSERT INTO bookings_new (id, item_id, user_id, booked_by_user_id, booking_date,
  is_guest, guest_name, guest_email, created_at, updated_at)
SELECT id, item_id, user_id, booked_by_user_id, booking_date,
  is_guest, CASE WHEN is_guest = 1 THEN user_name ELSE '' END, guest_email,
  created_at, updated_at
FROM bookings;

DROP TABLE bookings;
ALTER TABLE bookings_new RENAME TO bookings;
