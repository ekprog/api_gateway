package core

type AccessRole int32

const (
	RoleGuest      AccessRole = 0
	RoleUser       AccessRole = 1
	RoleSuperAdmin AccessRole = 10
)
