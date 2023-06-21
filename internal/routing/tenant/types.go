package tenant

type TntType string

const (
	VacantType TntType = "vacant"
	AdminType  TntType = "admin"
	UserType   TntType = "user"
)
