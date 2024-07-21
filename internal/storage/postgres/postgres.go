package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/iTcatt/segmenter/internal/config"
	"github.com/iTcatt/segmenter/internal/models"
	"github.com/iTcatt/segmenter/internal/storage"

	"github.com/jackc/pgx/v5"
)

const (
	createUsersSQL   = `CREATE TABLE if NOT EXISTS users(user_id INT PRIMARY KEY);`
	createSegmentSQL = `
		CREATE TABLE if NOT EXISTS segment(
			segment_id serial PRIMARY KEY,
			segment_name text NOT NULL
		);`
	createUserSegmentSQL = `
		CREATE TABLE if NOT EXISTS user_segment(
			user_id INT,
			segment_id INT,
			FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE ,
			FOREIGN KEY (segment_id) REFERENCES segment (segment_id) ON DELETE CASCADE
		);`

	joinUsersAndSegmentSQL = `
		SELECT s.segment_id 
		FROM user_segment us
		JOIN segment s on us.segment_id = s.segment_id
		WHERE us.user_id = $1 AND s.segment_name = $2;`
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(cfg config.DatabaseConfig) (*Storage, error) {
	dbPath := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	ticker := time.NewTicker(1 * time.Second)
	deadline := time.After(cfg.Timeout)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			conn, err := pgx.Connect(context.Background(), dbPath)
			if err != nil {
				continue
			}

			if err = conn.Ping(context.Background()); err != nil {
				continue
			}

			log.Println("Successful database connection")
			return &Storage{conn: conn}, nil
		case <-deadline:
			return nil, fmt.Errorf("timed out waiting for postgres connection")
		}
	}
}

// StartUp create tables: users, segment, user_segment
func (s *Storage) StartUp() error {
	_, err := s.conn.Exec(context.Background(), createUsersSQL)
	if err != nil {
		return err
	}
	log.Println("Table users created successfully!")

	_, err = s.conn.Exec(context.Background(), createSegmentSQL)
	if err != nil {
		return err
	}
	log.Println("Table segment created successfully!")

	_, err = s.conn.Exec(context.Background(), createUserSegmentSQL)
	if err != nil {
		return err
	}
	log.Println("Table user_segment created successfully!")

	return nil
}

func (s *Storage) CreateSegment(ctx context.Context, name string) error {
	isCreated, err := s.isSegmentCreated(ctx, name)
	if err != nil {
		return err
	}
	if isCreated {
		return storage.ErrAlreadyExist
	}

	insertSQL := "INSERT INTO segment(segment_name) VALUES($1);"
	if _, err = s.conn.Exec(ctx, insertSQL, name); err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteSegment(ctx context.Context, name string) error {
	log.Println("[DEBUG] Delete segment:", name)
	tag, err := s.conn.Exec(ctx, "DELETE FROM segment WHERE segment_name = $1", name)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return storage.ErrNotExist
	}
	return nil
}

func (s *Storage) CreateUser(ctx context.Context, id int) error {
	var tempUserID int
	requestSQL := "SELECT user_id FROM users WHERE user_id = $1"
	row := s.conn.QueryRow(ctx, requestSQL, id)
	err := row.Scan(&tempUserID)
	if err == nil {
		return storage.ErrAlreadyExist
	}

	insertSQL := "INSERT INTO users(user_id) VALUES($1)"
	_, err = s.conn.Exec(ctx, insertSQL, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	tag, err := s.conn.Exec(ctx, "DELETE FROM users WHERE user_id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return storage.ErrNotExist
	}
	return nil
}

func (s *Storage) AddUserToSegment(ctx context.Context, userID int, segment string) error {
	segmentID, err := s.getSegmentIDByName(ctx, segment)
	if err != nil {
		return err
	}

	var tempSegmentID int
	row := s.conn.QueryRow(ctx, joinUsersAndSegmentSQL, userID, segmentID)
	if err := row.Scan(&tempSegmentID); err == nil {
		return storage.ErrAlreadyExist
	}

	_, err = s.conn.Exec(ctx, "INSERT INTO user_segment(user_id, segment_id) VALUES($1, $2);", userID, segmentID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteUserFromSegment(ctx context.Context, userID int, segment string) error {
	segmentID, err := s.getSegmentIDByName(ctx, segment)
	if err != nil {
		return err
	}

	deleteSQL := "DELETE FROM user_segment us WHERE us.user_id = $1 AND us.segment_id = $2"
	_, err = s.conn.Exec(ctx, deleteSQL, userID, segmentID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUser(ctx context.Context, id int) (models.User, error) {
	var tempUserID int
	row := s.conn.QueryRow(ctx, "SELECT user_id FROM users WHERE user_id = $1", id)
	err := row.Scan(&tempUserID)
	if err != nil {
		return models.User{}, storage.ErrNotExist
	}
	user := models.User{
		ID:       id,
		Segments: []string{},
	}
	getSegmentsSQL := `
		SELECT s.segment_name 
		FROM segment s
		JOIN user_segment us ON s.segment_id = us.segment_id
		WHERE us.user_id = $1;`

	rows, err := s.conn.Query(ctx, getSegmentsSQL, id)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var segmentName string
		err := rows.Scan(&segmentName)
		if err != nil {
			return models.User{}, err
		}
		user.Segments = append(user.Segments, segmentName)
	}

	return user, nil
}

func (s *Storage) IsUserCreated(ctx context.Context, userID int) (bool, error) {
	row := s.conn.QueryRow(ctx, "SELECT 1 FROM users WHERE user_id = $1", userID)
	err := row.Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Storage) isSegmentCreated(ctx context.Context, segmentName string) (bool, error) {
	var temp int
	row := s.conn.QueryRow(ctx, "SELECT 1 FROM segment WHERE segment_name = $1", segmentName)
	err := row.Scan(&temp)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Storage) getSegmentIDByName(ctx context.Context, name string) (int, error) {
	var segmentID int
	row := s.conn.QueryRow(ctx, `SELECT segment_id FROM segment WHERE segment_name = $1;`, name)
	if err := row.Scan(&segmentID); err != nil {
		return 0, storage.ErrNotExist
	}
	return segmentID, nil
}
