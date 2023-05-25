package app

import (
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/reactivejson/users-svc/internal/domain"
	"github.com/reactivejson/users-svc/internal/notifier"
	"github.com/reactivejson/users-svc/internal/repository"
	"lukechampine.com/blake3"

	"log"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

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
	user.Password = hashPswd(user.Password)
	err := s.Storage.Save(user)
	if err != nil {
		return errors.Wrap(err, "failed to save user")
	}

	err = s.Notifier.NotifyUserChange(notifier.Created, user)
	if err != nil {
		log.Printf("Failed to send notification.")

	}

	return nil
}

// UpdateUser updates an existing user.
func (s *UserService) UpdateUser(user *domain.User) error {
	user.Password = hashPswd(user.Password)
	err := s.Storage.Update(user)
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	err = s.Notifier.NotifyUserChange(notifier.Updated, user)
	if err != nil {
		log.Printf("Failed to send notification.")
	}

	return nil
}

// DeleteUser deletes a user.
func (s *UserService) DeleteUser(userID string) error {
	err := s.Storage.Delete(userID)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	err = s.Notifier.NotifyUserChange(notifier.Deleted, &domain.User{ID: userID})
	if err != nil {
		log.Printf("Failed to send notification.")
	}
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

func hashPswd(password string) string {
	hash := blake3.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
