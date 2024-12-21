package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/exp/slog"

	"github.com/xaviercrochet/turbo-octo-adventure/api/app"
)

var (
	// flags to be provided for running the example server
	domain = flag.String("domain", "", "your ZITADEL instance domain (in the form: <instance>.zitadel.cloud or <yourdomain>)")
	key    = flag.String("key", "", "path to your key.json")
	port   = flag.String("port", "8090", "port to run the server on (default is 8090)")
)

/*
 This example demonstrates how to secure an HTTP API with ZITADEL using the provided authorization (AuthZ) middleware.

 It will serve the following 3 different endpoints:
 (These are meant to demonstrate the possibilities and do not follow REST best practices):

 - /api/healthz (can be called by anyone)
 - /api/feed (requires authorization)
 - /api/select_feed (requires authorization with granted `admin` role)
*/

func main() {
	flag.Parse()
	ctx := context.Background()

	serverOptions := app.NewServerOptions(*domain, *key, *port)
	router := http.NewServeMux()
	if err := app.SetupRoutes(ctx, router, serverOptions); err != nil {
		slog.Error("could not start server", "error", err)
		os.Exit(1)
	}

	// start the server on the specified port (default http://localhost:8101)
	lis := fmt.Sprintf(":%s", *port)
	slog.Info("server listening, press ctrl+c to stop", "addr", "http://localhost"+lis)
	err := http.ListenAndServe(lis, router)
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server terminated", "error", err)
		os.Exit(1)
	}
}
