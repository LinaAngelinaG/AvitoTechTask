package postgresql

import (
	"AvitoTechTask/internal/configuration"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type Client interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, attempsVal int, cnf *configuration.Config) (pool *pgxpool.Pool, err error) {
	psqlConnection := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cnf.Storage.Username, cnf.Storage.Password, cnf.Storage.Host, cnf.Storage.Port, cnf.Storage.Database)
	for attempsVal > 0 {
		ctx, endC := context.WithTimeout(ctx, time.Second*5)
		defer endC()
		pool, err = pgxpool.Connect(ctx, psqlConnection)
		if err != nil {
			time.Sleep(time.Second * 5)
			attempsVal--
			continue
		}
		break
	}
	if err != nil {
		log.Fatal("impossible to connect to postgresql")
	}
	return pool, nil
}
