// internal/app/user_service.go
package app

import (
	"github.com/pkg/errors"
	"github.com/reactivejson/usr-svc/internal/domain"
	"github.com/reactivejson/usr-svc/internal/notifier"
	"github.com/reactivejson/usr-svc/internal/repository"
)

// UserService represents the user service.
type UserService struct {
	Storage  repository.UserRepository
	Notifier notifier.UserNotifier
}

// NewUserService creates a new instance of UserService.
func NewUserService(repository repository.UserRepository, notifier notifier.UserNotifier) *UserService {

	return &UserService{
		Storage:  repository,
		Notifier: notifier,
	}
}

// AddUser adds a new user.
func (s *UserService) AddUser(user *domain.User) error {
	err := s.Storage.Save(user)
	if err != nil {
		return errors.Wrap(err, "failed to save user")
	}

	err = s.Notifier.NotifyUserChange(user)
	if err != nil {
		// Handle error while notifying
	}

	return nil
}

// UpdateUser updates an existing user.
func (s *UserService) UpdateUser(user *domain.User) error {
	err := s.Storage.Update(user)
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	err = s.Notifier.NotifyUserChange(user)
	if err != nil {
		// Handle error while notifying
	}

	return nil
}

// DeleteUser deletes a user.
func (s *UserService) DeleteUser(userID string) error {
	err := s.Storage.Delete(userID)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	// Notifying is not necessary for deletion

	return nil
}

// GetUsers retrieves a paginated list of users filtered by criteria.
func (s *UserService) GetUsers(country string, page, pageSize int) ([]*domain.User, error) {
	users, err := s.Storage.FindByCountry(country, page, pageSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users")
	}

	return users, nil
}
