package tenant

import (
	"context"
	"errors"
)

type tenantcontextKeyType string

var tenantContextKey tenantcontextKeyType = "tenant_context_key"

type user interface {
	GetID() int64
}

type admin interface {
	GetID() int64
}

type Tenant struct {
	id      string
	tntType TntType
	admin   admin
	user    user
}

func New() *Tenant {
	return &Tenant{tntType: VacantType}
}

func (t *Tenant) SetID(id string) *Tenant {
	t.id = id
	return t
}

func (t *Tenant) SetUser(user user) *Tenant {
	t.user = user
	return t
}

func (t *Tenant) SetAdmin(admin admin) *Tenant {
	t.admin = admin
	return t
}

func (t *Tenant) SetTenantType(val TntType) *Tenant {
	t.tntType = val
	return t
}

func (t *Tenant) IsUser() bool {
	return t.tntType == UserType
}

func (t *Tenant) IsAdmin() bool {
	return t.tntType == AdminType
}

func (t *Tenant) User() user {
	return t.user
}

func (t *Tenant) Admin() admin {
	return t.admin
}

func Attach(ctx context.Context, tnt *Tenant) (context.Context, error) {
	if tnt == nil {
		return nil, errors.New("tenant is nil")
	}
	return context.WithValue(ctx, tenantContextKey, tnt), nil
}

// From will get the attached tenant from context
func From(ctx context.Context) *Tenant {
	if tnt, ok := ctx.Value(tenantContextKey).(*Tenant); ok {
		return tnt
	}
	return nil
}
