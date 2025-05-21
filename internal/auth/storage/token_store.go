package storage

import (
	"context"

	"github.com/akatranlp/sentinel/token"
	"github.com/akatranlp/sentinel/token/base_store"
	"github.com/green-ecolution/backend/internal/auth/mapper"
	sqlc "github.com/green-ecolution/backend/internal/auth/storage/_sqlc"
)

var _ basestore.Repository = (*TokenSessionRepository)(nil)
var _ basestore.SessionDeleter = (*TokenSessionRepository)(nil)

type TokenSessionRepository struct {
	store *sqlc.Queries
	TokenSessionMappers
}

type TokenSessionMappers struct {
	mapper mapper.InternalTokenSessionRepoMapper
}

func NewTokenSessionMappers(mapper mapper.InternalTokenSessionRepoMapper) TokenSessionMappers {
	return TokenSessionMappers{
		mapper: mapper,
	}
}

func NewTokenSessionRepository(queries *sqlc.Queries, mappers TokenSessionMappers) *TokenSessionRepository {
	return &TokenSessionRepository{
		store:               queries,
		TokenSessionMappers: mappers,
	}
}

func (r *TokenSessionRepository) CreateSession(ctx context.Context, session token.Session) (token.Session, error) {
	sess, err := r.store.CreateSession(ctx, &sqlc.CreateSessionParams{
		SessionID:  session.SessionID,
		Expiry:     session.Expiry,
		RefreshJti: session.RefreshJTI,
	})
	if err != nil {
		return token.Session{}, err
	}

	return *r.mapper.FromSql(sess), nil
}

func (r *TokenSessionRepository) GetSessionByID(ctx context.Context, sid string) (token.Session, error) {
	sess, err := r.store.GetSessionById(ctx, sid)
	if err != nil {
		return token.Session{}, err
	}
	return *r.mapper.FromSql(sess), nil
}

func (r *TokenSessionRepository) UpdateSession(ctx context.Context, session token.Session) error {
	return r.store.UpdateSession(ctx, &sqlc.UpdateSessionParams{
		SessionID:  session.SessionID,
		Expiry:     session.Expiry,
		RefreshJti: session.RefreshJTI,
	})
}

func (r *TokenSessionRepository) DeleteSessionByID(ctx context.Context, sid string) error {
	return r.store.DeleteSession(ctx, sid)
}

func (r *TokenSessionRepository) DeleteSessionsAfterExpiry(ctx context.Context) error {
	return r.store.DeleteSessionAfterExpiry(ctx)
}
