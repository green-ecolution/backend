package mapper

import (
	"github.com/akatranlp/sentinel/account"
	"github.com/google/uuid"
	sqlc "github.com/green-ecolution/backend/internal/auth/storage/_sqlc"
)

// goverter:converter
// goverter:extend UUIDToUserID
type InternalAuthUserRepoMapper interface {
	// goverter:map ID UserID
	FromSql(src *sqlc.AuthUser) *account.User
	FromSqlList(src []*sqlc.AuthUser) []*account.User
}

// goverter:converter
// goverter:extend github.com/green-ecolution/backend/internal/utils:TimeToTime
// goverter:extend UUIDToUserID
type InternalAuthAccountRepoMapper interface {
	// goverter:map . AccountID
	FromSql(src *sqlc.AuthAccount) *account.Account

	FromSqlList(src []*sqlc.AuthAccount) []*account.Account
}

func UUIDToUserID(id uuid.UUID) account.UserID {
	return account.UserID(id.String())
}
