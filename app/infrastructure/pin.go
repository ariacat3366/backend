package infrastructure

import (
	"app/helpers"
	"app/models"
	"app/repository"
	"database/sql"
)

type Pin struct {
	DB *sql.DB
}

func NewPinRepository(db *sql.DB) repository.PinRepository {
	return &Pin{
		DB: db,
	}
}

func (p *Pin) CreatePin(pin *models.Pin, boardID int) (*models.Pin, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}

	const query = `
INSERT INTO pins (user_id, title, description, url, image_url, is_private) VALUES (?, ?, ?, ?, ?, ?);
`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, helpers.TryRollback(tx, err)
	}

	result, err := stmt.Exec(pin.UserID, pin.Title, pin.Description, pin.URL, pin.ImageURL, pin.IsPrivate)
	if err = helpers.CheckDBExecError(result, err); err != nil {
		return nil, helpers.TryRollback(tx, err)
	}

	pinID, err := result.LastInsertId()
	if err != nil {
		return nil, helpers.TryRollback(tx, err)
	}

	const query2 = `
INSERT INTO boards_pins (board_id, pin_id) VALUES (?, ?);
`

	stmt, err = tx.Prepare(query2)
	if err != nil {
		return nil, helpers.TryRollback(tx, err)
	}

	result, err = stmt.Exec(boardID, pinID)
	if err = helpers.CheckDBExecError(result, err); err != nil {
		return nil, helpers.TryRollback(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, helpers.TryRollback(tx, err)
	}

	return pin, nil
}

func (p *Pin) UpdatePin(pin *models.Pin) error {
	const query = `
UPDATE pins SET title = ?, description = ?, url = ?, image_url = ?, is_private = ?;
`

	stmt, err := p.DB.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(pin.Title, pin.Description, pin.URL, pin.ImageURL, pin.IsPrivate)
	if err := helpers.CheckDBExecError(result, err); err != nil {
		return err
	}

	return nil
}

func (p *Pin) DeletePin(pinID int) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	const query = `
DELETE FROM pins WHERE id = ?;
`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return helpers.TryRollback(tx, err)
	}

	result, err := stmt.Exec(pinID)
	if err := helpers.CheckDBExecError(result, err); err != nil {
		return helpers.TryRollback(tx, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return helpers.TryRollback(tx, err)
	}

	const query2 = `
DELETE FROM boards_pins WHERE pin_id = ?;
`

	stmt, err = tx.Prepare(query2)
	if err != nil {
		return helpers.TryRollback(tx, err)
	}

	result, err = stmt.Exec(pinID)
	if err := helpers.CheckDBExecError(result, err); err != nil {
		return helpers.TryRollback(tx, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return helpers.TryRollback(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return helpers.TryRollback(tx, err)
	}

	return nil
}

func (p *Pin) GetPin(pinID int) (*models.Pin, error) {
	const query = `
SELECT
    p.id,
    p.user_id,
    p.title,
    p.description,
    p.url,
    p.image_url,
    p.is_private,
    p.created_at,
    p.updated_at
FROM
    pins p
WHERE
    p.id = ?;
`

	stmt, err := p.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(pinID)

	pin := &models.Pin{}
	err = row.Scan(
		&pin.ID,
		&pin.UserID,
		&pin.Title,
		&pin.Description,
		&pin.URL,
		&pin.ImageURL,
		&pin.IsPrivate,
		&pin.CreatedAt,
		&pin.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return pin, nil
}

func (p *Pin) GetPinsByBoardID(boardID int, page int) ([]*models.Pin, error) {
	const query = `
SELECT
    p.id,
    p.user_id,
    p.title,
    p.description,
    p.url,
    p.image_url,
    p.is_private,
    p.created_at,
    p.updated_at
FROM
    pins AS p
    JOIN boards_pins AS bp ON p.id = bp.pin_id
WHERE
    bp.board_id = ?
LIMIT ? OFFSET ?;
`
	limit := 10
	offset := (page - 1) * limit

	stmt, err := p.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(boardID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pins []*models.Pin
	for rows.Next() {
		pin := &models.Pin{}
		err := rows.Scan(
			&pin.ID,
			&pin.UserID,
			&pin.Title,
			&pin.Description,
			&pin.URL,
			&pin.ImageURL,
			&pin.IsPrivate,
			&pin.CreatedAt,
			&pin.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pins = append(pins, pin)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pins, nil
}

func (p *Pin) GetPinsByUserID(userID int) ([]*models.Pin, error) {
	const query = `
SELECT
    p.id,
    p.user_id,
    p.title,
    p.description,
    p.url,
    p.image_url,
    p.is_private,
    p.created_at,
    p.updated_at
FROM
    pins p
WHERE
    p.user_id = ?;
`

	stmt, err := p.DB.Prepare(query)

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pins []*models.Pin
	for rows.Next() {
		pin := &models.Pin{}
		err := rows.Scan(
			&pin.ID,
			&pin.UserID,
			&pin.Title,
			&pin.Description,
			&pin.URL,
			&pin.ImageURL,
			&pin.IsPrivate,
			&pin.CreatedAt,
			&pin.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pins = append(pins, pin)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pins, nil
}
