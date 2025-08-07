CREATE TABLE IF NOT EXISTS Session_Windows (
    id integer PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER,
    window_id INTEGER,
    FOREIGN KEY(session_id) REFERENCES Sessions(id),
    FOREIGN KEY(window_id) REFERENCES Windows(id)
)
