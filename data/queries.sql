CREATE TABLE IF NOT EXISTS beer (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    type INTEGER NOT NULL,
    style INTEGER NOT NULL
);