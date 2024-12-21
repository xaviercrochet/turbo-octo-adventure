package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/xaviercrochet/turbo-octo-adventure/api/musicbrainz"
	"github.com/xaviercrochet/turbo-octo-adventure/pkg/net"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

var (
	// store, in memory, the username from which the feed will be retrieved from the musicbrainz api
	selectedUsername = "xcrochet"
)

type SelectedFeed struct {
	Name string `json:"name"`
}

// configuration options for the server
type ServerOptions struct {
	domain      string
	keyFilePath string
	port        string
}

func NewServerOptions(domain, keyFilePath, port string) *ServerOptions {
	return &ServerOptions{
		domain:      domain,
		keyFilePath: keyFilePath,
		port:        port,
	}
}

/*
- Setup the authentication context and its middleware and the routes of the api
*/

func SetupRoutes(ctx context.Context, router *http.ServeMux, options *ServerOptions) error {
	//setup authorziation context
	authZ, err := authorization.New(ctx, zitadel.New(options.domain), oauth.DefaultAuthorization(options.keyFilePath))
	if err != nil {
		return fmt.Errorf("zitadel sdk could not initialize: %v", err)
	}

	// initialize the authorization middleware
	mw := middleware.New(authZ)

	// This endpoint is accessible by anyone and will always return "200 OK" to indicate the API is running
	router.Handle("/api/healthz", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			slog.Info("/api/healthz", "verb", r.Method)
			err = jsonResponse(w, "OK", http.StatusOK)
			if err != nil {
				slog.Error("error writing response", "error", err)
			}
		}))

	/*

	   Update selectedFeed

	   Request body: see SelectedFeed

	   - user need to be authenticated
	   - user is authorized with admin role

	   Response:
	   - 401 if user is not authenticated
	   - 403 if user is not authorized
	   - 404 if http verb is not POST or username doesn't exist
	*/

	router.Handle("/api/select_feed", mw.RequireAuthorization()(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			slog.Info("/api/select_feed", "verb", r.Method)
			if r.Method != http.MethodPost {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}

			authCtx := mw.Context(r.Context())
			if !authCtx.IsGrantedRole("admin") {
				slog.Warn("user doesn't have access to the resource", "id", authCtx.UserID(), "username", authCtx.Username)
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			// deserialize the request payload
			body, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Warn("could not read request body", "id", authCtx.UserID(), "username", authCtx.Username, "error", err)
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			var selectedFeed SelectedFeed
			err = json.Unmarshal(body, &selectedFeed)
			if err != nil {
				slog.Warn("could not deserialize request body", "id", authCtx.UserID(), "username", authCtx.Username, "error", err)
				http.Error(w, "failed to deserialize request body", http.StatusBadRequest)
				return
			}

			// update the username from wich '/feed' will retrieve the musicbrainz feed from
			selectedUsername = selectedFeed.Name

			// OK
			err = jsonResponse(w, "OK", http.StatusOK)
			if err != nil {
				slog.Error("error writing response", "error", err)
			}
		})))

	/*
	   Retrieve music feed from feed API, based on selectedUsername
	   - user need to be authenticated
	   - user is authorized with any role

	   Response:
	   - 401 if not authenticated
	   - 403 if not authorized
	   - 404 if http verb is not GET
	*/

	router.Handle("/api/feed", mw.RequireAuthorization()(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			slog.Info("/api/feed", "verb", r.Method)
			if r.Method != http.MethodGet {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}

			authCtx := mw.Context(r.Context())
			slog.Info("retrieving user feed", "id", authCtx.UserID(), "username", authCtx.Username, "feed_username", selectedUsername)

			// retrieve music feed from musicbrainz API
			feed, err := musicbrainz.GetFeed(selectedUsername)

			// handle client error if any
			if err == net.ErrNotFound {
				http.Error(w, "feed not found", http.StatusNotFound)
			} else if err != nil {
				slog.Warn("musicbrainz api call failed", "error", err)
				http.Error(w, "musicbrainz api call failed", http.StatusInternalServerError)
			}

			/*
			   Return the feed and weather the user has sufficiant authorization to update selectedUsername
			*/
			resp := &FeedResponse{
				Feed:        feed,
				WriteAccess: authCtx.IsGrantedRole("admin"),
			}

			err = jsonResponse(w, resp, http.StatusOK)
			if err != nil {
				slog.Error("error writing response", "error", err)
			}
		})))

	return nil
}

func jsonResponse(w http.ResponseWriter, resp any, status int) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

type FeedResponse struct {
	WriteAccess bool              `json:"write_access"`
	Feed        *musicbrainz.Feed `json:"feed"`
}
