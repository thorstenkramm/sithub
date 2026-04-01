CREATE TABLE floor_plan_positions (
  id TEXT PRIMARY KEY,
  floor_plan TEXT NOT NULL,
  item_id TEXT NOT NULL,
  label TEXT NOT NULL DEFAULT '',
  x REAL NOT NULL,
  y REAL NOT NULL,
  width REAL NOT NULL,
  height REAL NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  UNIQUE (floor_plan, item_id)
);

CREATE INDEX idx_fpp_floor_plan ON floor_plan_positions(floor_plan);
