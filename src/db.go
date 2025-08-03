package src

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type TaskRecord struct {
	ID      int
	Title   string
	Elapsed int64
}

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		elapsed INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func LoadTasks(db *sql.DB) ([]TaskRecord, error) {
	rows, err := db.Query(`SELECT id, title, elapsed FROM tasks ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []TaskRecord
	for rows.Next() {
		var t TaskRecord
		err := rows.Scan(&t.ID, &t.Title, &t.Elapsed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func InsertTask(db *sql.DB, title string) (int64, error) {
	res, err := db.Exec(`INSERT INTO tasks (title, elapsed) VALUES (?, 0)`, title)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func UpdateTaskElapsed(db *sql.DB, id int, elapsed int64) error {
	_, err := db.Exec(`UPDATE tasks SET elapsed = ? WHERE id = ?`, elapsed, id)
	return err
}

func UpdateTask(db *sql.DB, id int, title string) error {
	_, err := db.Exec(`UPDATE tasks SET title = ? WHERE id = ?`, title, id)
	return err
}

func DeleteTask(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}
