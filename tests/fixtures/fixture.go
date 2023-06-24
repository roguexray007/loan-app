package fixtures

import (
	users "github.com/roguexray007/loan-app/internal/app/users"
	"github.com/roguexray007/loan-app/internal/routing/tenant"
)

var (
	LoanUser = &users.User{
		Username: "loanuser",
		Pass:     "loanpass",
	}
)

func Init() {
	InitUser()
}

func InitUser() {
	UserCreate(LoanUser)
}

func GetLoanUserTnt() *tenant.Tenant {
	LoanUserTnt := &tenant.Tenant{}
	LoanUserTnt.SetTenantType(tenant.UserType).SetUser(LoanUser)
	return LoanUserTnt
}
