package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/luischitala/go-server/models"
)

type PostgresRepository struct {
	db *sql.DB
}

//Constructor

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

//Function to insert 1 user to db
func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	//Will execute sql sentence
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", user.Id, user.Email, user.Password)
	return err
}

//Function to insert 1 user to db
func (repo *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	//Will execute sql sentence
	_, err := repo.db.ExecContext(ctx, "INSERT INTO posts (id, post_content, user_id) VALUES ($1, $2, $3)", post.Id, post.PostContent, post.UserId)
	return err
}

func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1 ", id)
	//Once it finishes close the connection
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var user = models.User{}
	for rows.Next() {
		//If there are not errors in the copy process, return the user
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = $1 ", email)
	//Once it finishes close the connection
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var user = models.User{}
	for rows.Next() {
		//If there are not errors in the copy process, return the user
		if err = rows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, created_at, user_id FROM posts WHERE id = $1 ", id)
	//Once it finishes close the connection
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var post = models.Post{}
	for rows.Next() {
		//If there are not errors in the copy process, return the user
		if err = rows.Scan(&post.Id, &post.PostContent, &post.CreatedAt, &post.UserId); err == nil {
			return &post, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &post, nil
}

func (repo *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	//Only allow edit for owners
	_, err := repo.db.ExecContext(ctx, "UPDATE posts SET post_content = $1 WHERE id = $2 and user_id = $3", post.PostContent, post.Id, post.UserId)
	return err
}

func (repo *PostgresRepository) DeletePost(ctx context.Context, id string, userId string) error {
	//Only allow edit for owners
	_, err := repo.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1 and user_id = $2", id, userId)
	return err
}
