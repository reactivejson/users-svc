package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/reactivejson/users-svc/internal/app"
	"github.com/reactivejson/users-svc/internal/notifier"
	"github.com/reactivejson/users-svc/internal/repository"
	grpcc "github.com/reactivejson/users-svc/pkg/grpc"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

func setupLog() *log.Logger {
	return log.Default()
	//TODO proper logging setup
}

func SetupEnvConfig() *EnvConfig {

	cfg := &EnvConfig{}
	if err := envconfig.Process("", cfg); err != nil {
		fmt.Errorf("could not parse config: %w", err)
	}
	return cfg
}

func setupDB() setupFn {
	return func(c *Context) (err error) {
		if c.UserRepository == nil {
			// Initialize repository
			// Initialize the database connection
			dsn := "host=localhost user=postgres password=postgres dbname=user_svc port=5432 sslmode=disable"
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("failed to connect to database: %v", err)
			}

			/*err = db.AutoMigrate(&domain.User{})
			if err != nil {
				log.Fatalf("failed to auto-migrate user table: %v", err)
			}*/

			c.UserRepository = repository.NewUserRepository(db)
		}
		return nil
	}
}

func setupKafka() setupFn {
	return func(c *Context) (err error) {
		if c.userNotifier == nil {
			userNotifier, err := notifier.NewKafkaNotifier(c.cfg.KafkaBrokers, c.cfg.KafkaTopic)
			if err != nil {
				log.Fatal("Failed to create Kafka producer:", err)
			}
			c.userNotifier = userNotifier
		}
		return nil
	}
}

func setupGrpcServer() setupFn {
	return func(c *Context) (err error) {
		if c.grpcServer == nil {
			// Create a new gRPC server
			server := grpc.NewServer()

			// Register the UserServiceServer with the server
			grpcc.RegisterUserServiceServer(server, &grpcc.UserServiceServerImpl{
				UserService: c.UserService,
			})

			// Start listening on a TCP port
			listener, err := net.Listen("tcp", c.cfg.GrpcPort)
			if err != nil {
				log.Fatalf("Failed to listen: %v", err)
			}

			// Serve the gRPC server
			log.Println("gRPC server is running...")
			if err := server.Serve(listener); err != nil {
				log.Fatalf("Failed to serve: %v", err)
			}
			c.grpcServer = server
		}
		return nil
	}
}

func setupHttpServer() setupFn {
	return func(c *Context) (err error) {
		if c.httpServer == nil {
			// Set up HTTP server
			router := gin.Default()

			router.GET("/health", app.HealthHandler)
			router.POST("/users", c.UserService.AddUserHandler)
			router.PUT("/users/:id", c.UserService.UpdateUserHandler)
			router.DELETE("/users/:id", c.UserService.DeleteUserHandler)
			router.GET("/users", c.UserService.GetUsersHandler)

			c.httpServer = router
			log.Fatal(router.Run(":8080"))
		}
		return nil
	}
}

func setupUserService() setupFn {
	return func(c *Context) (err error) {
		if c.UserService == nil {
			c.UserService = app.NewUserService(c.UserRepository, c.userNotifier)
		}
		return nil
	}
}
