-- SQLite doesn't support DROP COLUMN directly in older versions
-- These columns will remain but be unused if migration is rolled back

-- For SQLite 3.35.0+ (2021-03-12), we could use:
-- ALTER TABLE areas DROP COLUMN description;
-- ALTER TABLE areas DROP COLUMN floor_plan;
-- ALTER TABLE rooms DROP COLUMN description;
-- ALTER TABLE rooms DROP COLUMN floor_plan;
-- ALTER TABLE desks DROP COLUMN warning;

-- For compatibility, we use a no-op down migration
SELECT 1;
