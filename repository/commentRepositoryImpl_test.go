package repository

import (
	"context"
	"database/databases"
	"database/entity"
	"fmt"
	"testing"
)

func TestCommentImpl(t *testing.T) {
	db := databases.GetConnection()
	defer db.Close()

	ctx := context.Background()
	comment := entity.Comment{
		ID:       1,
		Email:    "diasnour@gmail.com",
		Password: "random",
	}

	commentRepository := NewCommentRepository(db)
	comment, err := commentRepository.Insert(ctx, comment)
	if err != nil {
		panic(err)
	}
	fmt.Println("ID :", comment.ID)
}

func TestFindById(t *testing.T) {
	db := databases.GetConnection()
	defer db.Close()

	ctx := context.Background()

	commentRepository := NewCommentRepository(db)

	comment, err := commentRepository.FindById(ctx, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println(comment)
}
