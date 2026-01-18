PRAGMA foreign_keys=OFF;

CREATE TABLE areas_new (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

INSERT INTO areas_new (id, name, created_at, updated_at)
SELECT id, name, created_at, updated_at
FROM areas;

DROP TABLE areas;

ALTER TABLE areas_new RENAME TO areas;

PRAGMA foreign_keys=ON;
