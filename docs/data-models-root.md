# Data Models (root)

## Scan Method
- Quick scan (pattern-based). No source files were read.
- Searched for migrations, schemas, and ORM configs.

## Detected Evidence
- No schema or migration files detected.

## Known Data Store (from README)
- SQLite3 (embedded database).

## Likely Entities (inferred from README)
- Areas, Rooms, Desks, Desk Equipment
- Bookings (single day/week/period)
- Users and Groups (for access control)
- Guests (non-account bookings)

## Next Steps to Complete
- Provide migration/schema files (or enable deep scan) to document tables and relationships.
- Confirm how bookings, rooms, and desks relate in the schema.
