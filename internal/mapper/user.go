package mapper

import (
	"fmt"

	"github.com/fatalistix/golang-jwt-server/internal/domain/model"
	"github.com/fatalistix/golang-jwt-server/internal/infrastructure/persistence/entity"
	"gorm.io/gorm"
)

type UserMapper struct {
}

func MakeUserMapper() UserMapper {
	return UserMapper{}
}

func (m UserMapper) ToEntity(u model.User) entity.User {
	roles := make([]entity.Role, 0, len(u.Roles))
	for _, v := range u.Roles {
		roles = append(roles, entity.Role{
			Name: v.String(),
		})
	}

	return entity.User{
		Model: gorm.Model{
			ID: u.Id,
		},
		Username: u.Username,
		Password: u.Password,
		Roles:    roles,
	}
}

func (m UserMapper) FromEntityToDomain(u entity.User) (model.User, error) {
	const op = "mapper.UserMapper.FromEntityToDomain"

	roles := make([]model.Role, 0, len(u.Roles))
	for _, v := range u.Roles {
		modelRole, err := model.MakeRole(v.Name)
		if err != nil {
			return model.User{}, fmt.Errorf("%s: %w", op, err)
		}

		roles = append(roles, modelRole)
	}

	return model.User{
		Id:       u.ID,
		Username: u.Username,
		Password: u.Password,
		Roles:    roles,
	}, nil
}
