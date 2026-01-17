package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (string, error) {
	sql := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(sql, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", err
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	idString := strconv.FormatInt(insertId, 10)
	return idString, nil
}

func Tasks(limit int) ([]*Task, error) {
	sql := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
	rows, err := db.Query(sql, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*Task, 0)
	for rows.Next() {
		t := Task{}
		var dbId int64
		err := rows.Scan(&dbId, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		t.ID = strconv.FormatInt(dbId, 10)
		result = append(result, &t)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if result == nil {
		result = []*Task{}
	}

	return result, nil
}

func GetTask(id string) (*Task, error) {
	taskId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, sql.ErrNoRows
	}

	sql := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	t := Task{}
	var dbId int64
	err = db.QueryRow(sql, taskId).Scan(&dbId, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}

	t.ID = strconv.FormatInt(dbId, 10)
	return &t, nil
}

func UpdateTask(task *Task) error {
	taskId, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("incorrect id for updating task")
	}

	sql := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	result, err := db.Exec(sql, task.Date, task.Title, task.Comment, task.Repeat, taskId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	taskId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("incorrect id for updating task")
	}

	sql := `UPDATE scheduler SET date = ? WHERE id = ?`
	result, err := db.Exec(sql, next, taskId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}

	return nil
}

func DeleteTask(id string) error {
	taskId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("incorrect id for deleting task")
	}

	sql := `DELETE FROM scheduler WHERE id = ?`
	result, err := db.Exec(sql, taskId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("incorrect id for deleting task")
	}

	return nil
}
