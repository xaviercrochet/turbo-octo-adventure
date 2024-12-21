package main

import (
	"context"
	_ "embed"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/xaviercrochet/turbo-octo-adventure/web"
)

var (
	// flags to be provided for running the example server
	domain      = flag.String("domain", "", "your ZITADEL instance domain (in the form: https://<instance>.zitadel.cloud or https://<yourdomain>)")
	apiHostname = flag.String("apiHostname", "localhost", "hostname of the api")
	apiPort     = flag.String("apiPort", "8090", "port of the api")
	key         = flag.String("key", "", "encryption key")
	clientID    = flag.String("clientID", "", "clientID provided by ZITADEL")
	redirectURI = flag.String("redirectURI", "", "redirectURI registered at ZITADEL")
	port        = flag.String("port", "8089", "port to run the server on (default is 8089)")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	base64Key, err := base64.StdEncoding.DecodeString(*key)
	if err != nil {
		slog.Error("unable to decode aes key", "error", err)
		os.Exit(1)
	}

	if err != nil {
		slog.Error("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}

	router := http.NewServeMux()
	options := web.NewServerOptions(base64Key, *apiHostname, *apiPort, *domain, *clientID, *redirectURI)
	if err := web.SetupRoutes(ctx, router, options); err != nil {
		slog.Error("could not setup routes", "error", err)
		os.Exit(1)

	}

	lis := fmt.Sprintf(":%s", *port)
	slog.Info("server listening, press ctrl+c to stop", "addr", "http://localhost"+lis)
	err = http.ListenAndServe(lis, router)
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server terminated", "error", err)
		os.Exit(1)
	}
}
