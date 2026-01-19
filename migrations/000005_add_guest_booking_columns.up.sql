-- Add guest booking support
ALTER TABLE bookings ADD COLUMN is_guest INTEGER NOT NULL DEFAULT 0;
ALTER TABLE bookings ADD COLUMN guest_email TEXT NOT NULL DEFAULT '';
