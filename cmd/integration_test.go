//go:build integration
// +build integration

package main

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */
import (
	"context"
	"fmt"
	conf "github.com/reactivejson/users-svc/cmd/app"
	"github.com/reactivejson/users-svc/internal/app"
	"github.com/reactivejson/users-svc/internal/notifier"
	"github.com/reactivejson/users-svc/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"testing"

	"github.com/reactivejson/users-svc/pkg/grpc"
	"github.com/stretchr/testify/assert"
)

func TestUserServiceServerImpl_TestIntegration(t *testing.T) {

	cfg := conf.SetupEnvConfig()
	fmt.Printf("KafkaBrokers %s\n", cfg.KafkaBrokers)
	fmt.Printf("GrpcPort %s\n", cfg.GrpcPort)

	repo := setupDB()
	kafka := setupKafka(cfg)

	userService := app.NewUserService(repo, kafka)
	server := grpc.UserServiceServerImpl{
		UserService: userService,
	}

	// Test case 1: Successful user addition
	req := &grpc.UserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Password:  "password",
		Email:     "johndoe@example.com",
		Country:   "US",
	}

	addResp, err := server.AddUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, addResp)
	assert.NotEmpty(t, addResp.Id)

	// Test case 2: Error in user addition (duplicate email)
	req.FirstName = "something"

	resp, err := server.AddUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	// Test case 3: Successful retrieval of users
	getUsersReq := &grpc.GetUsersRequest{
		Country:  "US",
		Page:     1,
		PageSize: 10,
	}
	usersResp, err := server.GetUsers(context.Background(), getUsersReq)

	assert.NoError(t, err)
	assert.NotNil(t, usersResp)
	assert.Len(t, usersResp.Users, 1) // Assuming we added one user in the previous test case

	// Test case 4: Error in user retrieval (invalid page size)
	getUsersReq.PageSize = -1

	usersResp, err = server.GetUsers(context.Background(), getUsersReq)

	assert.Error(t, err)
	assert.Nil(t, usersResp)

	// Test case 5: Successful removal of user
	removeUserReq := &grpc.UserRequest{
		Id: addResp.Id, // Use the ID generated from the successful user addition
	}

	removeUserResp, err := server.RemoveUser(context.Background(), removeUserReq)

	assert.NoError(t, err)
	assert.NotNil(t, removeUserResp)

	// Test case 6: Error in user removal (non-existent user)
	//removeUserReq.Id = "non_existent_id"

	//removeUserResp, err = server.RemoveUser(context.Background(), removeUserReq)

	//assert.Error(t, err)
	//assert.Nil(t, removeUserResp)

}

func setupDB() repository.UserRepository {

	// Initialize repository
	// Initialize the database connection
	dsn := "host=localhost user=postgres password=postgres dbname=user_svc port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	/*err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("failed to auto-migrate user table: %v", err)
	}*/

	return repository.NewUserRepository(db)
}

func setupKafka(cfg *conf.EnvConfig) *notifier.KafkaNotifier {
	userNotifier, err := notifier.NewKafkaNotifier(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}
	return userNotifier
}
