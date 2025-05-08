-- SQLite schema (migrations/0001_init.sql)

BEGIN;

CREATE TABLE IF NOT EXISTS tasks (
  id                INTEGER PRIMARY KEY,
  title             TEXT    NOT NULL UNIQUE,
  description       TEXT,
  done              BOOLEAN NOT NULL DEFAULT FALSE CHECK(done IN (0,1)),
  creation_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  modification_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deletion_date     DATETIME,
  CHECK(deletion_date IS NULL OR done = 1),
  CHECK(creation_date <= modification_date),
  CHECK(deletion_date IS NULL OR modification_date <= deletion_date)
);

-- trigger to auto‐bump modification_date
CREATE TRIGGER IF NOT EXISTS trg_tasks_moddate
AFTER UPDATE ON tasks
FOR EACH ROW
WHEN NEW.modification_date = OLD.modification_date
BEGIN
  UPDATE tasks
  SET modification_date = CURRENT_TIMESTAMP
  WHERE id = OLD.id;
END;

-- partial indexes
CREATE INDEX IF NOT EXISTS idx_tasks_active
  ON tasks(id)
  WHERE done = FALSE AND deletion_date IS NULL;
CREATE INDEX IF NOT EXISTS idx_tasks_done_unarchived
  ON tasks(id)
  WHERE done = TRUE AND deletion_date IS NULL;

-- full‐text search setup
CREATE VIRTUAL TABLE IF NOT EXISTS tasks_fts USING fts5(
  title, description,
  content='tasks', content_rowid='id'
);

CREATE TRIGGER IF NOT EXISTS trg_tasks_fts_insert
AFTER INSERT ON tasks
BEGIN
  INSERT INTO tasks_fts(rowid, title, description)
    VALUES (new.id, new.title, new.description);
END;

CREATE TRIGGER IF NOT EXISTS trg_tasks_fts_update
AFTER UPDATE ON tasks
BEGIN
  INSERT INTO tasks_fts(tasks_fts, rowid, title, description)
    VALUES('delete', old.id, old.title, old.description);
  INSERT INTO tasks_fts(rowid, title, description)
    VALUES (new.id, new.title, new.description);
END;

CREATE TRIGGER IF NOT EXISTS trg_tasks_fts_delete
AFTER DELETE ON tasks
BEGIN
  INSERT INTO tasks_fts(tasks_fts, rowid, title, description)
    VALUES('delete', old.id, old.title, old.description);
END;

COMMIT;