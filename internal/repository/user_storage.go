// internal/repository/user_storage.go
package repository

import (
	"github.com/pkg/errors"

	"github.com/reactivejson/usr-svc/internal/domain"
)

// InMemoryUserStorage is an in-memory implementation of UserStorage.
type InMemoryUserStorage struct {
	users []*domain.User
}

// NewInMemoryUserStorage creates a new instance of InMemoryUserStorage.
func NewInMemoryUserStorage() *InMemoryUserStorage {
	return &InMemoryUserStorage{
		users: make([]*domain.User, 0),
	}
}

// Save saves a user.
func (s *InMemoryUserStorage) Save(user *domain.User) error {
	s.users = append(s.users, user)
	return nil
}

// Update updates a user.
func (s *InMemoryUserStorage) Update(user *domain.User) error {
	for i, u := range s.users {
		if u.ID == user.ID {
			s.users[i] = user
			return nil
		}
	}
	return errors.New("user not found")
}

// Delete deletes a user.
func (s *InMemoryUserStorage) Delete(userID string) error {
	for i, user := range s.users {
		if user.ID == userID {
			// Remove the user from the slice
			s.users = append(s.users[:i], s.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

// FindByCountry finds users by country.
func (s *InMemoryUserStorage) FindByCountry(country string, page, pageSize int) ([]*domain.User, error) {
	// Implement pagination logic here
	// Note: This is a simplified example; you should implement proper pagination based on the page and pageSize parameters
	return s.users, nil
}
