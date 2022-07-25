package thing

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

func (r *SQLiteRepository) MigrateRating() error {
	query := `
	CREATE TABLE IF NOT EXISTS(ratings
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT NOT NULL UNIQUE,
		candidates TEXT NOT NULL UNIQUE,
		winner TEXT NOT NULL,
		loser TEXT NOT NULL
	);
	`

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) CreateRating(rating Rating) (*Rating, error) {
	query := `
	INSERT INTO ratings(
		uuid, 
		candidates, 
		winner,
		loser
	) values(
		?,
		?,
		?,
		?
	);
	`
	res, err := r.db.Exec(query, rating.UUID, rating.Candidates, rating.Winner, rating.Loser)

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
	rating.ID = id

	return &rating, nil
}

func (r *SQLiteRepository) AllRatings() ([]Rating, error) {
	rows, err := r.db.Query("SELECT * FROM ratings")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []Rating

	for rows.Next() {
		var rating Rating
		if err := rows.Scan(&rating.ID, &rating.UUID, &rating.Candidates, &rating.Winner, &rating.Loser); err != nil {
			return nil, err
		}
		all = append(all, rating)
	}

	return all, nil
}

func (r *SQLiteRepository) GetRatingByUUID(uuid string) (*Rating, error) {
	row := r.db.QueryRow("SELECT * FROM ratings WHERE uuid = ?", uuid)

	var rating Rating
	if err := row.Scan(&rating.ID, &rating.UUID, &rating.Candidates, &rating.Winner, &rating.Loser); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &rating, nil
}

func (r *SQLiteRepository) UpdateRating(id int64, updated Rating) (*Rating, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}
	res, err := r.db.Exec("UPDATE ratings SET winner = ?, loser = ? WHERE id = ?", updated.Winner, updated.Loser, id)
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

func (r *SQLiteRepository) DeleteRating(id int64) error {
	res, err := r.db.Exec("DELETE FROM ratings WHERE id = ?", id)
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
