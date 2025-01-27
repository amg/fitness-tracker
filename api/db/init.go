package db

import (
	"context"
	gendb "fitness-tracker/db/generated"
	"fitness-tracker/env"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Queries *gendb.Queries
	Pool    *pgxpool.Pool
}

// InitDB initializes a Unix socket connection pool for
// a Cloud SQL instance of MySQL.
func InitDB(config env.Config) (*DB, error) {
	var (
		dbUser         = config.SecEnv.PostgresUser()
		dbPwd          = config.SecEnv.PostgresPassword()
		dbName         = config.SecEnv.PostgresDbName()
		unixSocketPath = config.Env.PostgresUrl() // e.g. '/project:region:instance'
	)

	// in dev we got postgres running in a docker but in staging GCP Cloud SQL is taking care of it
	var connString string
	switch config.Env.(type) {
	case env.DevEnv:
		connString = fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable",
			dbUser,
			dbPwd,
			unixSocketPath,
			dbName)
	case env.StagingEnv:
		// knowing right connection string is a bit painful so use env at container level when playing with this
		//  for now left split to remain consistent with local docker
		connString = fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=disable", unixSocketPath, dbUser, dbPwd, dbName)
	default:
		panic("db: unsupported env")
	}

	cx := context.Background()
	pool, err := pgxpool.New(cx, connString)
	if err != nil {
		return nil, fmt.Errorf("db.open: %w;", err)
	}

	return &DB{Queries: gendb.New(pool), Pool: pool}, nil
}
