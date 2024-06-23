CREATE TABLE example_cities (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  population INTEGER NOT NULL,
  founding_iso8601 TEXT NOT NULL
);

INSERT INTO example_cities (id, name, population, founding_iso8601)
VALUES
  ("15bb1737-26cc-4a9e-9bb2-852d1606010f", "Copenhagen", 1476988, "1254-01-01T00:00:00Z"),
  ("11abc989-2350-444b-9d62-3f243fa988d7", "Stockholm", 1720000, "1252-01-01T00:00:00Z");
