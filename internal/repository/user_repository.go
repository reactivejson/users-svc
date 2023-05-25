package repository

import (
	"github.com/google/uuid"
	model "github.com/reactivejson/users-svc/internal/domain"
	"gorm.io/gorm"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

type UserRepository interface {
	Save(user *model.User) error
	Update(user *model.User) error
	Delete(userID string) error
	FindByCountry(country string, page, pageSize int) ([]*model.User, error)
}
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Save(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Model(user).Updates(user).Error
}

func (r *userRepository) Delete(userID string) error {
	id, _ := uuid.Parse(userID)
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) FindByCountry(country string, page, pageSize int) ([]*model.User, error) {
	var users []*model.User
	offset := (page - 1) * pageSize

	err := r.db.Where("country = ?", country).Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
