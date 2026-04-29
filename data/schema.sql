CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL CHECK (length(hashed_password) = 60),
    created TEXT NOT NULL,
    updated TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS fasts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    start_time TEXT NOT NULL,
    end_time TEXT,
    goal INTEGER NOT NULL CHECK (goal > 18000 AND goal < 604800),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    CHECK (end_time IS NULL OR end_time > start_time)
);

CREATE UNIQUE INDEX IF NOT EXISTS one_active_fast_per_user
ON fasts (user_id)
WHERE end_time IS NULL;

CREATE INDEX IF NOT EXISTS fasts_user_id_idx ON fasts (user_id);

CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
