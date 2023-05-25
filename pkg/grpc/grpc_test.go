package grpc_test

import (
	"context"
	"github.com/reactivejson/users-svc/internal/domain"
	"github.com/reactivejson/users-svc/internal/notifier"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/reactivejson/users-svc/internal/app"
	"github.com/reactivejson/users-svc/pkg/grpc"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */
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
	if m.users == nil {
		m.users = make(map[string]*domain.User)
	}
	m.users[user.ID] = user
	return m.updateError
}

func (m *mockStorage) Delete(userID string) error {
	if m.users == nil {
		m.users = make(map[string]*domain.User)
	}
	delete(m.users, userID)
	return m.deleteError
}

func (m *mockStorage) FindByCountry(country string, page, pageSize int) ([]*domain.User, error) {
	if m.users == nil {
		m.users = make(map[string]*domain.User)
	}
	v := make([]*domain.User, 0, len(m.users))
	for _, value := range m.users {
		v = append(v, value)
	}
	return v, nil
}

type mockNotifier struct {
	notifyError error
}

func (m *mockNotifier) NotifyUserChange(event notifier.UserEventType, user *domain.User) error {
	return m.notifyError
}

func TestUserService_AddUser(t *testing.T) {
	mockStorage := &mockStorage{}
	mockNotifier := &mockNotifier{}
	// Create a new instance of UserServiceServerImpl
	userService := app.NewUserService(mockStorage, mockNotifier) // Assuming you have a constructor function to create UserService instance
	server := grpc.UserServiceServerImpl{
		UserService: userService,
	}

	// Define the test case
	testCases := []struct {
		name     string
		request  *grpc.UserRequest
		expected *grpc.UserResponse
		err      error
	}{
		{
			name: "AddUser_Success",
			request: &grpc.UserRequest{
				Id:        "1",
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Password:  "password",
				Email:     "john.doe@example.com",
				Country:   "USA",
			},
			expected: &grpc.UserResponse{
				Id: "1",
			},
			err: nil,
		},
		// Add more test cases for different scenarios
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Make the gRPC method call
			response, err := server.AddUser(context.Background(), tc.request)

			// Assertions
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, response)
		})
	}
}

func TestUserService_GetUsers(t *testing.T) {
	mockStorage := &mockStorage{}
	mockNotifier := &mockNotifier{}
	// Create a new instance of UserServiceServerImpl
	userService := app.NewUserService(mockStorage, mockNotifier) // Assuming you have a constructor function to create UserService instance
	server := grpc.UserServiceServerImpl{
		UserService: userService,
	}
	user := &grpc.UserRequest{
		Id:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john.doe@example.com",
		Country:   "USA",
		// Add more user entries as per your test scenario
	}

	// Define the test case
	testCases := []struct {
		name     string
		request  *grpc.GetUsersRequest
		expected *grpc.GetUsersResponse
		err      error
	}{
		{
			name: "GetUsers_Success",
			request: &grpc.GetUsersRequest{
				Country:  "USA",
				Page:     1,
				PageSize: 10,
			},
			expected: &grpc.GetUsersResponse{
				Users: []*grpc.UserRequest{
					user,
					// Add more user entries as per your test scenario
				},
			},
			err: nil,
		},
		// Add more test cases for different scenarios
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server.AddUser(context.Background(), user)
			// Make the gRPC method call
			response, err := server.GetUsers(context.Background(), tc.request)

			// Assertions
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, response)
		})
	}
}

func TestUserService_RemoveUser(t *testing.T) {
	mockStorage := &mockStorage{}
	mockNotifier := &mockNotifier{}
	// Create a new instance of UserServiceServerImpl
	userService := app.NewUserService(mockStorage, mockNotifier) // Assuming you have a constructor function to create UserService instance
	server := grpc.UserServiceServerImpl{
		UserService: userService,
	}

	// Define the test case
	testCases := []struct {
		name     string
		request  *grpc.UserRequest
		expected *grpc.UserResponse
		err      error
	}{
		{
			name: "RemoveUser_Success",
			request: &grpc.UserRequest{
				Id: "1",
			},
			expected: &grpc.UserResponse{},
			err:      nil,
		},
		// Add more test cases for different scenarios
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Make the gRPC method call
			response, err := server.RemoveUser(context.Background(), tc.request)

			// Assertions
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, response)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	mockStorage := &mockStorage{}
	mockNotifier := &mockNotifier{}
	// Create a new instance of UserServiceServerImpl
	userService := app.NewUserService(mockStorage, mockNotifier) // Assuming you have a constructor function to create UserService instance
	server := grpc.UserServiceServerImpl{
		UserService: userService,
	}

	// Define the test case
	testCases := []struct {
		name     string
		request  *grpc.UserRequest
		expected *grpc.UserResponse
		err      error
	}{
		{
			name: "UpdateUser_Success",
			request: &grpc.UserRequest{
				Id:        "1",
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Password:  "password",
				Email:     "john.doe@example.com",
				Country:   "USA",
				// Update other fields as per your test scenario
			},
			expected: &grpc.UserResponse{},
			err:      nil,
		},
		// Add more test cases for different scenarios
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Make the gRPC method call
			response, err := server.UpdateUser(context.Background(), tc.request)

			// Assertions
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, response)
		})
	}
}
