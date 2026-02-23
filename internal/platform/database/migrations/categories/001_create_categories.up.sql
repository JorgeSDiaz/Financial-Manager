CREATE TABLE IF NOT EXISTS categories (
    id        TEXT    PRIMARY KEY,
    name      TEXT    NOT NULL,
    type      TEXT    NOT NULL CHECK(type IN ('expense', 'income')),
    is_system INTEGER NOT NULL DEFAULT 0
);
