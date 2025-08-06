CREATE TABLE IF NOT EXISTS Session_Windows (
    session_id INTEGER,
    window_id INTEGER,
    PRIMARY KEY(session_id, window_id),
    FOREIGN KEY(session_id) REFERENCES Sessions(id),
    FOREIGN KEY(window_id) REFERENCES Windows(id)
)
