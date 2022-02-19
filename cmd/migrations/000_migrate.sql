
CREATE TABLE IF NOT EXISTS "users" (
    "id" INTEGER PRIMARY KEY,
    "email" TEXT,
    "password_hash" TEXT
);

CREATE TABLE IF NOT EXISTS "todos" (
    "id" INTEGER PRIMARY KEY,
    "checked" INTEGER,
    "text" TEXT,
    "user_id" INTEGER REFERENCES "users"("id")
);

CREATE TABLE IF NOT EXISTS "sessions" (
    "id" INTEGER PRIMARY KEY,
    "email" TEXT,
    "cookie" TEXT
);