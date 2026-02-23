CREATE TABLE IF NOT EXISTS transactions (
    id         TEXT PRIMARY KEY,
    account_id TEXT NOT NULL,
    amount     REAL NOT NULL,
    created_at TEXT NOT NULL
);
