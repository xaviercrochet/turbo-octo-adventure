package web

import (
	"context"
	"embed"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"text/template"

	mw "github.com/xaviercrochet/turbo-octo-adventure/pkg/middleware"
	"github.com/xaviercrochet/turbo-octo-adventure/pkg/net"
	"github.com/xaviercrochet/turbo-octo-adventure/web/feed_api"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	openid "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

//go:embed "templates/*.html"
var templates embed.FS

// config values for the server
type ServerOptions struct {
	base64Key   []byte
	domain      string
	clientID    string
	redirectURI string
	apiHostname string
	apiPort     string
}

func NewServerOptions(base64Key []byte, apiHostname, apiPort, domain, clientID, redirectURI string) *ServerOptions {
	return &ServerOptions{
		base64Key:   base64Key,
		domain:      domain,
		clientID:    clientID,
		redirectURI: redirectURI,
		apiHostname: apiHostname,
		apiPort:     apiPort,
	}
}

/*
- Setup the authentication context and its middleware and the routes of the web application
*/
func SetupRoutes(ctx context.Context, router *http.ServeMux, options *ServerOptions) error {

	// load html tempates
	t, err := template.New("").ParseFS(templates, "templates/*.html")
	if err != nil {
		return fmt.Errorf("unable to parse template: %v", err)
	}

	//setup authentication context
	authN, err := authentication.New(ctx, zitadel.New(options.domain), string(options.base64Key),
		openid.DefaultAuthentication(options.clientID, options.redirectURI, string(options.base64Key)),
	)
	if err != nil {
		return fmt.Errorf("zitadel sdk could not initialize: %v", err)
	}

	//initialize the authentication middleware
	authMw := authentication.Middleware(authN)

	// default authentication routes provided by the sdk
	router.Handle("/auth/", mw.LogMiddleware(authN))

	/*
	   This endpoint
	   - is only accessible with a valid authentication
	   - is only accessible for admin users
	   - only accepts POST requests
	   - integrate the  /select_feed feed api endpoint to change from which user the feed is retrieved for

	   if the request is successfull, the user is redirected to /feed
	*/
	router.Handle("/select_feed", mw.LogMiddleware(authMw.RequireAuthentication()(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		authCtx := authMw.Context(req.Context())

		// deserialize request payload
		err := req.ParseForm()
		if err != nil {
			http.Error(w, "failed to parse form data", http.StatusBadRequest)
			return
		}

		// validate user input
		name := req.FormValue("name")
		if name == "" {
			http.Error(w, "name can't be empty", http.StatusBadRequest)
		}
		// sanitize user input
		name = html.EscapeString(name)

		// http client that integrate the feed api
		feedClient := feed_api.NewFeedClient(options.apiHostname, options.apiPort)
		if err := feedClient.SelectFeed(name, authCtx.Tokens.AccessToken); err == net.ErrNoAccess {
			slog.Error("select feed api call failed", "error", err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		} else if err == net.ErrNotAuthenticated {
			slog.Error("select feed api call failed", "error", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else if err != nil {
			slog.Error("select feed api call failed", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// browser expect a http status 3XX if we want to redirect after a successfull post
		http.Redirect(w, req, "/feed", http.StatusSeeOther)
	}))))

	/*
	   This endpoint
	   - is only accessible with a valid authentication
	   - is only accessible for users with any role
	   - only accepts GET requests
	   - integrate the  /feed feed api endpoint to retrieve the music feed
	   - renders feed.html

	*/

	router.Handle("/feed", mw.LogMiddleware(authMw.RequireAuthentication()(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Using the [middleware.Context] function we can gather information about the authenticated user.
		// This example will just print a JSON representation of the UserInfo of the typed [*oidc.UserInfoContext].
		authCtx := authMw.Context(req.Context())
		feedPage := NewFeedPage(authCtx.UserInfo.GivenName, authCtx.UserInfo.FamilyName)

		/*
		  check if health API is healthy
		  ideally, this should be part of a middleware
		*/

		feedClient := feed_api.NewFeedClient(options.apiHostname, options.apiPort)
		if ok, err := feedClient.CheckHealth(); !ok {
			feedPage.Health = false

			if err != nil {
				slog.Error("feed api is down or unresponsive", "error", err)
			}
		} else {
			// only query for feed if feed API is healthy
			feed, err := feedClient.GetFeed(authCtx.Tokens.AccessToken)
			// if the token is not valid anymore, log out the user
			if err == net.ErrNotAuthenticated {
				http.Redirect(w, req, "/auth/logout", http.StatusFound)
				return
			} else if err != nil {
				slog.Error("feed api call failed", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			feedPage.Feed = feed
		}

		err = t.ExecuteTemplate(w, "feed.html", feedPage)
		if err != nil {
			slog.Error("error writing feed response", "error", err)
		}
	}))))

	// This endpoint is accessible by anyone, but it will check if there already is a valid session (authentication).
	// If there is an active session, the information will be put into the context for later retrieval.
	router.Handle("/", mw.LogMiddleware(authMw.CheckAuthentication()(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// redirect the user to /feed in case he is already authenticated
		if authentication.IsAuthenticated(req.Context()) {
			http.Redirect(w, req, "/feed", http.StatusFound)
			return
		}

		err = t.ExecuteTemplate(w, "home.html", nil)
		if err != nil {
			slog.Error("error writing home page response", "error", err)
		}
	}))))
	return nil
}

// Represent the state of the feed.html page
type FeedPage struct {
	// informations about the current logged in user
	LoggedInUser string
	// Is the feed api running
	Health bool
	// the feed data
	Feed *feed_api.FeedResponse
}

func NewFeedPage(firstName, lastName string) *FeedPage {
	return &FeedPage{
		LoggedInUser: fmt.Sprintf("%s %s", firstName, lastName),
		Health:       true,
	}
}
