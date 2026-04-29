package fasts

import (
	"database/sql"
	"errors"
	"time"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db}
}

type scanner interface {
	Scan(dest ...any) error
}

func scanFast(s scanner) (fastRecord, error) {
	var f fastRecord
	var startTime string
	var endTime sql.NullString
	var createdAt string
	var updatedAt string

	err := s.Scan(&f.ID, &f.UserID, &startTime, &endTime, &f.GoalSeconds, &createdAt, &updatedAt)
	if err != nil {
		return fastRecord{}, err
	}

	f.StartTime, err = time.Parse(time.RFC3339, startTime)
	if err != nil {
		return fastRecord{}, err
	}

	f.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return fastRecord{}, err
	}

	f.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return fastRecord{}, err
	}

	if endTime.Valid {
		t, err := time.Parse(time.RFC3339, endTime.String)
		if err != nil {
			return fastRecord{}, err
		}
		f.EndTime = &t
	}

	return f, nil
}

func (s *Store) GetByID(id, userID int) (Fast, error) {
	query := `SELECT * FROM fasts WHERE id = ? and user_id = ?`
	row := s.db.QueryRow(query, id, userID)

	var f fastRecord

	f, err := scanFast(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Fast{}, ErrNoRecord
		}
		return Fast{}, err
	}

	return compute(f), nil
}

func (s *Store) Start(goal, userID int) (int, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO fasts (start_time, user_id, goal, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?);`
	result, err := s.db.Exec(query, now, userID, goal, now, now)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Store) End(userID int) error {
	now := time.Now().UTC().Format(time.RFC3339)

	query := `UPDATE fasts SET end_time = ?, updated_at = ?
			WHERE user_id = ? AND end_time IS NULL;
`
	_, err := s.db.Exec(query, now, now, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Delete(id, userID int) error {
	query := `DELETE FROM fasts WHERE id = ? AND user_id = ?;`
	_, err := s.db.Exec(query, id, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetActiveFast(userID int) (Fast, error) {
	query := `SELECT * FROM fasts WHERE end_time IS NULL AND user_id = ? LIMIT 1;`
	row := s.db.QueryRow(query, userID)
	var f fastRecord

	f, err := scanFast(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Fast{}, ErrNoRecord
		}
		return Fast{}, err
	}
	return compute(f), nil
}

func (s *Store) GetHistory(userID int) ([]Fast, error) {
	query := `SELECT * FROM fasts WHERE user_id = ? AND end_time IS NOT NULL ORDER BY id DESC LIMIT 10`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fasts []Fast

	var f fastRecord

	for rows.Next() {
		f, err = scanFast(rows)
		if err != nil {
			return nil, err
		}
		fasts = append(fasts, compute(f))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fasts, nil
}
