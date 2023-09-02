package postgres

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/iTcatt/avito-task/internal/storage"
	"github.com/iTcatt/avito-task/internal/types"
	"github.com/jackc/pgx/v5"
)

type PostgresStorage struct {
	conn *pgx.Conn
}

const (
	createUsersSQL   = "create table if not exists users(user_id int primary key);"
	createSegmentSQL = `
create table if not exists segment(
	segment_id serial primary key,
	segment_name text not null
);`
	createUserSegmentSQL = `
create table if not exists usersegment(
	user_id int,
	segment_id int,
	foreign key (user_id) references users (user_id),
	foreign key (segment_id) references segment (segment_id)
);`
)

func NewPostgresStorage() (*PostgresStorage, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	log.Println("Successful database connection")

	return &PostgresStorage{conn: conn}, nil
}

// create tables: user, segment, usersegment
func (ps *PostgresStorage) StartUp() error {
	_, err := ps.conn.Exec(context.Background(), createUsersSQL)
	if err != nil {
		return err
	}
	log.Println("Table users created successfully!")

	_, err = ps.conn.Exec(context.Background(), createSegmentSQL)
	if err != nil {
		return err
	}
	log.Println("Table segment created successfully!")

	_, err = ps.conn.Exec(context.Background(), createUserSegmentSQL)
	if err != nil {
		return err
	}
	log.Println("Table usersegment created successfully!")

	return nil
}

func (ps *PostgresStorage) CreateSegment(name string) error {
	requestSQL := "SELECT segment_name FROM segment WHERE segment_name = $1;"
	row := ps.conn.QueryRow(context.Background(), requestSQL, name)

	var result string
	err := row.Scan(&result)
	if err == nil {
		return storage.ErrAlreadyExist
	}
	
	insertSQL := "insert into segment(segment_name) values($1);"
	_, err = ps.conn.Exec(context.Background(), insertSQL, name)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStorage) DeleteSegment(name string) error {
	return nil
}

func (ps *PostgresStorage) AddUser(id int, addedSegments []string, removedSegments []string) error {

	return nil
}

func (ps *PostgresStorage) GetSegments(id int) (types.User, error) {
	return types.User{}, nil
}
