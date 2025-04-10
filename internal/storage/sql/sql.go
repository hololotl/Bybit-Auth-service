package sql

import (
	"Bybit_Pet_Project/internal/domain/models"
	"database/sql"
	"fmt"
	"golang.org/x/net/context"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.sqlite.SaveUser"
	//TODO validate email for unique value
	var id int64
	r := s.db.QueryRow("INSERT INTO users (email, pass_hash) VALUES ($1, $2) returning id", email, passHash).Scan(&id)
	if r != nil {
		fmt.Println(r)
		return 0, fmt.Errorf("%s: %w", op, r)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	row := s.db.QueryRow("select * from users where email = $1", email)
	user := models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) App(ctx context.Context, id int64) (models.App, error) {
	const op = "storage.sqlite.App"
	row := s.db.QueryRow("select * from apps where id = $1", id)
	app := models.App{}
	err := row.Scan(&app.ID, &app.Name)
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}
