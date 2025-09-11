CREATE TABLE IF NOT EXISTS Panes (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    working_directory TEXT NOT NULL,
    split_direction INTEGER NOT NULL,
    split_ratio REAL NOT NULL
)
