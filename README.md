<p align="center">
  <img src="tesks_logo.png" width="120" alt="Tesks Logo">
</p>

# 📝 Tesks

**Tesks** is a simple, fast, and user-friendly command-line tool for managing tasks, built in Go using SQLite for storage. Tasks can be created, listed, updated, searched, and archived through intuitive commands — all from your terminal.

---

## ⚙️ Features

- ✅ Create, list, update, delete, and search tasks
- 📦 Archive and restore tasks
- ✔️ Mark tasks as completed
- 📊 Tabular task display with age-based coloring
- 💾 Uses SQLite and `sqlc` for efficient query generation
- 🔄 Auto-migration and task cleanup at startup
- 📁 Stores data in `$HOME/.local/share/tesks`

---

## 🚀 Installation

```bash
go install -tags "sqlite_fts5" github.com/ahmedsat/tesks@latest
````

---

## 📚 Usage

All commands are available via:

```bash
tesks [command]
```

### 🔧 Create a task

```bash
tesks create --title "Buy milk" --description "Before it expires"
```

### 📋 List tasks

```bash
tesks list --page 1 --size 10
```

### 🔍 Search tasks

```bash
tesks search --query "milk"
```

### 🛠 Update a task

```bash
tesks update 1 --title "Buy groceries" --description "Milk, bread, eggs"
```

### 🗑 Archive a task

```bash
tesks delete 1
```

### ✔ Mark as done

```bash
tesks markdone 1
```

### ♻ Restore an archived/completed task

```bash
tesks restore 1
```

### 📦 List archived or done tasks

```bash
tesks archived
tesks done
```

---

## 🧱 Database

* Uses SQLite
* SQL queries generated via [`sqlc`](https://github.com/kyleconroy/sqlc)
* Migrations auto-run from embedded `/migrations` folder

---

## 📦 Project Structure

* `main.go`: Entry point and command registration
* `sql/`: Auto-generated code by `sqlc`
* `migrations/`: Embedded SQL migrations
* Uses [`tablewriter`](https://github.com/olekukonko/tablewriter) for pretty output

---

## 🛠 Development

### Install dependencies

```bash
go mod tidy
```

### Generate SQL code

```bash
go generate
```

