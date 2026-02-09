-- Re-add user_name and booked_by_user_name columns, remove guest_name.
CREATE TABLE bookings_old (
  id TEXT PRIMARY KEY,
  item_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  user_name TEXT NOT NULL DEFAULT '',
  booked_by_user_id TEXT NOT NULL DEFAULT '',
  booked_by_user_name TEXT NOT NULL DEFAULT '',
  booking_date TEXT NOT NULL,
  is_guest INTEGER NOT NULL DEFAULT 0,
  guest_email TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  UNIQUE (item_id, booking_date)
);

INSERT INTO bookings_old (id, item_id, user_id, user_name, booked_by_user_id,
  booked_by_user_name, booking_date, is_guest, guest_email, created_at, updated_at)
SELECT id, item_id, user_id,
  CASE WHEN is_guest = 1 THEN guest_name ELSE '' END,
  booked_by_user_id, '', booking_date, is_guest, guest_email,
  created_at, updated_at
FROM bookings;

DROP TABLE bookings;
ALTER TABLE bookings_old RENAME TO bookings;
