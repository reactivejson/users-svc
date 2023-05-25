// main.go
package main

import (
	"context"
	"fmt"
	"github.com/reactivejson/users-svc/internal/domain"
	"log"
	"os"
	"os/signal"
	"syscall"

	conf "github.com/reactivejson/users-svc/cmd/app"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */
func main() {
	ctx, cancelCtxFn := context.WithCancel(context.Background())

	cfg := conf.SetupEnvConfig()

	fmt.Printf("KafkaBrokers %s\n", cfg.KafkaBrokers)
	fmt.Printf("GrpcPort %s\n", cfg.GrpcPort)

	// Intercepting shutdown signals.
	go waitForSignal(ctx, cancelCtxFn)

	appCtx := conf.NewContext(cfg)
	logExitMsg(conf.Setup(ctx, appCtx))
}

type mockNotifier struct {
	notifyError error
}

func (m *mockNotifier) NotifyUserChange(user *domain.User) error {
	return m.notifyError
}

func waitForSignal(ctx context.Context, cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-signals:
		log.Printf("received signal: %s, exiting gracefully", s)
		cancel()
	case <-ctx.Done():
		log.Printf("Service context done, serving remaining requests and exiting.")
	}
}

func logExitMsg(err error) {
	if err != nil {
		log.Fatalf("Service failed to setup: %s", err)
	}

	log.Print("Service exited successfully")
}
