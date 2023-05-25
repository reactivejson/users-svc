package grpc

import (
	"context"
	"errors"
	"github.com/reactivejson/users-svc/internal/app"
	"github.com/reactivejson/users-svc/internal/domain"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

// UserServiceServerImpl implements the gRPC UserServiceServer interface.
type UserServiceServerImpl struct {
	UserService *app.UserService
}

// AddUser implements the AddUser gRPC method.
func (s *UserServiceServerImpl) AddUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	user := &domain.User{
		ID:        req.Id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Password:  req.Password,
		Email:     req.Email,
		Country:   req.Country,
	}

	err := s.UserService.AddUser(user)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		Id: user.ID,
	}, nil
}

// GetUsers implements the GetUsers gRPC method.
func (s *UserServiceServerImpl) GetUsers(ctx context.Context, req *GetUsersRequest) (*GetUsersResponse, error) {
	if req.PageSize < 0 {
		return nil, errors.New("Error in user retrieval (invalid page size)")
	}
	users, err := s.UserService.GetUsers(req.Country, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var responseUsers []*UserRequest
	for _, user := range users {
		responseUser := &UserRequest{
			Id:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Country:   user.Country,
		}
		responseUsers = append(responseUsers, responseUser)
	}

	return &GetUsersResponse{
		Users: responseUsers,
	}, nil
}

// RemoveUser implements the DeleteUser gRPC method.
func (s *UserServiceServerImpl) RemoveUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	err := s.UserService.DeleteUser(req.Id)
	if err != nil {
		return nil, err
	}

	return &UserResponse{}, nil
}

// UpdateUser implements the UpdateUser gRPC method.
func (s *UserServiceServerImpl) UpdateUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	user := &domain.User{
		ID:        req.Id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Password:  req.Password,
		Email:     req.Email,
		Country:   req.Country,
	}

	err := s.UserService.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return &UserResponse{}, nil
}

// UpdateUser implements the UpdateUser gRPC method.
func (s *UserServiceServerImpl) mustEmbedUnimplementedUserServiceServer() {}
