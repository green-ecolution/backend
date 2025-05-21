package storage

import (
	"context"

	"github.com/akatranlp/sentinel/account"
	"github.com/akatranlp/sentinel/account/base_store"
	"github.com/google/uuid"
	"github.com/green-ecolution/backend/internal/auth/mapper"
	sqlc "github.com/green-ecolution/backend/internal/auth/storage/_sqlc"
	"github.com/green-ecolution/backend/internal/utils"
)

var _ basestore.Repository = (*AuthUserRepository)(nil)
var _ basestore.UserIDAndPRoviderGetter = (*AuthUserRepository)(nil)

type AuthUserRepository struct {
	store *sqlc.Queries
	AuthUserMappers
}

type AuthUserMappers struct {
	userMapper    mapper.InternalAuthUserRepoMapper
	accountMapper mapper.InternalAuthAccountRepoMapper
}

func NewAuthUserMappers(userMapper mapper.InternalAuthUserRepoMapper, accountMapper mapper.InternalAuthAccountRepoMapper) AuthUserMappers {
	return AuthUserMappers{
		userMapper:    userMapper,
		accountMapper: accountMapper,
	}
}

func NewAuthUserRepository(queries *sqlc.Queries, mappers AuthUserMappers) *AuthUserRepository {
	return &AuthUserRepository{
		store:           queries,
		AuthUserMappers: mappers,
	}
}

func (r *AuthUserRepository) CreateUser(ctx context.Context, user account.User) (account.User, error) {
	userId, err := r.store.CreateAuthUser(ctx, &sqlc.CreateAuthUserParams{
		Name:          user.Name,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Picture:       user.Picture,
	})
	if err != nil {
		return account.User{}, err
	}
	user.UserID = account.UserID(userId.String())

	return user, nil
}

func (r *AuthUserRepository) GetUserByID(ctx context.Context, id account.UserID) (account.User, error) {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return account.User{}, err
	}

	user, err := r.store.GetAuthUserById(ctx, userID)
	if err != nil {
		return account.User{}, account.ErrUserNotFound
	}

	return *r.userMapper.FromSql(user), nil
}

func (r *AuthUserRepository) GetUserByAccountID(ctx context.Context, id account.AccountID) (account.User, error) {
	user, err := r.store.GetAuthUserByAccountId(ctx, &sqlc.GetAuthUserByAccountIdParams{
		Provider:   id.Provider,
		ProviderID: id.ProviderID,
	})
	if err != nil {
		return account.User{}, account.ErrUserNotFound
	}

	return *r.userMapper.FromSql(user), nil
}

func (r *AuthUserRepository) UpdateUser(ctx context.Context, id account.UserID, user account.User) error {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return err
	}

	return r.store.UpdateAuthUser(ctx, &sqlc.UpdateAuthUserParams{
		ID:            userID,
		Name:          user.Name,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Picture:       user.Picture,
	})
}

func (r *AuthUserRepository) CreateAccount(ctx context.Context, acc account.Account) (account.Account, error) {
	userID, err := uuid.Parse(string(acc.UserID))
	if err != nil {
		return account.Account{}, err
	}

	_, err = r.store.CreateAuthAccount(ctx, &sqlc.CreateAuthAccountParams{
		Provider:          acc.Provider,
		ProviderID:        acc.ProviderID,
		AccessToken:       acc.AccessToken,
		Expiry:            acc.Expiry,
		RefreshToken:      acc.RefreshToken,
		RefreshExpiry:     acc.RefreshExpiry,
		IDToken:           acc.IDToken,
		TokenType:         acc.TokenType,
		Email:             acc.Email,
		EmailVerified:     acc.EmailVerified,
		Name:              acc.Name,
		PreferredUsername: acc.PreferredUsername,
		Picture:           acc.Picture,
		Nickname:          acc.Nickname,
		Profile:           acc.Profile,
		UserID:            userID,
	})
	if err != nil {
		return account.Account{}, err
	}

	return acc, nil
}

func (r *AuthUserRepository) GetAccountByID(ctx context.Context, id account.AccountID) (account.Account, error) {
	acc, err := r.store.GetAuthAccountById(ctx, &sqlc.GetAuthAccountByIdParams{
		Provider:   id.Provider,
		ProviderID: id.ProviderID,
	})
	if err != nil {
		return account.Account{}, account.ErrAccountNotFound
	}
	return *r.accountMapper.FromSql(acc), nil
}

func (r *AuthUserRepository) GetAccountByUserIDAndProvider(ctx context.Context, id account.UserID, provider string) (account.Account, error) {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return account.Account{}, err
	}

	acc, err := r.store.GetAuthAccountByUserIdAndProvider(ctx, &sqlc.GetAuthAccountByUserIdAndProviderParams{
		UserID:   userID,
		Provider: provider,
	})
	if err != nil {
		return account.Account{}, account.ErrAccountNotFound
	}
	return *r.accountMapper.FromSql(acc), nil
}

func (r *AuthUserRepository) GetAccountsByUserID(ctx context.Context, id account.UserID) ([]account.Account, error) {
	userID, err := uuid.Parse(string(id))
	if err != nil {
		return nil, err
	}
	accs, err := r.store.GetAuthAccountsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	return utils.Map(r.accountMapper.FromSqlList(accs), func(acc *account.Account) account.Account {
		return *acc
	}), nil
}

func (r *AuthUserRepository) UpdateAccount(ctx context.Context, id account.AccountID, acc account.Account) error {
	return r.store.UpdateAuthAccount(ctx, &sqlc.UpdateAuthAccountParams{
		Provider:          acc.Provider,
		ProviderID:        acc.ProviderID,
		AccessToken:       acc.AccessToken,
		Expiry:            acc.Expiry,
		RefreshToken:      acc.RefreshToken,
		RefreshExpiry:     acc.RefreshExpiry,
		IDToken:           acc.IDToken,
		TokenType:         acc.TokenType,
		Email:             acc.Email,
		EmailVerified:     acc.EmailVerified,
		Name:              acc.Name,
		PreferredUsername: acc.PreferredUsername,
		Picture:           acc.Picture,
		Nickname:          acc.Nickname,
		Profile:           acc.Profile,
	})
}

func (r *AuthUserRepository) DeleteAccountByID(ctx context.Context, id account.AccountID) error {
	return r.store.DeleteAuthAccount(ctx, &sqlc.DeleteAuthAccountParams{
		Provider:   id.Provider,
		ProviderID: id.ProviderID,
	})
}
