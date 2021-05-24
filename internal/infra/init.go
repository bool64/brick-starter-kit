package infra

import (
	"database/sql"
	"fmt"
	"github.com/vearutop/gooselite"

	"github.com/Masterminds/squirrel"
	"github.com/bool64/brick"
	"github.com/bool64/brick-template/internal/domain/greeting"
	"github.com/bool64/brick-template/internal/infra/schema"
	"github.com/bool64/brick-template/internal/infra/service"
	"github.com/bool64/brick-template/internal/infra/storage"
	"github.com/bool64/brick/database"
	"github.com/bool64/brick/jaeger"
	"github.com/bool64/sqluct"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/swaggest/rest/response/gzip"
	"github.com/vearutop/gooselite/iofs"
)

// NewServiceLocator creates application service locator.
func NewServiceLocator(cfg service.Config) (*service.Locator, error) {
	bl, err := brick.NewBaseLocator(cfg.BaseConfig)
	if err != nil {
		return nil, err
	}

	l := &service.Locator{BaseLocator: bl}

	if err = jaeger.Setup(cfg.Jaeger, cfg.ServiceName, bl); err != nil {
		return nil, err
	}

	schema.SetupOpenapiCollector(l.OpenAPI)

	l.HTTPServerMiddlewares = append(l.HTTPServerMiddlewares, gzip.Middleware)

	if err = setupDatabase(l, cfg.Database); err != nil {
		return nil, err
	}

	l.GreetingMakerProvider = &storage.GreetingSaver{
		Upstream: &greeting.SimpleMaker{},
		Storage:  l.Storage,
	}

	return l, nil
}

func setupDatabase(l *service.Locator, cfg database.Config) error {
	c, err := mysql.ParseDSN(cfg.DSN)
	if err != nil {
		return err
	}

	conn, err := mysql.NewConnector(c)
	if err != nil {
		return err
	}

	l.Storage, err = database.SetupStorage(cfg, l.CtxdLogger(), "mysql", conn, storage.Migrations)
	if err != nil {
		return err
	}

	conn = database.WithTracing(conn)
	conn = database.WithQueriesLogging(conn, l.CtxdLogger())

	db := sql.OpenDB(conn)
	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	st := sqluct.NewStorage(sqlx.NewDb(sql.OpenDB(conn), "AAAAmSSDyDDDsql"))

	st.Format = squirrel.Question
	st.IdentifierQuoter = sqluct.QuoteBackticks

	if cfg.InitConn {
		if err = st.DB().Ping(); err != nil {
			return fmt.Errorf("failed to ping database: %w", err)
		}
	}

	if cfg.ApplyMigrations {
		if err := gooselite.SetDialect("mysql"); err != nil {
			return err
		}

		// Apply migrations.
		if err := iofs.Up(db, storage.Migrations, "migrations"); err != nil {
			return fmt.Errorf("failed to run up migrations: %v", err)
		}
	}

	l.Storage = st

	return nil
}
