-- sqlc queries (queries/tasks.sql)


-- name: GetTask :one
SELECT *
FROM tasks
WHERE id = ?
  AND done = FALSE
  AND deletion_date IS NULL
LIMIT 1;

-- name: ListTasks :many
SELECT *
FROM tasks
WHERE done = FALSE
  AND deletion_date IS NULL
ORDER BY id;

-- name: ListTasksPage :many
SELECT *
FROM tasks
WHERE done = FALSE
  AND deletion_date IS NULL
ORDER BY id
LIMIT ? OFFSET ?;

-- name: ListDoneTasks :many
SELECT *
FROM tasks
WHERE done = TRUE
  AND deletion_date IS NULL
ORDER BY id;

-- name: ListArchivedTasks :many
SELECT *
FROM tasks
WHERE deletion_date IS NOT NULL
ORDER BY deletion_date DESC;

-- name: SearchTasksByTitle :many
SELECT *
FROM tasks
WHERE title LIKE '%' || ? || '%'
  AND deletion_date IS NULL
ORDER BY id;

-- name: SearchTasksByDescription :many
SELECT *
FROM tasks
WHERE description LIKE '%' || ? || '%'
  AND deletion_date IS NULL
ORDER BY id;

-- name: CreateTask :one
INSERT INTO tasks (title, description)
VALUES (?, ?)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET title       = ?,
    description = ?
WHERE id = ?
RETURNING *;

-- name: MarkTaskDone :one
UPDATE tasks
SET done = TRUE
WHERE id = ?
RETURNING *;

-- name: RestoreTask :one
UPDATE tasks
SET done = FALSE
WHERE id = ?
  AND deletion_date IS NULL
RETURNING *;

-- name: DeleteTask :one
UPDATE tasks
SET deletion_date = CURRENT_TIMESTAMP,
    done          = TRUE
WHERE id = ?
RETURNING *;

-- name: DeleteOldTasks :exec
UPDATE tasks
SET deletion_date = CURRENT_TIMESTAMP
WHERE done = TRUE
  AND deletion_date IS NULL
  AND modification_date < datetime('now', '-30 days');

-- name: DeleteOlderTasks :exec
DELETE FROM tasks
WHERE deletion_date IS NOT NULL
  AND deletion_date < datetime('now', '-60 days');