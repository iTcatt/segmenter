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
	
	joinUsersAndSegmentSQL = `
		select us.user_id, s.segment_id from usersegment us
			join segment s on us.segment_id = s.segment_id
   			where us.user_id = $1 and s.segment_name = $2;`
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

	var temp string
	err := row.Scan(&temp)
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

func (ps *PostgresStorage) AddUser(id int) error {
	requestSQL := "select user_id from users where user_id = $1"
	row := ps.conn.QueryRow(context.Background(), requestSQL, id)
	var temp string
	err := row.Scan(&temp)
	if err == nil {
		return storage.ErrAlreadyExist
	}

	insertSQL := "insert into users(user_id) values($1)"
	_, err = ps.conn.Exec(context.Background(), insertSQL, id)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStorage) AddUserToSegment(id int, segment string) error {
	var temp_user_id, temp_segment_id int

	row := ps.conn.QueryRow(context.Background(), joinUsersAndSegmentSQL, id, segment)
	err := row.Scan(&temp_user_id, &temp_segment_id)
	if err == nil {
		return storage.ErrAlreadyExist
	}
	// получаю segment_id по названию сегмента
	row = ps.conn.QueryRow(context.Background(), "select segment_id from segment where segment_name = $1;", segment)
	var segment_id int
	err = row.Scan(&segment_id)
	if err != nil {
		return storage.ErrNotExist
	}

	insertSQL := "insert into usersegment(user_id, segment_id) values($1, $2);"
	_, err = ps.conn.Exec(context.Background(), insertSQL, id, segment_id)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PostgresStorage) DeleteUserFromSegment(id int, segment string) error {
	var user_id, segment_id int
	
	row := ps.conn.QueryRow(context.Background(), joinUsersAndSegmentSQL, id, segment)
	err := row.Scan(&user_id, &segment_id)
	if err != nil {
		return storage.ErrNotExist
	}

	deleteSQL := "delete from usersegment us where us.user_id = $1 and us.segment_id = $2"
	_, err = ps.conn.Exec(context.Background(), deleteSQL, user_id, segment_id)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PostgresStorage) GetUserSegments(id int) (types.User, error) {
	var result types.User
	result.Id = id
	getSegmentsSQL := `
		select s.segment_name from segment s
			join usersegment us on s.segment_id = us.segment_id
			where us.user_id = $1;`

	rows, err := ps.conn.Query(context.Background(), getSegmentsSQL, id)
	if err != nil {
		return types.User{}, err
	} 
	defer rows.Close()

	for rows.Next() {
		var segment_name string
		err := rows.Scan(&segment_name)
		if err != nil {
			return types.User{}, err
		}
		result.Segments = append(result.Segments, segment_name)
	}
	if result.Segments == nil {
		return result, storage.ErrNotExist
	}
	return result, nil
}
