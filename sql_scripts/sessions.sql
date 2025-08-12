CREATE TABLE IF NOT EXISTS Sessions (
    id INTEGER PRIMARY KEY,
    name NOT NULL UNIQUE,
    working_directory TEXT NOT NULL
)
