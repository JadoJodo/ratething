package thing

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

func (r *SQLiteRepository) MigrateThings() error {
	query := `
	CREATE TABLE IF NOT EXISTS things(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		category TEXT NOT NULL
	);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) CreateThing(thing Thing) (*Thing, error) {
	uuid := uuid.New()

	query := `
	INSERT INTO things(
		uuid, 
		name, 
		category
	) values(
		?,
		?,
		?
	);
	`
	res, err := r.db.Exec(query, uuid, thing.Name, thing.Category)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	thing.ID = id
	thing.UUID = uuid

	return &thing, nil
}

func (r *SQLiteRepository) AllThings() ([]Thing, error) {
	rows, err := r.db.Query("SELECT * FROM things")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []Thing

	for rows.Next() {
		var thing Thing
		if err := rows.Scan(&thing.ID, &thing.UUID, &thing.Name, &thing.Category); err != nil {
			return nil, err
		}
		all = append(all, thing)
	}

	return all, nil
}

func (r *SQLiteRepository) GetThingByUUID(uuid string) (*Thing, error) {
	row := r.db.QueryRow("SELECT * FROM things WHERE uuid = ?", uuid)

	var thing Thing
	if err := row.Scan(&thing.ID, &thing.UUID, &thing.Name, &thing.Category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &thing, nil
}

func (r *SQLiteRepository) UpdateThing(id int64, updated Thing) (*Thing, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}
	res, err := r.db.Exec("UPDATE things SET name = ?, category = ? WHERE id = ?", updated.Name, updated.Category, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

func (r *SQLiteRepository) DeleteThing(id int64) error {
	res, err := r.db.Exec("DELETE FROM thing WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
