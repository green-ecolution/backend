package mapper

import (
	"github.com/akatranlp/sentinel/token"
	sqlc "github.com/green-ecolution/backend/internal/auth/storage/_sqlc"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/backend/internal/utils:TimeToTime
type InternalTokenSessionRepoMapper interface {
	// goverter:map RefreshJti RefreshJTI
	FromSql(src *sqlc.TokenSession) *token.Session

	FromSqlList(src []*sqlc.TokenSession) []*token.Session
}
