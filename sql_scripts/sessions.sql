CREATE TABLE IF NOT EXISTS Sessions (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    working_directory TEXT NOT NULL
);
