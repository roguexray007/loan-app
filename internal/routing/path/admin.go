package path

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/roguexray007/loan-app/internal/routing/tenant"
)

const (
	// this can be moved to a config or db. currently for simplicity hardcoded in code
	AdminUser = "admin"
	AdminPass = "adminpass"
	AdminID   = int64(100000)
)

type Admin struct {
	path *Endpoint
	user string
	ID   int64
}

func (a *Admin) GetID() int64 {
	return a.ID
}

func (a *Admin) Request(req *http.Request) error {
	tnt := tenant.From(req.Context())
	if tnt == nil {
		return fmt.Errorf("auth: tenant is required")
	}
	user, pass, ok := req.BasicAuth()
	if !ok {
		return fmt.Errorf("auth: no admin credentials")
	}
	err := a.Authenticate(user, pass)
	if err != nil {
		return err
	}
	// Mark tenant as admin only after successful validation
	tnt.SetTenantType(tenant.AdminType)
	tnt.SetAdmin(a)
	return nil
}

func (admin *Admin) Authenticate(user, pass string) error {
	if AdminUser == user {
		result := subtle.ConstantTimeCompare([]byte(AdminPass), []byte(pass))
		if result == 1 {
			admin.user = user
			admin.ID = AdminID
			return nil
		}
	}

	return fmt.Errorf("auth: invalid admin credentials")
}
