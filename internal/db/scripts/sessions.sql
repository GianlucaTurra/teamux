CREATE TABLE IF NOT EXISTS Session_files (
    id INTEGER PRIMARY KEY,
    file TEXT NOT NULL UNIQUE,
    name NOT NULL UNIQUE
)
