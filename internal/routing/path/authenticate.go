package path

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"

	users "github.com/roguexray007/loan-app/internal/app/users"
	"github.com/roguexray007/loan-app/internal/routing/tenant"
)

type Authenticate struct {
	path   *Endpoint
	tenant *tenant.Tenant
	user   string
	pass   string
}

func (auth *Authenticate) Request(req *http.Request) error {
	auth.tenant = tenant.From(req.Context())
	if auth.tenant == nil {
		return fmt.Errorf("auth: tenant is required")
	}

	if auth.path.IsPublic() {
		return auth.Public(req.Context())
	}

	user, pass, ok := req.BasicAuth()
	if !ok {
		return fmt.Errorf("auth: basic auth is required")
	}
	auth.user = user
	auth.pass = pass

	switch {
	case auth.path.IsPrivate():
		return auth.Private(req.Context())
	default:
		return fmt.Errorf("Path: auth not supported, check with admin")
	}
}

func (auth *Authenticate) Private(ctx context.Context) error {
	user := users.User{}
	err := users.GetRepo().FindByUsername(ctx, &user, auth.user)
	if err != nil {
		return fmt.Errorf("auth: unable to find user")
	}

	passSame := subtle.ConstantTimeCompare([]byte(auth.pass), []byte(user.Pass))
	if passSame != 1 {
		return fmt.Errorf("auth: wrong pass")
	}
	auth.tenant.SetUser(&user)
	auth.tenant.SetTenantType(tenant.UserType)

	return nil
}

func (auth *Authenticate) Public(ctx context.Context) error {
	return nil
}
