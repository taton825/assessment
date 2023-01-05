CREATE TABLE IF NOT EXISTS "expenses" (
    "id" SERIAL PRIMARY KEY,
    "title" TEXT,
    "amount" FLOAT,
    "note" TEXT,
    "tags" TEXT[]
);

INSERT INTO "expenses" ("title", "amount", "note", "tags") values ('test-title', 45, 'test-note', ARRAY ['test tag','array']);