-- Add indexes for common booking queries.
CREATE INDEX IF NOT EXISTS idx_bookings_booking_date ON bookings (booking_date);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings (user_id);
