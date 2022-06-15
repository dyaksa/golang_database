package repository

import (
	"context"
	"database/entity"
	"database/sql"
	"errors"
)

type commentRepositoryImpl struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepositoryImpl{DB: db}
}

func (repo *commentRepositoryImpl) Insert(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	script := "INSERT INTO customer(email, password) VALUES(?, ?)"

	result, err := repo.DB.ExecContext(ctx, script, comment.Email, comment.Password)
	if err != nil {
		return comment, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}
	comment.ID = int32(lastId)
	return comment, nil
}

func (repo *commentRepositoryImpl) FindById(ctx context.Context, id int32) (entity.Comment, error) {
	script := "SELECT id, email, password FROM customer WHERE id = ? LIMIT 1"
	comment := entity.Comment{}

	rows, err := repo.DB.QueryContext(ctx, script, id)
	if err != nil {
		return comment, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&comment.ID, &comment.Email, &comment.Password)
		if err != nil {
			return comment, err
		}
	} else {
		return comment, errors.New("data not found")
	}

	return comment, nil
}

func (repo *commentRepositoryImpl) FindAll(ctx context.Context) ([]entity.Comment, error) {
	script := "SELECT id, email, password FROM customer"
	comments := []entity.Comment{}

	rows, err := repo.DB.QueryContext(ctx, script)
	if err != nil {
		return comments, err
	}
	for rows.Next() {
		comment := entity.Comment{}
		rows.Scan(&comment.ID, &comment.Email, &comment.Password)
		comments = append(comments, comment)
	}
	return comments, nil
}
