package repository

import (
	"belajar-golang-database/entity"
	"context"
	"database/sql"
	"errors"
	"strconv"
)

type CommentRepositoryImpl struct {
	DB *sql.DB
}

func NewCommentRepository(conn *sql.DB) CommentRepository {
	return &CommentRepositoryImpl{
		DB: conn,
	}
}

func (repo *CommentRepositoryImpl) Insert(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	scrip := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	result, err := repo.DB.ExecContext(ctx, scrip, comment.Email, comment.Comment)
	if err != nil {
		return comment, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}
	comment.Id = int32(id)

	return comment, nil
}

func (repo *CommentRepositoryImpl) FindById(ctx context.Context, id int32) (entity.Comment, error) {

	script := "SELECT id, email, comment FROM comments WHERE id = ?"
	rows, err := repo.DB.QueryContext(ctx, script, id)
	comment := entity.Comment{}
	if err != nil {
		return comment, err
	}

	defer rows.Close()
	if rows.Next() {
		// ada
		rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		return comment, nil
	} else {
		// tidak ada data
		return comment, errors.New("id " + strconv.Itoa(int(id)) + " not found")
	}
}

func (repo *CommentRepositoryImpl) FindAll(ctx context.Context) ([]entity.Comment, error) {
	script := "SELECT id, email, comment FROM comments"
	rows, err := repo.DB.QueryContext(ctx, script)
	comments := []entity.Comment{}
	if err != nil {
		return comments, err
	}

	defer rows.Close()
	for rows.Next() {
		comment := entity.Comment{}
		rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		comments = append(comments, comment)
	}

	return comments, nil
}
