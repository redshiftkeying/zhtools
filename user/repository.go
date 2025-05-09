//go:generate sqlc generate

package user

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"
)

//go:embed schema.sql
var ddl string

func NewUserRepo(ctx context.Context) error {
	db, err := sql.Open("sqlite", "main.sqlite")
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
