package app

import (
	"github.com/gin-gonic/gin"
	"github.com/reactivejson/users-svc/internal/app"
	"github.com/reactivejson/users-svc/internal/notifier"
	"github.com/reactivejson/users-svc/internal/repository"
	"google.golang.org/grpc"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

type EnvConfig struct {
	Port         uint   `envconfig:"PORT"               required:"false" default:"8080"`
	GrpcPort     string `envconfig:"GRPCPORT"               required:"false" default:":50051"`
	KafkaBrokers string `envconfig:"KAFKA_BOKERS"      required:"false" default:"localhost:9092"`
	KafkaTopic   string `envconfig:"KAFKA_TOPIC"      required:"false" default:"user_events"`
}

// Context is application's content
type Context struct {
	Closers []func()

	cfg            *EnvConfig
	UserRepository repository.UserRepository
	userNotifier   *notifier.KafkaNotifier
	UserService    *app.UserService
	grpcServer     *grpc.Server
	httpServer     *gin.Engine
}

// NewContext instantiates new rte context object.
func NewContext(cfg *EnvConfig) *Context {
	return &Context{
		Closers: []func(){},
		cfg:     cfg,
	}
}
