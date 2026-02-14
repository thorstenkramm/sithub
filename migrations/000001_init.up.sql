-- SitHub database schema
-- Spaces (areas, item groups, items) are loaded from YAML configuration only.

CREATE TABLE users (
  id TEXT PRIMARY KEY,
  email TEXT NOT NULL,
  display_name TEXT NOT NULL,
  password_hash TEXT NOT NULL DEFAULT '',
  user_source TEXT NOT NULL CHECK (user_source IN ('internal', 'entraid')),
  entra_id TEXT NOT NULL DEFAULT '',
  is_admin INTEGER NOT NULL DEFAULT 0,
  last_login TEXT NOT NULL DEFAULT '',
  access_token TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_entra_id ON users(entra_id);

CREATE TABLE bookings (
  id TEXT PRIMARY KEY,
  item_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  booked_by_user_id TEXT NOT NULL DEFAULT '',
  booking_date TEXT NOT NULL,
  is_guest INTEGER NOT NULL DEFAULT 0,
  guest_name TEXT NOT NULL DEFAULT '',
  guest_email TEXT NOT NULL DEFAULT '',
  note TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  UNIQUE (item_id, booking_date)
);

CREATE INDEX idx_bookings_booking_date ON bookings(booking_date);
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
