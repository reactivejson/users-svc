// internal/app/user_service_test.go
package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/reactivejson/usr-svc/internal/domain"
)

type mockStorage struct {
	users       map[string]*domain.User
	saveError   error
	updateError error
	deleteError error
}

func (m *mockStorage) Save(user *domain.User) error {
	if m.users == nil {
		m.users = make(map[string]*domain.User)
	}
	m.users[user.ID] = user
	return m.saveError
}

func (m *mockStorage) Update(user *domain.User) error {
	m.users[user.ID] = user
	return m.updateError
}

func (m *mockStorage) Delete(userID string) error {
	delete(m.users, userID)
	return m.deleteError
}

func (m *mockStorage) FindByCountry(country string, page, pageSize int) ([]*domain.User, error) {
	v := make([]*domain.User, 0, len(m.users))
	for _, value := range m.users {
		v = append(v, value)
	}
	return v, nil
}

type mockNotifier struct {
	notifyError error
}

func (m *mockNotifier) NotifyUserChange(user *domain.User) error {
	return m.notifyError
}

func TestAddUser(t *testing.T) {
	user := &domain.User{ID: "1", FirstName: "Alice", LastName: "Bob", Country: "UK"}

	mockStorage := &mockStorage{}
	mockNotifier := &mockNotifier{}

	service := NewUserService(mockStorage, mockNotifier)

	err := service.AddUser(user)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockStorage.users))
	assert.Equal(t, user, mockStorage.users["1"])
}

func TestUpdateUser(t *testing.T) {
	user := &domain.User{ID: "1", FirstName: "Alice", LastName: "Bob", Country: "UK"}

	mockStorage := &mockStorage{}
	mockNotifier := &mockNotifier{}

	service := NewUserService(mockStorage, mockNotifier)

	service.AddUser(user)
	err := service.UpdateUser(&domain.User{ID: "1", FirstName: "Alice2", LastName: "Bob2", Country: "Finland"})

	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockStorage.users))
}

func TestDeleteUser(t *testing.T) {
	userID := "1"

	mockStorage := &mockStorage{}
	mockNotifier := &mockNotifier{}

	service := NewUserService(mockStorage, mockNotifier)

	err := service.DeleteUser(userID)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(mockStorage.users))
}

func TestGetUsers(t *testing.T) {
	usersMap := make(map[string]*domain.User)
	users := []*domain.User{
		{ID: "1", FirstName: "Alice", LastName: "Bob", Country: "UK"},
		{ID: "2", FirstName: "John", LastName: "Doe", Country: "UK"},
		{ID: "3", FirstName: "Jane", LastName: "Smith", Country: "US"},
	}
	usersMap["1"] = users[0]
	usersMap["2"] = users[1]
	usersMap["3"] = users[2]

	mockStorage := &mockStorage{users: usersMap}
	mockNotifier := &mockNotifier{}

	service := NewUserService(mockStorage, mockNotifier)

	result, err := service.GetUsers("UK", 1, 2)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(result))
}

func TestAddUser_StorageError(t *testing.T) {
	user := &domain.User{ID: "1", FirstName: "Alice", LastName: "Bob", Country: "UK"}

	mockStorage := &mockStorage{saveError: errors.New("storage error")}
	mockNotifier := &mockNotifier{}

	service := NewUserService(mockStorage, mockNotifier)

	err := service.AddUser(user)

	assert.Error(t, err)
}

func TestUpdateUser_StorageError(t *testing.T) {
	user := &domain.User{ID: "1", FirstName: "Alice", LastName: "Bob", Country: "UK"}

	mockStorage := &mockStorage{updateError: errors.New("storage error")}
	mockNotifier := &mockNotifier{}

	service := NewUserService(mockStorage, mockNotifier)

	service.AddUser(user)
	err := service.UpdateUser(&domain.User{ID: "1", FirstName: "Alice2", LastName: "Bob2", Country: "Finland"})

	assert.Error(t, err)
}

func TestDeleteUser_StorageError(t *testing.T) {
	userID := "1"

	mockStorage := &mockStorage{deleteError: errors.New("storage error")}
	mockNotifier := &mockNotifier{}

	service := NewUserService(mockStorage, mockNotifier)

	err := service.DeleteUser(userID)

	assert.Error(t, err)
}

func TestAddUser_NotifyError(t *testing.T) {
	user := &domain.User{ID: "1", FirstName: "Alice", LastName: "Bob", Country: "UK"}

	mockStorage := &mockStorage{}
	mockNotifier := &mockNotifier{notifyError: errors.New("notification error")}

	service := NewUserService(mockStorage, mockNotifier)

	err := service.AddUser(user)

	assert.NoError(t, err)
}

func TestUpdateUser_NotifyError(t *testing.T) {
	user := &domain.User{ID: "1", FirstName: "Alice", LastName: "Bob", Country: "UK"}

	mockStorage := &mockStorage{}
	mockNotifier := &mockNotifier{notifyError: errors.New("notification error")}

	service := NewUserService(mockStorage, mockNotifier)

	service.AddUser(user)
	err := service.UpdateUser(&domain.User{ID: "1", FirstName: "Alice2", LastName: "Bob2", Country: "Finland"})

	assert.NoError(t, err)
}

// Add similar tests for DeleteUser and GetUsers with NotifyError scenarios
