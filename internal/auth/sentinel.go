package auth

import (
	"context"
	"os"
	"time"

	accountbasestore "github.com/akatranlp/sentinel/account/base_store"
	"github.com/akatranlp/sentinel/openid"
	"github.com/akatranlp/sentinel/openid/enums"
	"github.com/akatranlp/sentinel/provider"
	tokenbasestore "github.com/akatranlp/sentinel/token/base_store"
	"github.com/akatranlp/sentinel/utils"
	"github.com/green-ecolution/backend/internal/storage"
)

const (
	clientID     = ""
	clientSecret = ""
)

func CreateSentinel(ctx context.Context, repos *storage.Repository) (*openid.IdentitiyProvider, error) {
	f, err := os.Open("jwk.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	userStore := accountbasestore.NewBaseUserStore(repos.AuthUserRepository)
	tokenStore, err := tokenbasestore.NewBaseTokenStore(repos.TokenRepository)
	if err != nil {
		return nil, err
	}

	ip, err := openid.NewIdentityProvider(
		"/auth",
		userStore,
		tokenStore,
		repos.SessionStore,
		openid.WithClients(openid.ClientRegistration{
			ClientID:            clientID,
			ClientSecret:        clientSecret,
			TokenExchangeSecret: "",
			Scope:               enums.ScopeValues(),
			RedirectURIs: []string{
				"http://localhost/callback",
				"https://oidcdebugger.com/debug",
			},
		}),
		openid.WithProviders(utils.Must(provider.ProviderFactory(provider.FactoryParams{
			Type:    "gitlab",
			Name:    "GitLab",
			Slug:    "gitlab",
			BaseURL: "https://gitlab.git-classrooms.dev",
			// TODO: Use a different ClientID and Secret for the ip and not gitlabs
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}))),
		openid.WithAccessTokenExpiration(15*time.Minute),
		openid.WithRefreshTokenExpiration(7*24*time.Hour),
		openid.WithSessionUnAuthedLifeTime(10*time.Minute),
		openid.WithSessionAuthedLifeTime(30*24*time.Hour),
		openid.WithSessionIdleTimeout(7*24*time.Hour),
		openid.WithSigningKeyReader(f),
		openid.WithAppURL("../"),
	)
	if err != nil {
		return nil, err
	}

	return ip, nil
}
