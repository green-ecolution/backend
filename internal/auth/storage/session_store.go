package storage

import (
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewSessionStore(pool *pgxpool.Pool) (scs.Store, error) {
	return pgxstore.New(pool), nil
}
