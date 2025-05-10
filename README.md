<p align="center">
  <img src="tesk_logo.png" width="120" alt="tesk Logo">
</p>

# tesk

**tesk** is a minimalist terminal-based task manager built in Go, using SQLite for storage and `sqlc` for fast, type-safe queries.  
Designed for speed, clarity, and simplicity.

---

## 🚀 Features

- Add and list tasks
- Mark tasks as done, undo, or delete them
- Archive and clean up old tasks
- Full-text search
- Age-based coloring in terminal table output

---

## 📦 Installation

```bash
git clone https://github.com/ahmedsat/tesk.git
cd tesk
go build -tags "sqlite_fts5" -o tesk
````

---

## 🛠️ Usage Examples

### Add a Task

```bash
$ ./tesk create -t "Buy groceries" -d "Milk, eggs, and bread"
```

### List Tasks

```bash
$ ./tesk list
```

```
| -----+----------------+------------------------+---------+
| ID   | Title          | DESCRIPTION            | Age     |
+------+----------------+------------------------+---------+
| 1    | Buy groceries  | Milk, eggs, and bread  | 2h ago  |
| 2    | Finish report  |                        | 30m ago |
+------+----------------+------------------------+---------+

```

### Mark Task as Done

```bash
$ ./tesk done 1
```

```
Task 1 marked as done.
```

### Search Tasks

```bash
$ ./tesk search -q groceries
```

```
+------+----------------+------------------------+---------+
| ID   | Title          | DESCRIPTION            | Age     |
+------+----------------+------------------------+---------+
| 1    | Buy groceries  | Milk, eggs, and bread  | 2h ago  |
+------+----------------+------------------------+---------+
```

### Archive Old Tasks

```bash
$ ./tesk archive
```

```
TODO: add something here
```

---

## ⚙️ Regenerate SQL (dev only)

```bash
go generate
```

