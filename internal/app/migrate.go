package app

import (
	"go-enterprise-blueprint/migrations"

	"github.com/code19m/errx"
	"github.com/pressly/goose/v3"
)

const (
	dialect       = "postgres"
	migrationsDir = "."
	tableName     = "_migrations"
)

func (a *app) applyMigrations() error {
	goose.SetBaseFS(migrations.MigrationsFS)
	goose.SetTableName(tableName)

	err := goose.SetDialect(dialect)
	if err != nil {
		return errx.Wrap(err)
	}

	err = goose.Up(a.dbConn.DB, migrationsDir)
	return errx.Wrap(err)
}
