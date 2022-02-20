package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func initDb(filename string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", filename))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Migrate
	if true {
		files, err := os.ReadDir("../migrations")

		if err != nil {
			return nil, err
		}

		for _, file := range files {
			data, err := os.ReadFile("../migrations/" + file.Name())
			if err != nil {
				return nil, err
			}

			_, err = db.Exec(string(data))

			if err != nil {
				return nil, err
			}
		}
	}

	return db, nil
}

func (app *application) TodoGetByUserId(ctx context.Context, userId uint) ([]Todo, error) {
	rows, err := app.db.QueryContext(ctx, `
	SELECT "id", "checked", "text", "user_id"
	FROM "todos"
	WHERE "user_id" = ?
	`, userId)

	if err != nil {
		return nil, err
	}

	var result []Todo

	for rows.Next() {
		var todo Todo
		if err = rows.Scan(&todo.ID, &todo.Checked, &todo.Text, &todo.UserID); err != nil {
			return nil, err
		}

		result = append(result, todo)
	}

	return result, nil
}

func (app *application) TodoSetCheck(ctx context.Context, checked bool, id int) error {
	result, err := app.db.ExecContext(ctx, `
	UPDATE "todos"
	SET "checked" = ?
	WHERE "id" = ?
	`, checked, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New(fmt.Sprintf("todo %d not found", id))
	}

	return nil
}

func (app *application) TodoCreate(ctx context.Context, todo *Todo) error {

	result, err := app.db.ExecContext(ctx, `
	INSERT INTO "todos" ("checked", "text", "user_id")
	VALUES (?, ?, ?)
	`, todo.Checked, todo.Text, todo.UserID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	todo.ID = uint(id)

	return nil
}

func (app *application) UserExist(ctx context.Context, email string) (bool, error) {
	var id uint
	err := app.db.QueryRowContext(ctx, `
	SELECT "id"
	FROM "users"
	WHERE "users"."email" = ?
	`, email).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return id != 0, err
}

func (app *application) UserCreate(ctx context.Context, user *User) error {
	result, err := app.db.ExecContext(ctx, `
	INSERT INTO "users" ("email", "password_hash")
	VALUES (?, ?)
	`, user.Email, user.PasswordHash)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	user.ID = uint(id)

	return nil
}

func (app *application) UserFindByEmail(ctx context.Context, email string) (User, error) {
	var user User

	err := app.db.QueryRowContext(ctx, `
	SELECT "id", "email", "password_hash"
	FROM "users"
	WHERE "users"."email" = ?
	`, email).Scan(&user.ID, &user.Email, &user.PasswordHash)

	return user, err
}

func (app *application) SessionCreate(ctx context.Context, session *Session) error {
	result, err := app.db.ExecContext(ctx, `
	INSERT INTO "sessions" ("cookie", "email")
	VALUES (?, ?)
	`, session.Cookie, session.Email)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	session.ID = uint(id)

	return nil
}

func (app *application) SessionFind(ctx context.Context, cookie string) (Session, error) {
	var session Session

	err := app.db.QueryRowContext(ctx, `
	SELECT "id", "cookie", "email"
	FROM "sessions"
	WHERE "cookie" = ?
	`, cookie).Scan(&session.ID, &session.Cookie, &session.Email)

	return session, err
}
