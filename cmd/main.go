// main.go
package main

import (
	"github.com/reactivejson/usr-svc/internal/domain"
	"github.com/reactivejson/usr-svc/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/reactivejson/usr-svc/internal/app"
	grpcc "github.com/reactivejson/usr-svc/pkg/grpc"
)

const (
	grpcPort     = ":50051"
	kafkaBrokers = "localhost:9092" // Modify with your Kafka brokers
	kafkaTopic   = "user_events"    // Modify with your Kafka topic
)

func main() {
	// Create a new instance of UserService
	// Initialize repository
	// Initialize the database connection
	dsn := "host=localhost user=your_user password=your_password dbname=your_database port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("failed to auto-migrate user table: %v", err)
	}

	userRepository := repository.NewUserRepository(db)

	//userRepo := repository.NewInMemoryUserStorage()
	userNotifier := &mockNotifier{}

	//userNotifier := &app.MockUserNotifier{} // Replace with actual implementation
	userService := app.NewUserService(userRepository, userNotifier)

	// Create a new gRPC server
	server := grpc.NewServer()

	// Register the UserServiceServer with the server
	grpcc.RegisterUserServiceServer(server, &grpcc.UserServiceServerImpl{
		UserService: userService,
	})

	// Start listening on a TCP port
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Serve the gRPC server
	log.Println("gRPC server is running...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type mockNotifier struct {
	notifyError error
}

func (m *mockNotifier) NotifyUserChange(user *domain.User) error {
	return m.notifyError
}
