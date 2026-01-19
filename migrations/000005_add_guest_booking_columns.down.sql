-- SQLite doesn't support DROP COLUMN directly prior to 3.35.0
-- For safety, we'll recreate the table without the guest columns

-- Step 1: Create new table without guest columns
CREATE TABLE bookings_new (
    id TEXT PRIMARY KEY,
    desk_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    user_name TEXT NOT NULL DEFAULT '',
    booked_by_user_id TEXT NOT NULL DEFAULT '',
    booked_by_user_name TEXT NOT NULL DEFAULT '',
    booking_date TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE(desk_id, booking_date)
);

-- Step 2: Copy data
INSERT INTO bookings_new (id, desk_id, user_id, user_name, booked_by_user_id, booked_by_user_name, booking_date, created_at, updated_at)
SELECT id, desk_id, user_id, user_name, booked_by_user_id, booked_by_user_name, booking_date, created_at, updated_at
FROM bookings;

-- Step 3: Drop old table
DROP TABLE bookings;

-- Step 4: Rename new table
ALTER TABLE bookings_new RENAME TO bookings;
