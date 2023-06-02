package dbrepo

import (
	"database/sql"

	"github.com/ishanshre/Book-Review-Platform/internals/config"
	"github.com/ishanshre/Book-Review-Platform/internals/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
