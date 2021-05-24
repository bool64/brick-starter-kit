package storage

import (
	"context"
	"time"

	"github.com/bool64/brick-template/internal/domain/greeting"
	"github.com/bool64/ctxd"
	"github.com/bool64/sqluct"
)

type GreetingSaver struct {
	Upstream greeting.Maker
	Storage  *sqluct.Storage
}

const GreetingsTable = "greetings"

type GreetingRow struct {
	ID        int       `db:"id,omitempty"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

func (gs *GreetingSaver) Hello(ctx context.Context, params greeting.Params) (string, error) {
	g, err := gs.Upstream.Hello(ctx, params)
	if err != nil {
		return g, err
	}

	q := gs.Storage.InsertStmt(GreetingsTable, GreetingRow{
		Message:   g,
		CreatedAt: time.Now(),
	})

	_, err = gs.Storage.Exec(ctx, q)
	if err != nil {
		return "", ctxd.WrapError(ctx, err, "failed to store greeting")
	}

	return g, nil
}

// GreetingMaker implements service provider.
func (s *GreetingSaver) GreetingMaker() greeting.Maker {
	return s
}
