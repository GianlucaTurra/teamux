CREATE TABLE IF NOT EXISTS Windows_Panes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    window_id INTEGER,
    pane_id INTEGER,
    FOREIGN KEY(window_id) REFERENCES Windows(id),
    FOREIGN KEY(pane_id) REFERENCES Panes(id)
);
