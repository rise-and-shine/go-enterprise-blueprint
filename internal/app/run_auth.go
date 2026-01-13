package app

import (
	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/pg"
)

func RunAuth(cfg Config) error {
	app := newAppWithConfig(cfg)
	defer app.shutdown()

	err := app.initObservability(cfg.Auth.Name, cfg.Auth.Version)
	if err != nil {
		return errx.Wrap(err)
	}

	dbConn, err := pg.NewBunDB(cfg.Postgres)
	if err != nil {
		return errx.Wrap(err)
	}
	app.dbCloser = dbConn.Close

	// TODO...

	return nil
}
