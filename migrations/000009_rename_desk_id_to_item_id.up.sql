-- Rename desk_id to item_id in bookings table (domain rename: desks → items)
ALTER TABLE bookings RENAME COLUMN desk_id TO item_id;
