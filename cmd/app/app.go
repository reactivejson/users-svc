package app

import (
	"context"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

// setupFn is an options function receiving app
type setupFn func(*Context) error

// Setup is preparing the app to be started
func Setup(ctx context.Context, cfg *Context) error {
	defer runClosers(cfg.Closers) // close closers at the end

	setupFuncs := []setupFn{
		setupDB(),
		setupKafka(),
		setupUserService(),
		setupGrpcServer(),
		setupHttpServer(),
	}
	return runSetupFncs(setupFuncs, cfg)
}
func runSetupFncs(fncs []setupFn, c *Context) error {
	for _, f := range fncs {
		if err := f(c); err != nil {
			return err
		}
	}
	return nil
}

func runClosers(fncs []func()) {
	for _, f := range fncs {
		f()
	}
}
