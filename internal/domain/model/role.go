package model

import "fmt"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func Roles() []Role {
	return []Role{
		RoleAdmin, RoleUser,
	}
}

func (r Role) String() string {
	return string(r)
}

func MakeRole(role string) (Role, error) {
	const op = "domain.model.MakeRole"

	switch role {
	case string(RoleAdmin):
		return RoleAdmin, nil
	case string(RoleUser):
		return RoleUser, nil
	default:
		return "", fmt.Errorf("%s: unknown role: %v", op, role)
	}
}
